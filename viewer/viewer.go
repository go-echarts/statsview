package viewer

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"text/template"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

// Metrics
type Metrics struct {
	Values []float64 `json:"values"`
	Time   string    `json:"time"`
}

type config struct {
	Interval   int
	MaxPoints  int
	Template   string
	ListenAddr string
	LinkAddr   string
	TimeFormat string
	Theme      Theme
}

type Theme string

const (
	ThemeWesteros Theme = types.ThemeWesteros
	ThemeMacarons Theme = types.ThemeMacarons
)

const (
	DefaultTemplate = `
$(function () { setInterval({{ .ViewID }}_sync, {{ .Interval }}); });
function {{ .ViewID }}_sync() {
    $.ajax({
        type: "GET",
        url: "http://{{ .Addr }}/debug/statsview/view/{{ .Route }}",
        dataType: "json",
        success: function (result) {
            let opt = goecharts_{{ .ViewID }}.getOption();

            let x = opt.xAxis[0].data;
            x.push(result.time);
            if (x.length > {{ .MaxPoints }}) {
                x = x.slice(1);
            }
            opt.xAxis[0].data = x;

            for (let i = 0; i < result.values.length; i++) {
                let y = opt.series[i].data;
                y.push({ value: result.values[i] });
                if (y.length > {{ .MaxPoints }}) {
                    y = y.slice(1);
                }
                opt.series[i].data = y;

                goecharts_{{ .ViewID }}.setOption(opt);
            }
        }
    });
}`
	DefaultMaxPoints  = 30
	DefaultTimeFormat = "15:04:05"
	DefaultInterval   = 2000
	DefaultAddr       = "localhost:18066"
	DefaultTheme      = ThemeMacarons
)

var defaultCfg = &config{
	Interval:   DefaultInterval,
	MaxPoints:  DefaultMaxPoints,
	Template:   DefaultTemplate,
	ListenAddr: DefaultAddr,
	LinkAddr:   DefaultAddr,
	TimeFormat: DefaultTimeFormat,
	Theme:      DefaultTheme,
}

type Option func(c *config)

// Addr returns the default server listening address
func Addr() string {
	return defaultCfg.ListenAddr
}

// LinkAddr returns the default html link address
func LinkAddr() string {
	return defaultCfg.LinkAddr
}

// Interval returns the default collecting interval of ViewManager
func Interval() int {
	return defaultCfg.Interval
}

// WithInterval sets the interval of collecting and pulling metrics
func WithInterval(interval int) Option {
	return func(c *config) {
		c.Interval = interval
	}
}

// WithMaxPoints sets the maximum points of each chart series
func WithMaxPoints(n int) Option {
	return func(c *config) {
		c.MaxPoints = n
	}
}

// WithTemplate sets the rendered template which fetching stats from the server and
// handling the metrics data
func WithTemplate(t string) Option {
	return func(c *config) {
		c.Template = t
	}
}

// WithAddr sets the listening address and link address
func WithAddr(addr string) Option {
	return func(c *config) {
		c.ListenAddr = addr
		c.LinkAddr = addr
	}
}

// WithLinkAddr sets the html link address
func WithLinkAddr(addr string) Option {
	return func(c *config) {
		c.LinkAddr = addr
	}
}

// WithTimeFormat sets the time format for the line-chart Y-axis label
func WithTimeFormat(s string) Option {
	return func(c *config) {
		c.TimeFormat = s
	}
}

// WithTheme sets the theme of the charts
func WithTheme(theme Theme) Option {
	return func(c *config) {
		c.Theme = theme
	}
}

func SetConfiguration(opts ...Option) {
	for _, opt := range opts {
		opt(defaultCfg)
	}
}

// Viewer is the abstraction of a Graph which in charge of collecting metrics from somewhere
type Viewer interface {
	Name() string
	View() *charts.Line
	Serve(w http.ResponseWriter, _ *http.Request)
	SetStatsMgr(smgr *StatsMgr)
}

type statsEntity struct {
	Stats *runtime.MemStats
	T     string
}

var memstats = &statsEntity{Stats: &runtime.MemStats{}}

type StatsMgr struct {
	last   int64
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewStatsMgr(ctx context.Context) *StatsMgr {
	s := &StatsMgr{}
	s.Ctx, s.Cancel = context.WithCancel(ctx)
	go s.polling()

	return s
}

func (s *StatsMgr) Tick() {
	s.last = time.Now().Unix() + int64(float64(Interval())/1000.0)*2
}

func (s *StatsMgr) polling() {
	ticker := time.NewTicker(time.Duration(Interval()) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.last > time.Now().Unix() {
				runtime.ReadMemStats(memstats.Stats)
				memstats.T = time.Now().Format(defaultCfg.TimeFormat)
			}
		case <-s.Ctx.Done():
			return
		}
	}
}

func genViewTemplate(vid, route string) string {
	tpl, err := template.New("view").Parse(defaultCfg.Template)
	if err != nil {
		panic("statsview: failed to parse template " + err.Error())
	}

	var c = struct {
		Interval  int
		MaxPoints int
		Addr      string
		Route     string
		ViewID    string
	}{
		Interval:  defaultCfg.Interval,
		MaxPoints: defaultCfg.MaxPoints,
		Addr:      defaultCfg.LinkAddr,
		Route:     route,
		ViewID:    vid,
	}

	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, c); err != nil {
		panic("statsview: failed to execute template " + err.Error())
	}

	return buf.String()
}

func fixedPrecision(n float64, p int) float64 {
	var r float64
	switch p {
	case 2:
		r, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", n), 64)
	case 6:
		r, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", n), 64)
	}
	return r
}

func newBasicView(route string) *charts.Line {
	graph := charts.NewLine()
	graph.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Time"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "600px",
			Height: "400px",
			Theme:  string(defaultCfg.Theme),
		}),
	)
	graph.SetXAxis([]string{}).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	graph.AddJSFuncs(genViewTemplate(graph.ChartID, route))
	return graph
}

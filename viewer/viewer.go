package viewer

import (
	"bytes"
	"net/http"
	"runtime"
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
	Addr       string
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
	DefaultMaxPoints  = 40
	DefaultTimeFormat = "15:04:05"
	DefaultInterval   = 1500
	DefaultAddr       = "localhost:18066"
	DefaultTheme      = ThemeMacarons
)

var defaultCfg = &config{
	Interval:   DefaultInterval,
	MaxPoints:  DefaultMaxPoints,
	Template:   DefaultTemplate,
	Addr:       DefaultAddr,
	TimeFormat: DefaultTimeFormat,
	Theme:      DefaultTheme,
}

type Option func(c *config)

// Addr returns the default server listening address
func Addr() string {
	return defaultCfg.Addr
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

// WithAddr sets the listening address
func WithAddr(addr string) Option {
	return func(c *config) {
		c.Addr = addr
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
}

type statsEntity struct {
	T     string
	Stats *runtime.MemStats
}

var rtStats = &statsEntity{Stats: &runtime.MemStats{}}

func StartRTCollect() {
	runtime.ReadMemStats(rtStats.Stats)
	rtStats.T = time.Now().Format(defaultCfg.TimeFormat)
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
		Addr:      defaultCfg.Addr,
		Route:     route,
		ViewID:    vid,
	}

	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, c); err != nil {
		panic("statsview: failed to execute template " + err.Error())
	}

	return buf.String()
}

func newBasicView(route string) *charts.Line {
	graph := charts.NewLine()
	graph.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
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

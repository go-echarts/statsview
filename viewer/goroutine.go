package viewer

import (
	"encoding/json"
	"net/http"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	// VGoroutine is the name of GoroutinesViewer
	VGoroutine = "goroutine"
)

// GoroutinesViewer collects the goroutine number metric via `runtime.NumGoroutine()`
type GoroutinesViewer struct {
	smgr  *StatsMgr
	graph *charts.Line
}

// NewGoroutinesViewer returns the GoroutinesViewer instance
// Series: Goroutines
func NewGoroutinesViewer() Viewer {
	graph := NewBasicView(VGoroutine)
	graph.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{Name: "Num"}),
		charts.WithTitleOpts(opts.Title{Title: "Goroutines"}),
	)
	graph.AddSeries("Goroutines", []opts.LineData{})
	graph.AddSeries("Threads", []opts.LineData{})
	graph.AddSeries("NumCPU", []opts.LineData{})
	graph.AddSeries("GOMAXPROCS", []opts.LineData{})

	return &GoroutinesViewer{graph: graph}
}

func (vr *GoroutinesViewer) SetStatsMgr(smgr *StatsMgr) {
	vr.smgr = smgr
}

func (vr *GoroutinesViewer) Name() string {
	return VGoroutine
}

func (vr *GoroutinesViewer) View() *charts.Line {
	return vr.graph
}

var threadProfile = pprof.Lookup("threadcreate")

func (vr *GoroutinesViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	vr.smgr.Tick()

	metrics := Metrics{
		Values: []float64{
			float64(runtime.NumGoroutine()),
			float64(threadProfile.Count()),
			float64(runtime.NumCPU()),
			float64(runtime.GOMAXPROCS(0)),
		},
		Time: time.Now().Format(DefaultCfg.TimeFormat),
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

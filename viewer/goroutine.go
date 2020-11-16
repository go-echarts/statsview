package viewer

import (
	"encoding/json"
	"net/http"
	"runtime"
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
	graph *charts.Line
}

// NewGoroutinesViewer returns the GoroutinesViewer instance
// Series: Goroutines
func NewGoroutinesViewer() Viewer {
	graph := newBasicView(VGoroutine)
	graph.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{Name: "Num"}),
		charts.WithTitleOpts(opts.Title{Title: "Goroutines"}),
	)
	graph.AddSeries("Goroutines", []opts.LineData{})

	return &GoroutinesViewer{graph: graph}
}

func (vr *GoroutinesViewer) Name() string {
	return VGoroutine
}

func (vr *GoroutinesViewer) View() *charts.Line {
	return vr.graph
}

func (vr *GoroutinesViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	metrics := Metrics{
		Values: []float64{float64(runtime.NumGoroutine())},
		Time:   time.Now().Format(defaultCfg.TimeFormat),
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

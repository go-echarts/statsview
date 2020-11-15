package viewer

import (
	"encoding/json"
	"net/http"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/opts"
)

const (
	VHeap = "heap"
)

// HeapViewer collects the heap-stats metrics from `runtime.ReadMemStats()`
// including `HeapAlloc` / `HeapInuse` / `HeapSys` / `HeapIdle`
type HeapViewer struct {
	graph *charts.Line
}

func NewHeapViewer() Viewer {
	graph := newBasicView(VHeap)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Heap"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Size", AxisLabel: &opts.AxisLabel{Formatter: "{value} MB"}}),
	)
	graph.AddSeries("Alloc", []opts.LineData{}).
		AddSeries("Inuse", []opts.LineData{}).
		AddSeries("Sys", []opts.LineData{}).
		AddSeries("Idle", []opts.LineData{})

	return &HeapViewer{graph: graph}
}

func (vr *HeapViewer) Name() string {
	return VHeap
}

func (vr *HeapViewer) View() *charts.Line {
	return vr.graph
}

func (vr *HeapViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	metrics := Metrics{
		Values: []float64{
			float64(rtStats.Stats.HeapAlloc) / 1024 / 1024,
			float64(rtStats.Stats.HeapInuse) / 1024 / 1024,
			float64(rtStats.Stats.HeapSys) / 1024 / 1024,
			float64(rtStats.Stats.HeapIdle) / 1024 / 1024,
		},
		Time: rtStats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

package viewer

import (
	"encoding/json"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	// VHeap is the name of HeapViewer
	VHeap = "heap"
)

// HeapViewer collects the heap-stats metrics via `runtime.ReadMemStats()`
type HeapViewer struct {
	smgr  *StatsMgr
	graph *charts.Line
}

// NewHeapViewer returns the HeapViewer instance
// Series: Alloc / Inuse / Sys / Idle
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
func (vr *HeapViewer) SetStatsMgr(smgr *StatsMgr) {
	vr.smgr = smgr
}
func (vr *HeapViewer) Name() string {
	return VHeap
}

func (vr *HeapViewer) View() *charts.Line {
	return vr.graph
}

func (vr *HeapViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	vr.smgr.Tick()

	metrics := Metrics{
		Values: []float64{
			fixedPrecision(float64(memstats.Stats.HeapAlloc)/1024/1024, 2),
			fixedPrecision(float64(memstats.Stats.HeapInuse)/1024/1024, 2),
			fixedPrecision(float64(memstats.Stats.HeapSys)/1024/1024, 2),
			fixedPrecision(float64(memstats.Stats.HeapIdle)/1024/1024, 2),
		},
		Time: memstats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

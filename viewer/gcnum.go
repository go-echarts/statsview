package viewer

import (
	"encoding/json"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	// VGCNum is the name of GCNumViewer
	VGCNum = "gcnum"
)

// GCNumViewer collects the GC number metric via `runtime.ReadMemStats()`
type GCNumViewer struct {
	graph *charts.Line
}

// NewGCNumViewer returns the GCNumViewer instance
// Series: GcNum
func NewGCNumViewer() Viewer {
	graph := newBasicView(VGCNum)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "GC Number"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Num"}),
	)
	graph.AddSeries("GcNum", []opts.LineData{})

	return &GCNumViewer{graph: graph}
}

func (vr *GCNumViewer) Name() string {
	return VGCNum
}

func (vr *GCNumViewer) View() *charts.Line {
	return vr.graph
}

func (vr *GCNumViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	metrics := Metrics{
		Values: []float64{float64(rtStats.Stats.NumGC)},
		Time:   rtStats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

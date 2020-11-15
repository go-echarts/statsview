package viewer

import (
	"encoding/json"
	"net/http"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/opts"
)

const (
	VGCCPUFraction = "gccpufraction"
)

// GCCPUFractionViewer collects the GC-CPU fraction metric from `runtime.ReadMemStats()`
type GCCPUFractionViewer struct {
	graph *charts.Line
}

func NewGCCPUFractionViewer() Viewer {
	graph := newBasicView(VGCCPUFraction)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "GC CPUFraction"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Percent", AxisLabel: &opts.AxisLabel{Formatter: "{value} %", Rotate: 30}}),
	)
	graph.AddSeries("Fraction", []opts.LineData{})

	return &GCCPUFractionViewer{graph: graph}
}

func (vr *GCCPUFractionViewer) Name() string {
	return VGCCPUFraction
}

func (vr *GCCPUFractionViewer) View() *charts.Line {
	return vr.graph
}

func (vr *GCCPUFractionViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	metrics := Metrics{
		Values: []float64{rtStats.Stats.GCCPUFraction},
		Time:   rtStats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

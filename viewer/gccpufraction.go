package viewer

import (
	"encoding/json"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	// VGCCPUFraction is the name of GCCPUFractionViewer
	VGCCPUFraction = "gccpufraction"
)

// GCCPUFractionViewer collects the GC-CPU fraction metric via `runtime.ReadMemStats()`
type GCCPUFractionViewer struct {
	smgr  *StatsMgr
	graph *charts.Line
}

// NewGCCPUFractionViewer returns the GCCPUFractionViewer instance
// Series: Fraction
func NewGCCPUFractionViewer() Viewer {
	graph := newBasicView(VGCCPUFraction)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "GC CPUFraction"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Percent", AxisLabel: &opts.AxisLabel{Formatter: "{value} %", Rotate: 35}}),
	)
	graph.AddSeries("Fraction", []opts.LineData{})

	return &GCCPUFractionViewer{graph: graph}
}

func (vr *GCCPUFractionViewer) SetStatsMgr(smgr *StatsMgr) {
	vr.smgr = smgr
}

func (vr *GCCPUFractionViewer) Name() string {
	return VGCCPUFraction
}

func (vr *GCCPUFractionViewer) View() *charts.Line {
	return vr.graph
}

func (vr *GCCPUFractionViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	vr.smgr.Tick()

	metrics := Metrics{
		Values: []float64{fixedPrecision(memstats.Stats.GCCPUFraction, 6)},
		Time:   memstats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

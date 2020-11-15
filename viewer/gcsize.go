package viewer

import (
	"encoding/json"
	"net/http"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/opts"
)

const (
	VGCSzie = "gcsize"
)

// GCSizeViewer collects the GC size metric from `runtime.ReadMemStats()`
type GCSizeViewer struct {
	graph *charts.Line
}

func NewGCSizeViewer() Viewer {
	graph := newBasicView(VGCSzie)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "GC Size"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Size", AxisLabel: &opts.AxisLabel{Formatter: "{value} MB"}}),
	)
	graph.AddSeries("GCSys", []opts.LineData{}).
		AddSeries("NextGC", []opts.LineData{})

	return &GCSizeViewer{graph: graph}
}

func (vr *GCSizeViewer) Name() string {
	return VGCSzie
}

func (vr *GCSizeViewer) View() *charts.Line {
	return vr.graph
}

func (vr *GCSizeViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	metrics := Metrics{
		Values: []float64{
			float64(rtStats.Stats.GCSys) / 1024 / 1024,
			float64(rtStats.Stats.NextGC) / 1024 / 1024,
		},
		Time: rtStats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}

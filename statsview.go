package statsview

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/templates"
	"github.com/rs/cors"

	"github.com/go-echarts/statsview/statics"
	"github.com/go-echarts/statsview/viewer"
)

// ViewManager
type ViewManager struct {
	Views []viewer.Viewer

	srv  *http.Server
	done chan struct{}
}

// Register registers views to the ViweManager
func (vm *ViewManager) Register(views ...viewer.Viewer) {
	vm.Views = append(vm.Views, views...)
}

// Start runs a http server and begin to collect metrics
func (vm *ViewManager) Start() {
	ticker := time.NewTicker(time.Duration(viewer.Interval()) * time.Millisecond)

	go func() {
		vm.srv.ListenAndServe()
	}()

	for {
		select {
		case <-ticker.C:
			viewer.StartRTCollect()
		case <-vm.done:
			vm.Stop()
			ticker.Stop()
			return
		}
	}
}

// Stop shutdown the http server gracefully
func (vm *ViewManager) Stop() {
	vm.done <- struct{}{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	vm.srv.Shutdown(ctx)
}

func init() {
	templates.PageTpl = `
{{- define "page" }}
<!DOCTYPE html>
<html>
    {{- template "header" . }}
<body>
<p>&nbsp;&nbsp;🚀 <a href="https://github.com/go-echarts/statsview"><b>StatsView</b></a> <em>is a real-time Golang runtime stats visualization profiler</em></p>
<style> .box { justify-content:center; display:flex; flex-wrap:wrap } </style>
<div class="box"> {{- range .Charts }} {{ template "base" . }} {{- end }} </div>
</body>
</html>
{{ end }}
`
}

// New creates a new ViewManager instance
func New() *ViewManager {
	page := components.NewPage()
	page.AssetsHost = fmt.Sprintf("http://%s/debug/statsview/statics/", viewer.Addr())
	page.Assets.JSAssets.Add("jquery.min.js")

	mgr := &ViewManager{
		done: make(chan struct{}),
		srv: &http.Server{
			Addr:           viewer.Addr(),
			ReadTimeout:    time.Minute,
			WriteTimeout:   time.Minute,
			MaxHeaderBytes: 1 << 20,
		},
	}

	mgr.Register(
		viewer.NewGoroutinesViewer(),
		viewer.NewHeapViewer(),
		viewer.NewStackViewer(),
		viewer.NewGCNumViewer(),
		viewer.NewGCSizeViewer(),
		viewer.NewGCCPUFractionViewer(),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	for _, v := range mgr.Views {
		page.AddCharts(v.View())
		mux.HandleFunc("/debug/statsview/view/"+v.Name(), v.Serve)
	}

	mux.HandleFunc("/debug/statsview", func(w http.ResponseWriter, _ *http.Request) {
		page.Render(w)
	})

	staticsPrev := "/debug/statsview/statics/"
	mux.HandleFunc(staticsPrev+"echarts.min.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.EchartJS))
	})

	mux.HandleFunc(staticsPrev+"jquery.min.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.JqueryJS))
	})

	mux.HandleFunc(staticsPrev+"themes/westeros.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.WesterosJS))
	})

	mux.HandleFunc(staticsPrev+"themes/macarons.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.MacaronsJS))
	})

	mgr.srv.Handler = cors.Default().Handler(mux)
	return mgr
}

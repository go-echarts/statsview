# 🚀 Statsview

Statsview is a real-time Golang runtime stats visualization profiler. It is built top on another open-source project, [go-echarts](https://github.com/go-echarts/go-echarts), which helps statsview to show its graphs on the browser.

<a href="https://github.com/go-echarts/statsview/pulls">
    <img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="Contributions welcome">
</a>
<a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-brightgreen.svg" alt="MIT License">
</a>
<a href="https://pkg.go.dev/github.com/go-echarts/statsview">
    <img src="https://godoc.org/github.com/go-echarts/statsview?status.svg" alt="GoDoc">
 </a>

## Installation

```shell
$ go get -u github.com/go-echarts/statsview/...
```

## Usage

Statsview is quite simple to use. All static assets have been packaged into the project thus it can be ran on any environment.

```golang
import (
    "time"

    "github.com/go-echarts/statsview"
)

func main() {
    go func() {
        mgr := statsview.New()

        // Start() runs a HTTP server at `localhost:18066` by default.
        mgr.Start()

        // Stop() will shutdown the http server gracefully
        // mgr.Stop()
    }()

    // busy working....
    time.Sleep(time.Minute)
}

// Visit your browser at http://localhost:18066/statsview/debug
```

## Configuration

Statsview gets a variety of configurations for the users. Everyone could customize their favorite charts style.

```golang
// WithInterval sets the interval(in millisecond) of collecting and pulling metrics
// default -> 1500
WithInterval(interval int)

// WithMaxPoints sets the maximum points of each chart series
// default -> 40
WithMaxPoints(n int)

// WithTemplate sets the rendered template which fetching stats from the server and
// handling the metrics data
WithTemplate(t string)

// WithAddr sets the listen address
// default -> "localhost:18066"
WithAddr(addr string)

// WithTimeFormat sets the time format for the line-chart Y-axis label
// default -> "15:04:05"
WithTimeFormat(s string)

// WithTheme sets the theme of the charts
// default -> Macarons
//
// Optional:
// * ThemeWesteros
// * ThemeMacarons
WithTheme(theme Theme)
```

#### Set the options

```golang
import (
    "github.com/go-echarts/statsview/viewer"
)

// set configurations before calling the `Start()` method
viewer.SetConfiguration(viewer.WithTheme(viewer.ThemeWalden), view.WithAddr("localhost:8087"))
```

## Viewers

Viewer is the abstraction of a Graph which in charge of collecting metrics from somewhere. Statsview provides some default viewers as below.

* `GCCPUFractionViewer`
* `GCNumViewer`
* `GCSizeViewer`
* `GoroutinesViewer`
* `HeapViewer`
* `StackViewer`

Viewer wraps a go-echarts [Line instance](https://github.com/go-echarts/go-echarts/blob/master/charts/line.go) that means you can use all the options/features on it. To be honest, I think that is the most charming thing about this project.

## Snapshot

#### ThemeMacarons

![Macarons](https://user-images.githubusercontent.com/19553554/99192859-45943400-27b0-11eb-8096-8a9e76fba3a1.png)

#### ThemeWesteros

![Westeros](https://user-images.githubusercontent.com/19553554/99193211-78d7c280-27b2-11eb-96c8-cbcb6792e68a.png)


## License

MIT [©chenjiandongx](https://github.com/chenjiandongx)

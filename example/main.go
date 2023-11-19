package main

import (
	"os"
	"os/signal"

	"github.com/smallnest/statsview"
)

func main() {
	mgr := statsview.New()

	go mgr.Start()

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt, os.Kill)
	<-sg
}

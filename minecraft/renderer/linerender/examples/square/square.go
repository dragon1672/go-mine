package main

import (
	"context"
	"flag"
	"time"

	"github.com/dragon1672/go-mine/minecraft/renderer/linerender"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	w, err := linerender.MakeWindow(1000, 800, func(t time.Time, dt time.Duration) error {
		return nil
	},
		func(renderer *linerender.LineRenderer) {
			renderer.SetColor(0, 0, 1)
			renderer.DrawLine(10, 10, 10, 100)
			renderer.DrawLine(10, 100, 100, 100)
			renderer.DrawLine(100, 100, 100, 10)
			renderer.DrawLine(100, 10, 10, 10)

			renderer.SetColor(1, 0, 0)
			renderer.DrawPoint(500, 500, 200)

			//renderer.SetColor(0, 1, 0)
			//renderer.DrawPoint(20, 30, 20)
		})
	if err != nil {
		glog.Fatalf("Error making window: %v", err)
	}
	if err := w.Start(ctx); err != nil {
		glog.Fatalf("Error running window: %v", err)
	}
}

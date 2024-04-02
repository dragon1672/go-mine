package tickers

import (
	"context"
	"time"

	"github.com/golang/glog"
)

func MakeTicker(ctx context.Context, d time.Duration, f func(t time.Time, dt time.Duration) (bool, error)) (start func(), cleanup func()) {
	updateTicker := time.NewTicker(d)
	updateTickerExit := make(chan struct{})
	var lastTime time.Time

	cleanup = func() {
		updateTicker.Stop()
		updateTickerExit <- struct{}{}
	}
	start = func() {
		go func() {
			for {
				select {
				case <-ctx.Done():
					glog.InfoContext(ctx, "context done, stopping ticker")
				case <-updateTickerExit:
					return
				case timestamp := <-updateTicker.C:
					dt := timestamp.Sub(lastTime)
					lastTime = timestamp
					ok, err := f(timestamp, dt)
					if err != nil {
						glog.ErrorContextf(ctx, "Ticker encountered error %v", err)
						cleanup()
						return
					}
					if !ok {
						glog.InfoContext(ctx, "Ticker safely exiting")
						cleanup()
						return
					}
				}
			}
		}()
	}
	return start, cleanup
}

func StartTicker(ctx context.Context, d time.Duration, f func(t time.Time, dt time.Duration) (bool, error)) (cleanup func()) {
	start, cleanup := MakeTicker(ctx, d, f)
	start()
	return cleanup
}

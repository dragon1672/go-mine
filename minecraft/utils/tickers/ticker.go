package tickers

import (
	"time"

	"github.com/golang/glog"
)

func MakeTicker(d time.Duration, f func(t time.Time, dt time.Duration) (bool, error)) (start func(), cleanup func()) {
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
				case <-updateTickerExit:
					return
				case timestamp := <-updateTicker.C:
					dt := timestamp.Sub(lastTime)
					lastTime = timestamp
					ok, err := f(timestamp, dt)
					if err != nil {
						glog.Errorf("Ticker encoutered error %v", err)
						cleanup()
						return
					}
					if !ok {
						glog.Info("Ticker safely exiting")
						cleanup()
						return
					}
				}
			}
		}()
	}
	return start, cleanup
}

func StartTicker(d time.Duration, f func(t time.Time, dt time.Duration) (bool, error)) (cleanup func()) {
	start, cleanup := MakeTicker(d, f)
	start()
	return cleanup
}

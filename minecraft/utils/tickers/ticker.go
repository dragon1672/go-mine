package tickers

import (
	"time"
)

func MakeTicker(d time.Duration, f func(t time.Time, dt time.Duration)) (start func(), cleanup func()) {
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
					f(timestamp, dt)
				}
			}
		}()
	}
	return start, cleanup
}

func StartTicker(d time.Duration, f func(t time.Time, dt time.Duration)) (cleanup func()) {
	start, cleanup := MakeTicker(d, f)
	start()
	return cleanup
}

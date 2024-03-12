package tickers

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestStartTicker(t *testing.T) {
	t.Parallel()
	t.Run("Base Case", func(t *testing.T) {
		ctx := context.Background()
		m := sync.Mutex{}
		var times []time.Time
		var dts []time.Duration
		waiter := make(chan struct{})
		_ = StartTicker(ctx, time.Nanosecond, func(t time.Time, dt time.Duration) (bool, error) {
			m.Lock()
			times = append(times, t)
			dts = append(dts, dt)
			m.Unlock()
			if len(times) < 5 {
				return true, nil
			}
			waiter <- struct{}{}
			return false, nil
		})
		// Wait for ticker
		<-waiter
		if len(times) != len(dts) && len(dts) != 5 {
			t.Errorf("Expected 5 entries for time and dts, got time: %d dts: %d", len(times), len(dts))
		}
		if times[2].Sub(times[1]) != dts[2] {
			t.Errorf("dts expected to measure the time between calls")
		}
	})
}

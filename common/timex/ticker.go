// Package timex
// Date: 2023/4/10 13:06
// Author: Amu
// Description:
package timex

import (
	"time"
)

type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

type realTicker struct {
	*time.Ticker
}

// NewTicker returns a Ticker.
func NewTicker(d time.Duration) Ticker {
	return &realTicker{
		Ticker: time.NewTicker(d),
	}
}

func (r *realTicker) Chan() <-chan time.Time {
	return r.C
}

func (r *realTicker) Stop() {
	r.Ticker.Stop()
}

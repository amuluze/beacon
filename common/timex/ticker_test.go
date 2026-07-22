// Package timex
// Date: 2023/4/10 13:31
// Author: Amu
// Description:
package timex

import (
	"testing"
	"time"
)

func TestRealTickerDoTick(t *testing.T) {
	ticker := NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	var count int
	for range ticker.Chan() {
		count++
		if count > 5 {
			break
		}
	}
}

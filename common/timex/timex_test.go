// Package timex
// Date: 2023/4/19 16:50
// Author: Amu
// Description:
package timex

import (
	"fmt"
	"testing"
	"time"
)

func TestInt64ToTime(t *testing.T) {
	timestamp := time.Now().Unix()
	res := Int64ToTime(timestamp)
	fmt.Println(res)
}

func TestCurrentDate(t *testing.T) {
	date := CurrentDate("2006.01.02")
	fmt.Println(date)
}

func TestCurrentTime(t *testing.T) {
	ct := CurrentTime()
	fmt.Println(ct)
}

func TestCurrentTimestamp(t *testing.T) {
	stamp := CurrentTimestamp()
	fmt.Println(stamp)
}

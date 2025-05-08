// Package timectl
// Date: 2024/06/30 13:17:50
// Author: Amu
// Description:
package timectl

import (
	"context"
	"testing"
)

func TestGetSystemTime(t *testing.T) {
	time, err := GetTime(context.TODO())
	t.Log(time, err)
}

func TestGetSystemTimezone(t *testing.T) {
	timezone, err := GetTimeZone(context.TODO())
	t.Log(timezone, err)
}

func TestGetSystemTimezoneList(t *testing.T) {
	timezoneList, err := GetTimeZoneList(context.TODO())
	t.Log(timezoneList, err)
}

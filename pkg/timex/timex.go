// Package timex
// Date: 2023/4/19 16:49
// Author: Amu
// Description:
package timex

import "time"

// Int64ToTime 秒值时间戳转 time.Time
func Int64ToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// CurrentDate 获取当前日期
func CurrentDate(format string) string {
	if format != "" {
		return time.Now().Format(format)
	}
	return time.Now().Format("2006-01-02")
}

// CurrentTime 获取当前时间
func CurrentTime() string {
	return time.Now().Format("15:04:05")
}

// CurrentTimestamp 获取当前秒级时间戳
func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

// Package utils
// Date: 2024/3/29 14:50
// Author: Amu
// Description:
package utils

import (
	"fmt"
	"math"
)

// ConvertBytesToReadable 将字节数格式化为带单位的可读字符串。
// 循环以 len(units)-1 为上界，保证超大值（PB/EB 及以上）不会越界 panic。
func ConvertBytesToReadable(bytes float64) string {
	var units = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	var index int
	for index = 0; index < len(units)-1; index++ {
		if bytes < 1024 {
			break
		}
		bytes /= 1024
	}
	return fmt.Sprintf("%.2f", bytes) + " " + units[index]
}

// Decimal float64 四舍五入，保留两位小数
func Decimal(f float64) float64 {
	power := math.Pow(10, 2)
	return math.Round(f*power) / power
}

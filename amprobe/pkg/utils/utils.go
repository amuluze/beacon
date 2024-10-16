// Package utils
// Date: 2024/3/29 14:50
// Author: Amu
// Description:
package utils

import (
	"fmt"
	"math"
)

func ConvertBytesToReadable(bytes float64) string {
	var units = []string{"B", "KB", "MB", "GB", "TB"}
	var index int
	for index = 0; index < len(units); index++ {
		if bytes < 1024 {
			break
		}
		bytes /= 1024
	}
	return fmt.Sprintf("%.2f", bytes) + " " + units[index]
}

/**
 * Decimal: float64 四舍五入，保留两位小数
 */

func Decimal(f float64) float64 {
	power := math.Pow(10, 2)
	return math.Round(f*power) / power
}

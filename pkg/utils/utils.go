// Package utils
// Date: 2024/3/29 14:50
// Author: Amu
// Description:
package utils

import "fmt"

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

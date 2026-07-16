// Package utils
// Date: 2026/07/16
// Author: Amu
// Description: tests for utils helpers
package utils

import "testing"

func TestConvertBytesToReadable(t *testing.T) {
	kb := 1024.0
	mb := kb * 1024
	gb := mb * 1024
	tb := gb * 1024
	pb := tb * 1024

	tests := []struct {
		name  string
		bytes float64
		want  string
	}{
		{"B", 512, "512.00 B"},
		{"零", 0, "0.00 B"},
		{"KB", 2048, "2.00 KB"},
		{"MB", mb * 1.5, "1.50 MB"},
		{"GB", gb * 3, "3.00 GB"},
		{"TB", tb * 2, "2.00 TB"},
		{"PB 不越界", pb, "1.00 PB"},
		{"远超 PB 不 panic，clamp 到 EB", pb * 1024 * 1024, "1024.00 EB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertBytesToReadable(tt.bytes); got != tt.want {
				t.Errorf("ConvertBytesToReadable(%v) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestDecimal(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want float64
	}{
		{"向下", 1.234, 1.23},
		{"向上", 1.236, 1.24},
		{"整数", 2.5, 2.5},
		{"负数", -1.236, -1.24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decimal(tt.in); got != tt.want {
				t.Errorf("Decimal(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

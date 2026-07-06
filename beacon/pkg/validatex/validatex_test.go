// Package validatex
// Date: 2026/6/26
// Author: Amu
// Description: unit tests for struct validation
package validatex

import (
	"testing"
)

type validStruct struct {
	Name   string `validate:"required,gte=1,lte=64"`
	Age    int    `validate:"gte=0,lte=150"`
	Status int    `validate:"oneof=0 1"`
}

type maxLenStruct struct {
	Name string `validate:"required,lte=5"`
}

type oneofStruct struct {
	Role string `validate:"required,oneof=admin user guest"`
}

func TestValidateStruct_Valid(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want bool // true = no error
	}{
		{
			name: "all valid fields",
			data: validStruct{Name: "alice", Age: 30, Status: 1},
			want: true,
		},
		{
			name: "zero age valid",
			data: validStruct{Name: "bob", Age: 0, Status: 0},
			want: true,
		},
		{
			name: "missing required name",
			data: validStruct{Name: "", Age: 30, Status: 1},
			want: false,
		},
		{
			name: "name exceeds max",
			data: maxLenStruct{Name: "toolongname"},
			want: false,
		},
		{
			name: "name within max",
			data: maxLenStruct{Name: "ok"},
			want: true,
		},
		{
			name: "valid oneof",
			data: oneofStruct{Role: "admin"},
			want: true,
		},
		{
			name: "invalid oneof",
			data: oneofStruct{Role: "super"},
			want: false,
		},
		{
			name: "empty oneof (required)",
			data: oneofStruct{Role: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err == nil) != tt.want {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, !tt.want)
			}
		})
	}
}

func TestValidateStruct_ReturnsFirstError(t *testing.T) {
	// When multiple fields are invalid, only the first error should be returned
	data := validStruct{Name: "", Age: 999, Status: 5}
	err := ValidateStruct(data)
	if err == nil {
		t.Error("expected error for multiple invalid fields, got nil")
	}
	// The function should return a single error, not panic or return nil
}

func TestValidateStruct_EmptyStruct(t *testing.T) {
	// A struct with no validate tags should pass
	type emptyStruct struct {
		Name string
	}
	err := ValidateStruct(emptyStruct{Name: "anything"})
	if err != nil {
		t.Errorf("empty struct should pass validation: %v", err)
	}
}

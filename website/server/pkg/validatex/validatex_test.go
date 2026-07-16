// Package validatex
// Date: 2026/07/16
// Author: Amu
// Description: tests for ValidateStruct
package validatex

import "testing"

type sample struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0,lte=150"`
}

func TestValidateStruct_OK(t *testing.T) {
	if err := ValidateStruct(sample{Name: "a", Age: 10}); err != nil {
		t.Fatalf("合法输入不应报错: %v", err)
	}
}

func TestValidateStruct_Required(t *testing.T) {
	if err := ValidateStruct(sample{Name: "", Age: 10}); err == nil {
		t.Fatal("缺少必填字段应报错")
	}
}

func TestValidateStruct_Range(t *testing.T) {
	if err := ValidateStruct(sample{Name: "a", Age: 200}); err == nil {
		t.Fatal("超出范围的 Age 应报错")
	}
}

func TestValidateStruct_Nil(t *testing.T) {
	if err := ValidateStruct(nil); err != nil {
		t.Fatalf("nil 输入应直接放行: %v", err)
	}
}

func TestValidateStruct_PointerOK(t *testing.T) {
	if err := ValidateStruct(&sample{Name: "a", Age: 1}); err != nil {
		t.Fatalf("结构体指针应正常校验: %v", err)
	}
}

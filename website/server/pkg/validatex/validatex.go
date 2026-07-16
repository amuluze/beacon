// Package validatex
// Date: 2024/3/6 13:21
// Author: Amu
// Description:
package validatex

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate 进程级单例，避免每次请求重新构造校验器。
var validate = validator.New()

func init() {
	// 错误信息中的字段名使用 json tag，与对外 API 字段对齐。
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
}

// ValidateStruct 在系统边界对请求结构体做模式校验。
// 永不 panic：入参非法时返回包裹错误，校验失败时返回首条可读错误。
func ValidateStruct(data interface{}) error {
	if data == nil {
		return nil
	}
	err := validate.Struct(data)
	if err == nil {
		return nil
	}
	// InvalidValidationError 表示入参本身非法（例如传入了非 struct 指针）。
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return fmt.Errorf("invalid request payload: %w", err)
	}
	if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
		e := errs[0]
		return fmt.Errorf("field %q failed %q validation", e.Field(), e.Tag())
	}
	return fmt.Errorf("validation failed: %w", err)
}

// Package validatex
// Date: 2024/3/6 13:21
// Author: Amu
// Description:
package validatex

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// validate is a package-level singleton. validator.Validate is concurrency-safe.
var validate = validator.New()

// ValidateStruct 用于api validate 请求参数。
// 非结构体或 nil 入参返回 validator 的 InvalidValidationError，不再 panic。
func ValidateStruct(data interface{}) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, validationErr := range validationErrors {
			return validationErr
		}
	}
	return err
}

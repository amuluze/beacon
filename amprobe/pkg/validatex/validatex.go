// Package validatex
// Date: 2024/3/6 13:21
// Author: Amu
// Description:
package validatex

import (
	"amprobe/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// ValidateStruct 用于api validate 请求参数
func ValidateStruct(data interface{}) error {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New400Error(err.Error())
		}
		return errors.New400Error(err.Error())
	}
	return nil
}

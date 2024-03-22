package common

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"dot/pkg/util"
)

type validatorImpl struct {
	validate *validator.Validate
}

func NewValidator() *validatorImpl {
	return &validatorImpl{util.NewValidator()}
}

func (v *validatorImpl) Validate(schema interface{}) error {
	return v.validate.Struct(schema)
}

type ErrorBindRequestBodyMeta struct {
	Field  string
	Format string
}

func ErrorBindBodyRequest(string error) ErrorBindRequestBodyMeta {
	errorBind := strings.Split(string.Error(), ",")
	field := ""
	if len(errorBind) > 3 {
		field = strings.Split(errorBind[3], "=")[1]
	}
	format := strings.Split(errorBind[1], "=")[1]
	return ErrorBindRequestBodyMeta{
		Field:  field,
		Format: fmt.Sprintf("must be in %s", format),
	}
}

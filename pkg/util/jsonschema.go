package util

import (
	"dot/core"
	"github.com/xeipuuv/gojsonschema"
)

func JsonSchemaValidate(schema string, document any) *core.CustomError {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	docLoader := gojsonschema.NewRawLoader(document)

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: err.Error(),
		}
	}

	if !result.Valid() {
		return &core.CustomError{
			Code:    core.BAD_REQUEST,
			Message: err.Error(),
		}
	}

	return nil
}

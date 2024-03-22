package common

import (
	"dot/core"
	"net/http"
)

var ResultMap = map[int]Result{
	core.OK:                           {Code: core.OK, StatusCode: 200, Message: core.OK_M},
	core.DATA_NOT_FOUND:               {Code: core.DATA_NOT_FOUND, StatusCode: 200, Message: core.DATA_NOT_FOUND_M},
	core.DATA_NOT_UPDATED:             {Code: core.DATA_NOT_UPDATED, StatusCode: 200, Message: core.DATA_NOT_UPDATED_M},
	core.CREATED:                      {Code: core.CREATED, StatusCode: 200, Message: core.CREATED_M},
	core.CREATED_BAD_REQUEST:          {Code: core.CREATED_BAD_REQUEST, StatusCode: 201, Message: core.CREATED_BAD_REQUEST_M},
	core.DELETED:                      {Code: core.DELETED, StatusCode: 200, Message: core.DELETED_M},
	core.BAD_REQUEST:                  {Code: core.BAD_REQUEST, StatusCode: 400, Message: core.BAD_REQUEST_M},
	core.INTERNAL_BAD_REQUEST:         {Code: core.INTERNAL_BAD_REQUEST, StatusCode: 400, Message: core.INTERNAL_BAD_REQUEST_M},
	core.DATA_ALREADY_EXISTS:          {Code: core.DATA_ALREADY_EXISTS, StatusCode: 400, Message: core.DATA_ALREADY_EXISTS_M},
	core.UNAUTHORIZED:                 {Code: core.UNAUTHORIZED, StatusCode: 401, Message: core.UNAUTHORIZED_M},
	core.JWT_TOKEN_REQUIRED:           {Code: core.JWT_TOKEN_REQUIRED, StatusCode: 401, Message: core.JWT_TOKEN_REQUIRED_M},
	core.JWT_REFRESH_TOKEN_REQUIRED:   {Code: core.JWT_REFRESH_TOKEN_REQUIRED, StatusCode: 401, Message: core.JWT_REFRESH_TOKEN_REQUIRED_M},
	core.INVALID_JWT_TOKEN:            {Code: core.INVALID_JWT_TOKEN, StatusCode: 401, Message: core.INVALID_JWT_TOKEN_M},
	core.INVALID_JWT_REFRESH_TOKEN:    {Code: core.INVALID_JWT_REFRESH_TOKEN, StatusCode: 401, Message: core.INVALID_JWT_REFRESH_TOKEN_M},
	core.X_CLIENT_ID_REQUIRED:         {Code: core.X_CLIENT_ID_REQUIRED, StatusCode: 401, Message: core.X_CLIENT_ID_REQUIRED_M},
	core.WRONG_USERNAME:               {Code: core.WRONG_USERNAME, StatusCode: 401, Message: core.WRONG_USERNAME_M},
	core.WRONG_PASSWORD:               {Code: core.WRONG_PASSWORD, StatusCode: 401, Message: core.WRONG_PASSWORD_M},
	core.INTERNAL_UNAUTHORIZED:        {Code: core.INTERNAL_UNAUTHORIZED, StatusCode: 401, Message: core.INTERNAL_UNAUTHORIZED_M},
	core.INTERNAL_FORBIDDEN:           {Code: core.INTERNAL_FORBIDDEN, StatusCode: 401, Message: core.INTERNAL_FORBIDDEN_M},
	core.INTERNAL_NOT_FOUND:           {Code: core.INTERNAL_NOT_FOUND, StatusCode: 401, Message: core.INTERNAL_NOT_FOUND_M},
	core.INTERNAL_CONFLICT:            {Code: core.INTERNAL_CONFLICT, StatusCode: 401, Message: core.INTERNAL_CONFLICT_M},
	core.INTERNAL_SERVICE_UNAVAILABLE: {Code: core.INTERNAL_SERVICE_UNAVAILABLE, StatusCode: 401, Message: core.INTERNAL_SERVICE_UNAVAILABLE_M},
	core.UNPROCESSABLE_ENTITY:         {Code: core.UNPROCESSABLE_ENTITY, StatusCode: 422, Message: core.UNPROCESSABLE_ENTITY_M},
	core.UPGRADE_REQUIRED:             {Code: core.UPGRADE_REQUIRED, StatusCode: 426, Message: core.UPGRADE_REQUIRED_M},
	core.INTERNAL_SERVER_ERROR:        {Code: core.INTERNAL_SERVER_ERROR, StatusCode: 500, Message: core.INTERNAL_SERVER_ERROR_M},
	core.INTERNAL_SERVICE_ERROR:       {Code: core.INTERNAL_SERVICE_ERROR, StatusCode: 500, Message: core.INTERNAL_SERVICE_ERROR_M},
	core.SERVICE_UNAVAILABLE:          {Code: core.SERVICE_UNAVAILABLE, StatusCode: 503, Message: core.SERVICE_UNAVAILABLE_M},
}

type Result struct {
	Code       int         `json:"code"`
	StatusCode int         `json:"status_code,omitempty"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors,omitempty"`
}

type Response struct {
	Result Result      `json:"result"`
	Data   interface{} `json:"data"`
}

type Meta struct {
	CurrentPage int `json:"current_page,omitempty" example:"1"`
	PerPage     int `json:"per_page,omitempty" example:"10"`
	From        int `json:"from,omitempty" example:"1"`
	To          int `json:"to,omitempty" example:"10"`
	Total       int `json:"total,omitempty" example:"100"`
	LastPage    int `json:"last_page,omitempty" example:"10"`
}

type List struct {
	Result Result `json:"result"`
	Data   any    `json:"data"`
	Meta   *Meta  `json:"meta,omitempty"`
}

func NewResponse(result Result, data any) Response {
	return Response{
		Result: result,
		Data:   data,
	}
}

func NewListResponse(result Result, data any) List {
	return List{
		Result: result,
		Data:   data,
	}
}

func NewListResponseWithMeta(result Result, data any, meta Meta) List {
	return List{
		Result: result,
		Data:   data,
		Meta:   &meta,
	}
}

func HandleError(err *core.CustomError) Result {
	if err == nil {
		return Result{
			Code:       http.StatusInternalServerError,
			StatusCode: http.StatusInternalServerError,
			Message:    "internal server error",
			Errors:     err.Errors,
		}
	}

	result, ok := ResultMap[err.Code]
	if !ok {
		result = Result{
			Code:       http.StatusInternalServerError,
			StatusCode: http.StatusInternalServerError,
			Message:    "internal server error",
			Errors:     nil,
		}
	}

	if err.Errors != nil {
		result.Errors = err.Errors
	}
	if err.Message != "" {
		result.Message = err.Message
	}

	return result
}

func HandleBind(err error) Result {
	errors := ErrorBindBodyRequest(err)
	res := ResultMap[core.BAD_REQUEST]
	res.Errors = []string{errors.Field + " " + errors.Format}
	return res
}

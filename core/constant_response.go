package core

const (
	OK                            = 200000
	DATA_NOT_FOUND                = 200001
	DATA_NOT_UPDATED              = 200002
	CREATED                       = 201000
	CREATED_BAD_REQUEST           = 201001
	DELETED                       = 204000
	BAD_REQUEST                   = 400000
	DATA_ALREADY_EXISTS           = 400001
	INTERNAL_BAD_REQUEST          = 400002
	UNAUTHORIZED                  = 401000
	JWT_TOKEN_REQUIRED            = 401001
	JWT_REFRESH_TOKEN_REQUIRED    = 401002
	INVALID_JWT_TOKEN             = 401003
	INVALID_JWT_REFRESH_TOKEN     = 401004
	X_CLIENT_ID_REQUIRED          = 401005
	WRONG_USERNAME                = 401006
	WRONG_PASSWORD                = 401007
	INTERNAL_UNAUTHORIZED         = 401008
	UNPROCESSABLE_ENTITY          = 422000
	INTERNAL_UNPROCESSABLE_ENTITY = 422001
	UPGRADE_REQUIRED              = 426000
	INTERNAL_FORBIDDEN            = 403000
	INTERNAL_NOT_FOUND            = 404000
	INTERNAL_CONFLICT             = 409000
	INTERNAL_SERVER_ERROR         = 500000
	INTERNAL_SERVICE_ERROR        = 500001
	SERVICE_UNAVAILABLE           = 503000
	INTERNAL_SERVICE_UNAVAILABLE  = 503001
)

const (
	OK_M                           = "ok"
	DATA_NOT_FOUND_M               = "data not found"
	DATA_NOT_UPDATED_M             = "data not updated"
	CREATED_M                      = "created"
	CREATED_BAD_REQUEST_M          = "bad request bad created"
	DELETED_M                      = "deleted"
	BAD_REQUEST_M                  = "bad request"
	INTERNAL_BAD_REQUEST_M         = "internal bad request"
	DATA_ALREADY_EXISTS_M          = "data already exists"
	UNAUTHORIZED_M                 = "unauthorized"
	JWT_TOKEN_REQUIRED_M           = "jwt token required"
	JWT_REFRESH_TOKEN_REQUIRED_M   = "jwt refresh token required"
	INVALID_JWT_TOKEN_M            = "invalid jwt token"
	INVALID_JWT_REFRESH_TOKEN_M    = "invalid jwt refresh token"
	X_CLIENT_ID_REQUIRED_M         = "x client id header request required"
	WRONG_USERNAME_M               = "wrong username"
	WRONG_PASSWORD_M               = "wrong password"
	INTERNAL_UNAUTHORIZED_M        = "internal unauthorized"
	INTERNAL_FORBIDDEN_M           = "internal forbidden"
	INTERNAL_NOT_FOUND_M           = "internal not found"
	INTERNAL_CONFLICT_M            = "internal conflict"
	INTERNAL_SERVICE_UNAVAILABLE_M = "internal service unavailable"
	UNPROCESSABLE_ENTITY_M         = "unprocessable entity"
	UPGRADE_REQUIRED_M             = "upgrade required"
	INTERNAL_SERVER_ERROR_M        = "internal server error"
	INTERNAL_SERVICE_ERROR_M       = "internal service error"
	SERVICE_UNAVAILABLE_M          = "service unavailable"
)

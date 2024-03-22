package mocks

import "dot/core"

type CacheRepositoryMock struct {
	Delete  *core.CustomError
	Get     string
	GetErr  *core.CustomError
	HSet    *core.CustomError
	Publish *core.CustomError
	Set     *core.CustomError
}

package mocks

import "dot/core"

type PubsubMock struct {
	Publish   *core.CustomError
	Subscribe *core.CustomError
}

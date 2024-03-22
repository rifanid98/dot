package broker

import "dot/core"

//go:generate mockery --name Pubsub --filename pubsub.go --output ./mocks
type Pubsub interface {
	Publish(ic *core.InternalContext, data []byte) *core.CustomError
	Subscribe(ic *core.InternalContext) *core.CustomError
}

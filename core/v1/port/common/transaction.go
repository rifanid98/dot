package common

import (
	"dot/core"
)

//go:generate mockery --name Transaction --filename transaction.go --output ./mocks
type Transaction interface {
	StartTransaction(ic *core.InternalContext) (Transaction, *core.InternalContext, *core.CustomError)
	CommitTransaction(ic *core.InternalContext) *core.CustomError
	AbortTransaction(ic *core.InternalContext) *core.CustomError
}

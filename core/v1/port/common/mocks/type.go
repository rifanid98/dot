package mocks

import (
	"dot/core"
)

type TransactionMock struct {
	AbortTransaction      *core.CustomError
	CommitTransaction     *core.CustomError
	StartTransactionTx    *Transaction
	StartTransactionTxCtx *core.InternalContext
	StartTransactionErr   *core.CustomError
}

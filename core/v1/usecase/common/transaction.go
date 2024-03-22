package common

import (
	"dot/core"
	"dot/pkg/util"

	portCommon "dot/core/v1/port/common"
)

var log = util.NewLogger()

func AbortTransaction(ic *core.InternalContext, tx portCommon.Transaction, firstErr *core.CustomError) *core.CustomError {
	err := tx.AbortTransaction(ic)
	if err != nil {
		log.Error(ic.ToContext(), "failed to abort transaction", err)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return firstErr
}

func CommitTransaction(ic *core.InternalContext, tx portCommon.Transaction) *core.CustomError {
	err := tx.CommitTransaction(ic)
	if err != nil {
		log.Error(ic.ToContext(), "failed to commit transaction", err)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}
	return nil
}

package mysql

import (
	"dot/core"
	"dot/core/v1/port/common"
	"gorm.io/gorm"
)

type transactionImpl struct {
	*gorm.DB
}

func NewTransaction(db *gorm.DB) *transactionImpl {
	return &transactionImpl{db}
}

func (t *transactionImpl) StartTransaction(ic *core.InternalContext) (common.Transaction, *core.InternalContext, *core.CustomError) {
	db := t.DB.Begin()
	newIc := ic.Clone()
	newIc.InjectData(map[string]any{
		"db": db,
	})

	log.Info(ic.ToContext(), "transaction started")
	return t, newIc, nil
}

func (t *transactionImpl) CommitTransaction(txCtx *core.InternalContext) *core.CustomError {
	db, cerr := t.extractDb(txCtx)
	if cerr != nil {
		return cerr
	}

	err := db.Commit().Error
	if err != nil {
		log.Error(txCtx.ToContext(), "failed to commit transaction", err)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	log.Info(txCtx.ToContext(), "transaction commited")
	return nil
}

func (t *transactionImpl) AbortTransaction(txCtx *core.InternalContext) *core.CustomError {
	db, cerr := t.extractDb(txCtx)
	if cerr != nil {
		return cerr
	}

	err := db.Rollback().Error
	if err != nil {
		log.Error(txCtx.ToContext(), "failed to rollback transaction", err)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	log.Info(txCtx.ToContext(), "transaction aborted")
	return nil
}

func (t *transactionImpl) extractDb(txCtx *core.InternalContext) (*gorm.DB, *core.CustomError) {
	ctxData := txCtx.GetData()
	if ctxData == nil {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to get internal context data; it's empty",
		}
	}

	db := ctxData["db"]
	if db == nil {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to get db from internal context; it's empty",
		}
	}

	s, ok := db.(*gorm.DB)
	if !ok {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to extract transaction db",
		}
	}

	return s, nil
}

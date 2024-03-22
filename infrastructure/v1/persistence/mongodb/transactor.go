package mongodb

import (
	"dot/core"
	"dot/core/v1/port/common"
)

type transactionImpl struct {
	Client
}

func NewTransaction(client Client) *transactionImpl {
	return &transactionImpl{Client: client}
}

func (t *transactionImpl) StartTransaction(ic *core.InternalContext) (common.Transaction, *core.InternalContext, *core.CustomError) {
	session, err := t.Client.StartSession(ic)

	newIc := ic.Clone()
	newIc.InjectData(map[string]any{
		"client":  session.Client(),
		"session": session,
	})

	err = session.StartTransaction()
	if err != nil {
		log.Error(ic.ToContext(), "failed to start transaction")
		return nil, nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	log.Info(ic.ToContext(), "transaction started")
	return t, newIc, nil
}

func (t *transactionImpl) AbortTransaction(txCtx *core.InternalContext) *core.CustomError {
	session, cerr := t.extractSession(txCtx)
	if cerr != nil {
		return cerr
	}

	err := session.AbortTransaction(txCtx.ToContext())
	if err != nil {
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	log.Info(txCtx.ToContext(), "transaction aborted")
	return nil
}

func (t *transactionImpl) CommitTransaction(txCtx *core.InternalContext) *core.CustomError {
	session, cerr := t.extractSession(txCtx)
	if cerr != nil {
		return cerr
	}

	err := session.CommitTransaction(txCtx.ToContext())
	if err != nil {
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	log.Info(txCtx.ToContext(), "transaction commited")
	return nil
}

func (t *transactionImpl) extractSession(txCtx *core.InternalContext) (Session, *core.CustomError) {
	ctxData := txCtx.GetData()
	if ctxData == nil {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to get internal context data; it's empty",
		}
	}

	session := ctxData["session"]
	if session == nil {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to get session from internal context; it's empty",
		}
	}

	s, ok := session.(Session)
	if !ok {
		return nil, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: "failed to extract transaction session",
		}
	}

	return s, nil
}

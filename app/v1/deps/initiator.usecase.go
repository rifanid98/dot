package deps

import (
	portCommon "dot/core/v1/port/common"
	"dot/core/v1/usecase/auth"
	"dot/core/v1/usecase/book"
	"dot/core/v1/usecase/subscriber"
)

func (d *dependency) initService() {
	d.initAuthUsecase()
	d.initBookUsecase()
	d.initSubscriberUsecase()
}

func (d *dependency) initAuthUsecase() {
	d.usecase.AuthUsecase = auth.NewAuthUsecase(
		d.repo.AccountRepository,
		d.repo.CacheRepository,
		d.base.Cfg,
	)
}

func (d *dependency) initBookUsecase() {
	var tx portCommon.Transaction
	switch d.base.Cfg.Db {
	case DB_MONGO_DB:
		tx = d.base.Mdbt
	case DB_MYSQL:
		tx = d.base.Mysqlt
	}
	d.usecase.BookUsecase = book.NewBookUsecase(
		d.repo.BookRepository,
		d.repo.CacheRepository,
		tx,
	)
}

func (d *dependency) initSubscriberUsecase() {
	d.usecase.SubscriberUsecase = subscriber.NewSubscriberUsecase()
}

package deps

import (
	mdb "dot/infrastructure/v1/persistence/mongodb/repository"
)

func (d *dependency) initMongodbRepository() {
	d.initMongodbAccountRepository()
	d.initMongodbBookRepository()
}

func (d *dependency) initMongodbAccountRepository() {
	d.repo.AccountRepository = mdb.NewAccountRepository(d.base.Mdb, d.base.Cfg)
}

func (d *dependency) initMongodbBookRepository() {
	d.repo.BookRepository = mdb.NewBookRepository(d.base.Mdb, d.base.Cfg)
}

package deps

import (
	mysql "dot/infrastructure/v1/persistence/mysql/repository"
)

func (d *dependency) initMysqlRepository() {
	d.initMysqlAccountRepository()
	d.initMysqlBookRepository()
}

func (d *dependency) initMysqlAccountRepository() {
	d.repo.AccountRepository = mysql.NewAccountRepository(d.base.Mysql, d.base.Cfg)
}

func (d *dependency) initMysqlBookRepository() {
	d.repo.BookRepository = mysql.NewBookRepository(d.base.Mysql, d.base.Cfg)
}

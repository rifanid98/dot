package deps

import (
	rdb "dot/infrastructure/v1/persistence/redisdb/repository"
)

func (d *dependency) initRepository() {
	switch d.base.Cfg.Db {
	case DB_MYSQL:
		d.initMysqlRepository()
	case DB_MONGO_DB:
		d.initMongodbRepository()
	default:
		panic("DATABASE_USE not initialized")
	}
	d.initCacheRepository()
}

func (d *dependency) initCacheRepository() {
	d.repo.CacheRepository = rdb.NewCacheRepository(d.base.Rdbc)
}

package deps

import (
	"dot/config"
	"dot/infrastructure/v1/broker/gcp"
	"dot/infrastructure/v1/persistence/mongodb"
	"dot/infrastructure/v1/persistence/mysql"
	"dot/infrastructure/v1/persistence/redisdb"
	"dot/pkg/api"
	"dot/pkg/util"
)

func (d *dependency) initBase() {
	d.initConfig()
	d.initMongodb()
	d.initMysql()
	d.initCache()
	d.initHttpClient()
	d.initRetrier()

	// [NEED GCP CREDENTIALS ACCESS]
	//d.initGcpBroker()
}

func (d *dependency) initConfig() {
	d.base.Cfg = config.GetConfig()
}

func (d *dependency) initMongodb() {
	mdb, mdbc, err := mongodb.New(d.base.Cfg.MongoDb)
	if err != nil {
		panic(err)
	}
	tx := mongodb.NewTransaction(mdbc)
	d.base.Mdb = mdb
	d.base.Mdbc = mdbc
	d.base.Mdbt = tx
}

func (d *dependency) initMysql() {
	db := mysql.New(d.base.Cfg.MySQL)
	tx := mysql.NewTransaction(db)
	d.base.Mysql = db
	d.base.Mysqlt = tx
}

func (d *dependency) initCache() {
	rdb, err := redisdb.New(d.base.Cfg.Redis)
	if err != nil {
		panic(err)
	}
	d.base.Rdbc = rdb
}

func (d *dependency) initHttpClient() {
	d.base.Httpc = api.NewHttpClient()
}

func (d *dependency) initGcpBroker() {
	gcpc, err := gcp.New(d.base.Cfg.GcpPubsub)
	if err != nil {
		panic(err)
	}
	d.base.Gcpc = gcpc
}

func (d *dependency) initRetrier() {
	d.base.Rtr = util.NewRetrier()
}

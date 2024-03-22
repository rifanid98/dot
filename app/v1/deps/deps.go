package deps

import (
	"dot/config"
	"dot/core/v1/port/account"
	"dot/core/v1/port/adapter"
	"dot/core/v1/port/auth"
	"dot/core/v1/port/book"
	"dot/core/v1/port/broker"
	"dot/core/v1/port/cache"
	"dot/core/v1/port/common"
	"dot/core/v1/port/retrier"
	"dot/core/v1/port/scheduler"
	"dot/core/v1/port/subscriber"
	"dot/infrastructure/v1/broker/gcp"
	"dot/infrastructure/v1/persistence/mongodb"
	"dot/infrastructure/v1/persistence/redisdb"
	"dot/pkg/api"
	"gorm.io/gorm"
)

const (
	DB_MONGO_DB = 1
	DB_MYSQL    = 2
)

type base struct {
	Cfg    *config.AppConfig
	Mdb    mongodb.Database
	Mdbc   mongodb.Client
	Mdbt   common.Transaction
	Mysql  *gorm.DB
	Mysqlt common.Transaction
	Rdbc   redisdb.Client
	Httpc  api.HttpDoer
	Rtr    retrier.Retrier
	Schlr  scheduler.Scheduler
	Gcpc   gcp.Client
}

type repository struct {
	account.AccountRepository
	book.BookRepository
	cache.CacheRepository
}

type apicall struct {
	adapter.XenditApiCall
}

type usecase struct {
	auth.AuthUsecase
	book.BookUsecase
	subscriber.SubscriberUsecase
}

type msgbroker struct {
	broker.Pubsub
}

type dependency struct {
	base    *base
	repo    *repository
	apicall *apicall
	usecase *usecase
	broker  *msgbroker
}

type IDependency interface {
	GetServices() *usecase
	GetRepositories() *repository
	GetBase() *base
	GetBroker() *msgbroker
}

func BuildDependency() *dependency {
	dep := &dependency{
		base:    &base{},
		repo:    &repository{},
		apicall: &apicall{},
		usecase: &usecase{},
		broker:  &msgbroker{},
	}
	dep.initBase()       // execute first
	dep.initRepository() // execute second
	dep.initApiCall()    // execute third
	dep.initBroker()     // execute fourth
	dep.initService()    // execute fifth
	dep.initScheduler()  // execute sixth
	return dep
}

func (d *dependency) GetBase() *base {
	return d.base
}

func (d *dependency) GetServices() *usecase {
	return d.usecase
}

func (d *dependency) GetRepositories() *repository {
	return d.repo
}

func (d *dependency) GetBroker() *msgbroker {
	return d.broker
}

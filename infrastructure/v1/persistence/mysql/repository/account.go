package repository

import (
	"dot/config"
	"dot/core"
	"dot/core/v1/entity"
	"dot/infrastructure/v1/persistence/mysql/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type accountRepositoryImpl struct {
	db    *gorm.DB
	cfg   *config.AppConfig
	table string
}

func NewAccountRepository(db *gorm.DB, cfg *config.AppConfig) *accountRepositoryImpl {
	return &accountRepositoryImpl{
		db:    db,
		cfg:   cfg,
		table: "account",
	}
}

func (a accountRepositoryImpl) InsertAccount(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError) {
	account.Id = uuid.NewString()
	now := time.Now()
	account.Created = &now
	account.Modified = &now
	accountDoc := new(model.Account).Bind(account)

	res := db(ic, a.db).WithContext(ic.ToContext()).Table(a.table).Create(accountDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed InsertAccount", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, &core.CustomError{
			Code: core.UNPROCESSABLE_ENTITY,
		}
	}

	return account, nil
}

func (a *accountRepositoryImpl) FindAccountByEmail(ic *core.InternalContext, email string) (*entity.Account, *core.CustomError) {
	accountDoc := new(model.Account)

	res := db(ic, a.db).WithContext(ic.ToContext()).
		Table(a.table).
		Limit(1).
		Where("email = ?", email).
		Find(accountDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed FindAccountByEmail", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, nil
	}

	return accountDoc.Entity()
}

func (a *accountRepositoryImpl) FindAccountById(ic *core.InternalContext, id string) (*entity.Account, *core.CustomError) {
	accountDoc := new(model.Account)

	res := db(ic, a.db).WithContext(ic.ToContext()).
		Table(a.table).
		Limit(1).
		Where("id = ?", id).
		Find(accountDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed FindAccountById", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, nil
	}

	return accountDoc.Entity()
}

func (a *accountRepositoryImpl) UpdateAccount(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError) {
	accountDoc := new(model.Account).Bind(account)

	res := db(ic, a.db).WithContext(ic.ToContext()).
		Table(a.table).
		Where("id = ?", account.Id).
		Updates(accountDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed UpdateAccount()", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, nil
	}

	return accountDoc.Entity()
}

func (a *accountRepositoryImpl) DeleteAccount(ic *core.InternalContext, id string) *core.CustomError {
	res := db(ic, a.db).WithContext(ic.ToContext()).
		Table(a.table).
		Where("id = ?", id).
		Delete(new(model.Account))
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed DeleteAccount()", res.Error)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "failed delete account",
		}
	}

	return nil
}

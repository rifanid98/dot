package model

import (
	"dot/core"
	"dot/pkg/helper"
	"time"

	"dot/core/v1/entity"
)

type Account struct {
	Id               string     `gorm:"id,primaryKey"`
	Email            string     `gorm:"email,omitempty"`
	PhoneNumber      string     `gorm:"phone_number,omitempty"`
	Username         string     `gorm:"username,omitempty"`
	Password         string     `gorm:"password,omitempty"`
	Otp              string     `gorm:"otp,omitempty"`
	ResetToken       string     `gorm:"reset_token"`
	ResetTokenExpire *time.Time `gorm:"reset_token_expire"`
	SentAccess       *time.Time `gorm:"sent_access,omitempty"`
	Age              int        `gorm:"age"`
	Gender           int        `gorm:"gender"`
	Metadata         string     `gorm:"metadata"`
	Verified         int        `gorm:"verified"`
	VerifiedDate     *time.Time `gorm:"verified_date"`
	Created          *time.Time `gorm:"created,omitempty"`
	Modified         *time.Time `gorm:"modified,omitempty"`

	Books []Book `gorm:"foreignKey:author"`
}

func (m *Account) TableName() string {
	return "account"
}

func (doc *Account) Bind(account *entity.Account) *Account {
	return &Account{
		Id:               account.Id,
		Email:            account.Email,
		PhoneNumber:      account.PhoneNumber,
		Username:         account.Username,
		Password:         account.Password,
		Otp:              account.Otp,
		ResetToken:       account.ResetToken,
		ResetTokenExpire: account.ResetTokenExpire,
		SentAccess:       account.SentAccess,
		Age:              account.Age,
		Gender:           account.Gender,
		Metadata:         helper.DataToString(account.Metadata),
		Verified:         account.Verified,
		VerifiedDate:     account.VerifiedDate,
		Created:          account.Created,
		Modified:         account.Modified,
	}
}

func (doc *Account) Entity() (*entity.Account, *core.CustomError) {
	metadata := make(map[string]any)
	if doc.Metadata != "" {
		cerr := helper.StringToStruct(doc.Metadata, &metadata)
		if cerr != nil {
			return nil, cerr
		}
	}

	return &entity.Account{
		Id:               doc.Id,
		Otp:              doc.Otp,
		Username:         doc.Username,
		Email:            doc.Email,
		Password:         doc.Password,
		PhoneNumber:      doc.PhoneNumber,
		Age:              doc.Age,
		Gender:           doc.Gender,
		Verified:         doc.Verified,
		Metadata:         metadata,
		VerifiedDate:     doc.VerifiedDate,
		ResetToken:       doc.ResetToken,
		ResetTokenExpire: doc.ResetTokenExpire,
		SentAccess:       doc.SentAccess,
		Created:          doc.Created,
		Modified:         doc.Modified,
	}, nil
}

type Accounts []Account

func (accs Accounts) Bind(accounts []entity.Account) Accounts {
	for i := range accounts {
		now := time.Now()
		acc := new(Account).Bind(&accounts[i])
		acc.Created = &now
		acc.Modified = &now
		if accounts[i].Id != "" {
			acc.Id = accounts[i].Id
		}
		accs = append(accs, *acc)
	}
	return accs
}

func (accs Accounts) Generics() []any {
	var data []any
	for i := range accs {
		data = append(data, accs[i])
	}
	return data
}

func (accs Accounts) Entities() ([]entity.Account, *core.CustomError) {
	var es []entity.Account
	for i := range accs {
		e, cerr := accs[i].Entity()
		if cerr != nil {
			return nil, cerr
		}
		es = append(es, *e)
	}
	return es, nil
}

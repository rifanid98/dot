package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"dot/core/v1/entity"
)

type Account struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	Email            string             `bson:"email,omitempty"`
	PhoneNumber      string             `bson:"phone_number,omitempty"`
	Username         string             `bson:"username,omitempty"`
	Password         string             `bson:"password,omitempty"`
	Otp              string             `bson:"otp,omitempty"`
	ResetToken       string             `bson:"reset_token"`
	ResetTokenExpire *time.Time         `bson:"reset_token_expire"`
	SentAccess       *time.Time         `bson:"sent_access,omitempty"`
	Age              int                `bson:"age"`
	Gender           int                `bson:"gender"`
	Metadata         map[string]any     `bson:"metadata"`
	Verified         int                `bson:"verified"`
	VerifiedDate     *time.Time         `bson:"verified_date"`
	Created          *time.Time         `bson:"created,omitempty"`
	Modified         *time.Time         `bson:"modified,omitempty"`
}

func (doc *Account) Bind(account *entity.Account) *Account {
	return &Account{
		Id:               GetObjectId(account.Id),
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
		Metadata:         account.Metadata,
		Verified:         account.Verified,
		VerifiedDate:     account.VerifiedDate,
		Created:          account.Created,
		Modified:         account.Modified,
	}
}

func (doc *Account) Entity() *entity.Account {
	return &entity.Account{
		Id:               GetObjectIdHex(doc.Id),
		Otp:              doc.Otp,
		Username:         doc.Username,
		Email:            doc.Email,
		Password:         doc.Password,
		PhoneNumber:      doc.PhoneNumber,
		Age:              doc.Age,
		Gender:           doc.Gender,
		Verified:         doc.Verified,
		Metadata:         doc.Metadata,
		VerifiedDate:     doc.VerifiedDate,
		ResetToken:       doc.ResetToken,
		ResetTokenExpire: doc.ResetTokenExpire,
		SentAccess:       doc.SentAccess,
		Created:          doc.Created,
		Modified:         doc.Modified,
	}
}

type Accounts []Account

func (accs Accounts) Bind(accounts []entity.Account) Accounts {
	for i := range accounts {
		acc := new(Account).Bind(&accounts[i])
		if accounts[i].Id != "" {
			acc.Id = GetObjectId(accounts[i].Id)
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

func (accs Accounts) Entities() []entity.Account {
	var es []entity.Account
	for i := range accs {
		es = append(es, *accs[i].Entity())
	}
	return es
}

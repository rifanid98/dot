package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"dot/config"
	"dot/core"
	"dot/core/v1/entity"
	"dot/infrastructure/v1/persistence/mongodb"
	"dot/infrastructure/v1/persistence/mongodb/model"
)

type accountRepositoryImpl struct {
	collection mongodb.Collection
	cfg        *config.AppConfig
}

func NewAccountRepository(db mongodb.Database, cfg *config.AppConfig) *accountRepositoryImpl {
	return &accountRepositoryImpl{
		collection: db.Collection("account"),
		cfg:        cfg,
	}
}

func (r *accountRepositoryImpl) InsertAccount(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError) {
	doc := new(model.Account).Bind(account)
	now := time.Now()
	doc.Created = &now
	doc.Modified = &now

	res, err := r.collection.InsertOne(ic.ToContext(), &doc)
	if err != nil {
		log.Error(ic.ToContext(), "failed to InsertAccount", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	doc.Id = res.InsertedID.(primitive.ObjectID)
	return doc.Entity(), nil
}

func (r *accountRepositoryImpl) FindAccountByEmail(ic *core.InternalContext, email string) (*entity.Account, *core.CustomError) {
	var data model.Account

	filter := bson.M{
		"email": email,
	}

	err := r.collection.FindOne(ic.ToContext(), filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		log.Error(ic.ToContext(), "failed to FindAccountByEmail", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}
	return data.Entity(), nil
}

func (r *accountRepositoryImpl) FindAccountById(ic *core.InternalContext, id string) (*entity.Account, *core.CustomError) {
	var doc model.Account

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "invalid account id",
		}
	}

	filter := bson.M{
		"_id": objId,
	}

	err = r.collection.FindOne(ic.ToContext(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		log.Error(ic.ToContext(), "failed FindAccountById : %v", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return doc.Entity(), nil
}

func (r *accountRepositoryImpl) UpdateAccount(ic *core.InternalContext, account *entity.Account) (*entity.Account, *core.CustomError) {
	doc := new(model.Account).Bind(account)
	now := time.Now()
	doc.Modified = &now

	objId, err := primitive.ObjectIDFromHex(account.Id)
	if err != nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "invalid account id",
		}
	}

	filter := bson.M{"_id": objId}
	set := bson.M{"$set": doc}
	_, err = r.collection.UpdateOne(ctx(ic), filter, set)
	if err != nil {
		log.Error(ic.ToContext(), "failed UpdateAccount : %v", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return doc.Entity(), nil
}

func (r *accountRepositoryImpl) DeleteAccount(ic *core.InternalContext, id string) *core.CustomError {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "invalid account id",
		}
	}

	filter := bson.M{"_id": objId}
	_, err = r.collection.DeleteOne(ctx(ic), filter)
	if err != nil {
		log.Error(ic.ToContext(), "failed DeleteAccount : %v", err)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return nil
}

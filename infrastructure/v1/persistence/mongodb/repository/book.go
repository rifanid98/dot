package repository

import (
	"dot/pkg/helper"
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

type bookRepositoryImpl struct {
	collection mongodb.Collection
	cfg        *config.AppConfig
}

func NewBookRepository(db mongodb.Database, cfg *config.AppConfig) *bookRepositoryImpl {
	return &bookRepositoryImpl{
		collection: db.Collection("book"),
		cfg:        cfg,
	}
}

func (r *bookRepositoryImpl) InsertBook(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError) {
	doc := new(model.Book).Bind(book)
	now := time.Now()
	doc.Created = &now
	doc.Modified = &now

	res, err := r.collection.InsertOne(ctx(ic), &doc)
	if err != nil {
		log.Error(ic.ToContext(), "failed to InsertBook", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	doc.Id = res.InsertedID.(primitive.ObjectID)
	return doc.Entity(), nil
}

func (r *bookRepositoryImpl) FindBookByName(ic *core.InternalContext, name string) (*entity.Book, *core.CustomError) {
	var data model.Book

	filter := bson.M{
		"name": name,
	}

	err := r.collection.FindOne(ctx(ic), filter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		log.Error(ic.ToContext(), "failed to FindBookByUsername", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}
	return data.Entity(), nil
}

func (r *bookRepositoryImpl) FindBookById(ic *core.InternalContext, id string) (*entity.Book, *core.CustomError) {
	var doc model.Book

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "invalid book id",
		}
	}

	filter := bson.M{
		"_id": objId,
	}

	err = r.collection.FindOne(ctx(ic), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		log.Error(ic.ToContext(), "failed FindBookById : %v", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return doc.Entity(), nil
}

func (r *bookRepositoryImpl) UpdateBook(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError) {
	doc := new(model.Book).Bind(book)
	now := time.Now()
	doc.Modified = &now

	objId, err := primitive.ObjectIDFromHex(book.Id)
	if err != nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "invalid book id",
		}
	}

	filter := bson.M{"_id": objId}
	set := bson.M{"$set": doc}
	_, err = r.collection.UpdateOne(ctx(ic), filter, set)
	if err != nil {
		log.Error(ic.ToContext(), "failed UpdateBook : %v", err)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return doc.Entity(), nil
}

func (r *bookRepositoryImpl) DeleteBook(ic *core.InternalContext, id string) *core.CustomError {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "invalid book id",
		}
	}

	filter := bson.M{"_id": objId}
	_, err = r.collection.DeleteOne(ctx(ic), filter)
	if err != nil {
		log.Error(ic.ToContext(), "failed DeleteBook : %v", err)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return nil
}

func (r *bookRepositoryImpl) FindBooks(ic *core.InternalContext, meta map[string]any) ([]entity.Book, int64, *core.CustomError) {
	page := helper.DataToInt(meta["page"])
	limit := helper.DataToInt(meta["limit"])
	offset := limit * (page - 1)

	facet := bson.D{
		{Key: "$facet", Value: bson.D{
			{Key: "metadata", Value: []bson.M{
				{"$count": "total"},
			}},
			{Key: "data", Value: []bson.M{
				{"$skip": offset},
				{"$limit": limit},
			}},
		}},
	}

	pipeline := mongo.Pipeline{facet}
	res, err := r.collection.Aggregate(ctx(ic), pipeline)
	if err != nil {
		log.Error(ic.ToContext(), "failed r.collection.Aggregate(ic.ToContext(), pipeline)", err)
		return nil, 0, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	var doc []struct {
		Metadata []map[string]interface{} `bson:"metadata"`
		Data     model.Books              `bson:"data"`
	}

	err = res.All(ic.ToContext(), &doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, 0, nil
		}

		log.Error(ic.ToContext(), "failed res.All(ic.ToContext(), &doc)", err)
		return nil, 0, &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	if len(doc[0].Data) < 1 {
		return nil, 0, nil
	}

	total := int64(doc[0].Metadata[0]["total"].(int32))
	return doc[0].Data.Entities(), total, nil
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"dot/core/v1/entity"
)

type Book struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Author   string             `bson:"author"`
	Created  *time.Time         `bson:"created"`
	Modified *time.Time         `bson:"modified"`
}

func (doc *Book) Bind(book *entity.Book) *Book {
	return &Book{
		Id:       GetObjectId(book.Id),
		Name:     book.Name,
		Author:   book.Author,
		Created:  book.Created,
		Modified: book.Modified,
	}
}

func (doc *Book) Entity() *entity.Book {
	return &entity.Book{
		Id:       GetObjectIdHex(doc.Id),
		Name:     doc.Name,
		Author:   doc.Author,
		Created:  doc.Created,
		Modified: doc.Modified,
	}
}

type Books []Book

func (accs Books) Bind(books []entity.Book) Books {
	for i := range books {
		acc := new(Book).Bind(&books[i])
		if books[i].Id != "" {
			acc.Id = GetObjectId(books[i].Id)
		}
		accs = append(accs, *acc)
	}
	return accs
}

func (accs Books) Generics() []any {
	var data []any
	for i := range accs {
		data = append(data, accs[i])
	}
	return data
}

func (accs Books) Entities() []entity.Book {
	var es []entity.Book
	for i := range accs {
		es = append(es, *accs[i].Entity())
	}
	return es
}

package model

import (
	"time"

	"dot/core/v1/entity"
)

type Book struct {
	Id       string     `gorm:"id,primaryKey"`
	Name     string     `gorm:"name"`
	Author   string     `gorm:"size:191"`
	Created  *time.Time `gorm:"created"`
	Modified *time.Time `gorm:"modified"`
}

func (m *Book) TableName() string {
	return "book"
}

func (doc *Book) Bind(book *entity.Book) *Book {
	return &Book{
		Id:       book.Id,
		Name:     book.Name,
		Author:   book.Author,
		Created:  book.Created,
		Modified: book.Modified,
	}
}

func (doc *Book) Entity() *entity.Book {
	return &entity.Book{
		Id:       doc.Id,
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
			acc.Id = books[i].Id
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

package book

import (
	"dot/core"
	"dot/core/v1/entity"
)

//go:generate mockery --name BookRepository --filename book_repository.go --output ./mocks
type BookRepository interface {
	InsertBook(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError)
	FindBooks(ic *core.InternalContext, meta map[string]any) ([]entity.Book, int64, *core.CustomError)
	FindBookById(ic *core.InternalContext, id string) (*entity.Book, *core.CustomError)
	FindBookByName(ic *core.InternalContext, name string) (*entity.Book, *core.CustomError)
	UpdateBook(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError)
	DeleteBook(ic *core.InternalContext, id string) *core.CustomError
}

package book

import (
	"dot/core"
	"dot/core/v1/entity"
)

//go:generate mockery --name BookUsecase --filename book_usecase.go --output ./mocks
type BookUsecase interface {
	BookCreate(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError)
	BookList(ic *core.InternalContext, meta map[string]any) ([]entity.Book, int64, *core.CustomError)
	BookGet(ic *core.InternalContext, id string) (*entity.Book, *core.CustomError)
	BookUpdate(ic *core.InternalContext, book *entity.Book) *core.CustomError
	BookDelete(ic *core.InternalContext, id string) *core.CustomError
}

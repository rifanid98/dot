package repository

import (
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"

	"dot/config"
	"dot/core"
	"dot/core/v1/entity"
	"dot/infrastructure/v1/persistence/mysql/model"
	"dot/pkg/helper"
)

type bookRepositoryImpl struct {
	db    *gorm.DB
	cfg   *config.AppConfig
	table string
}

func NewBookRepository(db *gorm.DB, cfg *config.AppConfig) *bookRepositoryImpl {
	return &bookRepositoryImpl{
		db:    db,
		cfg:   cfg,
		table: "book",
	}
}

func (b *bookRepositoryImpl) InsertBook(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError) {
	book.Id = uuid.NewString()
	now := time.Now()
	book.Created = &now
	book.Modified = &now
	bookDoc := new(model.Book).Bind(book)

	res := db(ic, b.db).WithContext(ic.ToContext()).Table(b.table).Create(bookDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed res := a.db.WithContext(ic.ToContext()).Create(bookDoc)", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, &core.CustomError{
			Code: core.UNPROCESSABLE_ENTITY,
		}
	}

	return book, nil
}

func (b *bookRepositoryImpl) FindBooks(ic *core.InternalContext, meta map[string]any) ([]entity.Book, int64, *core.CustomError) {
	page := helper.DataToInt(meta["page"])
	limit := helper.DataToInt(meta["limit"])
	offset := limit * (page - 1)

	var bookDocs model.Books

	query := db(ic, b.db).WithContext(ic.ToContext()).
		Table("book")

	var total int64
	count := query.Count(&total)
	if count.Error != nil {
		log.Error(ic.ToContext(), "failed query.Count(&count)", count.Error)
		return nil, 0, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	res := query.
		Offset(int(offset)).
		Limit(int(limit)).
		Order("created DESC").
		Find(&bookDocs)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed FindBooks()", res.Error)
		return nil, 0, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, 0, nil
	}

	return bookDocs.Entities(), total, nil
}

func (b *bookRepositoryImpl) FindBookById(ic *core.InternalContext, id string) (*entity.Book, *core.CustomError) {
	bookDoc := new(model.Book)

	res := db(ic, b.db).WithContext(ic.ToContext()).
		Table("book").
		Limit(1).
		Where("id = ?", id).
		Find(bookDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed FindBookById()", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, nil
	}

	return bookDoc.Entity(), nil
}

func (b *bookRepositoryImpl) FindBookByName(ic *core.InternalContext, name string) (*entity.Book, *core.CustomError) {
	bookDoc := new(model.Book)

	res := db(ic, b.db).WithContext(ic.ToContext()).
		Table(b.table).
		Limit(1).
		Where("name = ?", name).
		Find(bookDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed FindBookByName()", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, nil
	}

	return bookDoc.Entity(), nil
}

func (b *bookRepositoryImpl) UpdateBook(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError) {
	bookDoc := new(model.Book).Bind(book)

	res := db(ic, b.db).WithContext(ic.ToContext()).
		Table(b.table).
		Where("id = ?", book.Id).
		Updates(bookDoc)
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed UpdateBook()", res.Error)
		return nil, &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	if res.RowsAffected < 1 {
		return nil, nil
	}

	return bookDoc.Entity(), nil
}

func (b *bookRepositoryImpl) DeleteBook(ic *core.InternalContext, id string) *core.CustomError {
	res := db(ic, b.db).WithContext(ic.ToContext()).
		Table(b.table).
		Where("id = ?", id).
		Delete(new(model.Book))
	if res.Error != nil {
		log.Error(ic.ToContext(), "failed DeleteBook()", res.Error)
		return &core.CustomError{
			Code: core.INTERNAL_SERVER_ERROR,
		}
	}

	return nil
}

package book

import (
	"dot/core"
	"dot/core/v1/entity"
	"dot/pkg/helper"
	"dot/pkg/util"
	"fmt"
	"github.com/r3labs/diff/v3"
	"golang.org/x/exp/slices"
	"math"
	"time"

	portBook "dot/core/v1/port/book"
	portCache "dot/core/v1/port/cache"
	portCommon "dot/core/v1/port/common"
)

var log = util.NewLogger()

type bookUsecaseImpl struct {
	bookRepository  portBook.BookRepository
	cacheRepository portCache.CacheRepository
	transaction     portCommon.Transaction
}

func NewBookUsecase(
	bookRepository portBook.BookRepository,
	cacheRepository portCache.CacheRepository,
	transaction portCommon.Transaction,
) portBook.BookUsecase {
	return &bookUsecaseImpl{
		bookRepository:  bookRepository,
		cacheRepository: cacheRepository,
		transaction:     transaction,
	}
}

func (uc *bookUsecaseImpl) BookCreate(ic *core.InternalContext, book *entity.Book) (*entity.Book, *core.CustomError) {
	_book, cerr := uc.bookRepository.FindBookByName(ic, book.Name)
	if cerr != nil {
		return nil, cerr
	}
	if _book != nil {
		if _book.Name == book.Name {
			return nil, &core.CustomError{
				Code:    core.UNPROCESSABLE_ENTITY,
				Message: "book already exists",
			}
		}
	}

	_book, cerr = uc.bookRepository.InsertBook(ic, book)
	if cerr != nil {
		return nil, cerr
	}

	////[START WITH TRANSACTION]
	//tx, txCtx, cerr := uc.transaction.StartTransaction(ic)
	//
	//_book, cerr = uc.bookRepository.InsertBook(txCtx, book)
	//if cerr != nil {
	//	return nil, cerr
	//}
	//
	//cerr = common.AbortTransaction(txCtx, tx, cerr)
	//if cerr != nil {
	//	return nil, cerr
	//}
	//cerr = common.CommitTransaction(txCtx, tx)
	//if cerr != nil {
	//	return nil, cerr
	//}
	////[END WITH TRANSACTION]

	cerr = uc.redisClearBooks(ic)
	if cerr != nil {
		return nil, cerr
	}

	return _book, cerr
}

func (uc *bookUsecaseImpl) BookList(ic *core.InternalContext, meta map[string]any) ([]entity.Book, int64, *core.CustomError) {
	if helper.DataToString(meta["search"]) != "" || helper.DataToInt(meta["limit"]) != 10 {
		return uc.bookRepository.FindBooks(ic, map[string]any{
			"search": helper.DataToString(meta["search"]),
			"page":   helper.DataToInt(meta["page"]),
			"limit":  helper.DataToInt(meta["limit"]),
		})
	}

	// get total books
	redisKeyBooksTotal := "books::total"
	getTotal, cerr := uc.cacheRepository.Get(ic, redisKeyBooksTotal)
	if cerr != nil {
		return nil, 0, cerr
	}
	total := helper.DataToInt(getTotal)

	// get books list by page
	page := helper.DataToInt(meta["page"])
	redisKeyBooksList := fmt.Sprintf("books::page::%d", page)
	getBooks, cerr := uc.cacheRepository.Get(ic, redisKeyBooksList)
	if cerr != nil {
		return nil, 0, cerr
	}

	var books []entity.Book

	if getBooks != "" {
		cerr = helper.StringToStruct(getBooks, &books)
		if cerr != nil {
			return nil, 0, cerr
		}

		if books != nil {
			log.Info(ic.ToContext(), "books from redis")
			return books, total, nil
		}
	}

	// generate new redis data when redis is empty
	books, total, cerr = uc.bookRepository.FindBooks(ic, map[string]any{
		"page":  helper.DataToInt(meta["page"]),
		"limit": helper.DataToInt(meta["limit"]),
	})

	// set total books
	duration := time.Hour * 24
	cerr = uc.cacheRepository.Set(ic, redisKeyBooksTotal, helper.DataToString(total), &duration)
	if cerr != nil {
		log.Error(ic.ToContext(), "failed uc.cacheRepository.Set(ic, redisKeyBooksTotal, helper.DataToString(total), &duration)", cerr.Error())
		return nil, 0, cerr
	}

	// set books list
	cerr = uc.cacheRepository.Set(ic, redisKeyBooksList, helper.DataToString(books), &duration)
	if cerr != nil {
		log.Error(ic.ToContext(), "failed uc.cacheRepository.Set(ic, redisKeyBooksList, helper.DataToString(books), &duration)", cerr.Error())
		return nil, 0, cerr
	}

	return books, total, cerr
}

func (uc *bookUsecaseImpl) BookGet(ic *core.InternalContext, id string) (*entity.Book, *core.CustomError) {
	return uc.findBookById(ic, id)
}

func (uc *bookUsecaseImpl) BookUpdate(ic *core.InternalContext, book *entity.Book) *core.CustomError {
	_book, cerr := uc.findBookById(ic, book.Id)
	if cerr != nil {
		return cerr
	}

	changelog, err := diff.Diff(_book, book)
	if err != nil {
		return &core.CustomError{
			Code:    core.INTERNAL_SERVER_ERROR,
			Message: err.Error(),
		}
	}

	data := make(map[string]any)
	for _, change := range changelog {
		if !slices.Contains([]string{"Modified", "Created"}, change.Path[0]) && change.To != "" {
			data[helper.GetTagNamebyName(change.Path[0], _book)] = change.To
		}
	}

	cerr = helper.StringToStruct(helper.DataToString(data), _book)
	if cerr != nil {
		return cerr
	}

	book, cerr = uc.bookRepository.UpdateBook(ic, _book)
	if cerr != nil {
		return cerr
	}

	cerr = uc.redisClearBooks(ic)
	if cerr != nil {
		return cerr
	}

	return uc.cacheRepository.Delete(ic, fmt.Sprintf("books::%s", book.Id))
}

func (uc *bookUsecaseImpl) BookDelete(ic *core.InternalContext, id string) *core.CustomError {
	_, cerr := uc.findBookById(ic, id)
	if cerr != nil {
		return cerr
	}

	cerr = uc.bookRepository.DeleteBook(ic, id)
	if cerr != nil {
		return cerr
	}

	redisKeyBook := fmt.Sprintf("books::%s", id)
	cerr = uc.cacheRepository.Delete(ic, redisKeyBook)
	if cerr != nil {
		return cerr
	}

	return uc.redisClearBooks(ic)
}

func (uc *bookUsecaseImpl) findBookById(ic *core.InternalContext, id string) (*entity.Book, *core.CustomError) {
	redisKeyBook := fmt.Sprintf("books::%s", id)
	getBook, cerr := uc.cacheRepository.Get(ic, redisKeyBook)
	if cerr != nil {
		return nil, cerr
	}

	if getBook != "" {
		var book entity.Book
		cerr = helper.StringToStruct(getBook, &book)
		if cerr != nil {
			return nil, cerr
		}
		return &book, nil
	}

	book, cerr := uc.bookRepository.FindBookById(ic, id)
	if cerr != nil {
		return nil, cerr
	}

	if book == nil {
		return nil, &core.CustomError{
			Code:    core.UNPROCESSABLE_ENTITY,
			Message: "book not found",
		}
	}

	cerr = uc.redisSetBook(ic, book)
	if cerr != nil {
		return nil, cerr
	}

	return book, cerr
}

func (uc *bookUsecaseImpl) redisSetBook(ic *core.InternalContext, book *entity.Book) *core.CustomError {
	duration := time.Hour * 24
	redisKeyBook := fmt.Sprintf("books::%s", book.Id)
	cerr := uc.cacheRepository.Set(ic, redisKeyBook, helper.DataToString(book), &duration)
	if cerr != nil {
		log.Error(ic.ToContext(), "failed uc.cacheRepository.Set(ic, redisKeyBook, helper.DataToString(book), &duration)", cerr.Error())
		return cerr
	}

	return cerr
}

func (uc *bookUsecaseImpl) redisClearBooks(ic *core.InternalContext) *core.CustomError {
	// get total books
	redisKeyBooksTotal := "books::total"
	getTotal, cerr := uc.cacheRepository.Get(ic, redisKeyBooksTotal)
	if cerr != nil {
		return cerr
	}
	total := helper.DataToInt(getTotal)

	page := int(math.Ceil(float64(total) / float64(10)))
	for i := 1; i <= page; i++ {
		cerr = uc.cacheRepository.Delete(ic, fmt.Sprintf("books::page::%d", i))
		if cerr != nil {
			return cerr
		}
	}

	return uc.cacheRepository.Delete(ic, redisKeyBooksTotal)
}

package books

import (
	"dot/core/v1/port/book"
	"github.com/labstack/echo/v4"
	"net/http"

	"dot/core"
	"dot/interface/v1/general/common"
	"dot/pkg/util"
)

var log = util.NewLogger()

type Handler struct {
	bookUsecase book.BookUsecase
}

func New(service book.BookUsecase) *Handler {
	return &Handler{
		bookUsecase: service,
	}
}

// Create 	godoc
// @Summary		Create book.
// @Description	Create .
// @Tags		Books
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"Bearer token"
// @Param  		Body 			body		Book 	true	"body"
// @Failure		400
// @Failure		401
// @Failure		422
// @Failure		500
// @Router	/api/v1/books [post]
func (h *Handler) Create(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(Book)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(http.StatusBadRequest, common.NewResponse(res, nil))
	}

	result, cerr := h.bookUsecase.BookCreate(ic, request.Entity())
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	res := common.ResultMap[core.CREATED]
	return c.JSON(
		res.StatusCode,
		common.NewResponse(res, result),
	)
}

// List 	godoc
// @Summary		List books.
// @Description	List books.
// @Tags		Books
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string		true	"Bearer token"
// @Param		page			query		string		false 	"1"
// @Param		limit			query		string		false 	"10"
// @Failure		400
// @Failure		401
// @Failure		422
// @Failure		500
// @Router	/api/v1/books [get]
func (h *Handler) List(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(List)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(http.StatusBadRequest, common.NewResponse(res, nil))
	}

	if request.Page < 1 || request.Limit < 1 {
		request.Page = 1
		request.Limit = 10
	}

	result, total, cerr := h.bookUsecase.BookList(ic, map[string]any{
		"page":  request.Page,
		"limit": request.Limit,
	})
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	return c.JSON(
		http.StatusOK,
		common.NewListResponseWithMeta(common.ResultMap[core.OK], result, common.GetMeta(request.Page, request.Limit, total)),
	)
}

// Get 	godoc
// @Summary		Get book.
// @Description	Get book.
// @Tags		Books
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string		true	"Bearer token"
// @Param		id				path		string		false 	"65c1de91056ae9755c64ffba"
// @Failure		400
// @Failure		401
// @Failure		422
// @Failure		500
// @Router	/api/v1/books/{id} [get]
func (h *Handler) Get(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(Get)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(http.StatusBadRequest, common.NewResponse(res, nil))
	}

	result, cerr := h.bookUsecase.BookGet(ic, request.Id)
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	return c.JSON(
		http.StatusOK,
		common.NewResponse(common.ResultMap[core.OK], result),
	)
}

// Update 	godoc
// @Summary		Update book.
// @Description	Update .
// @Tags		Books
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"Bearer token"
// @Param		id				path		string	false 	"65c1de91056ae9755c64ffba"
// @Param  		Body 			body		Book 	true	"body"
// @Failure		400
// @Failure		401
// @Failure		422
// @Failure		500
// @Router	/api/v1/books/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(Book)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(http.StatusBadRequest, common.NewResponse(res, nil))
	}

	cerr := h.bookUsecase.BookUpdate(ic, request.Entity())
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	return c.JSON(
		http.StatusOK,
		common.NewResponse(common.ResultMap[core.OK], nil),
	)
}

// UpdatePartial 	godoc
// @Summary		Update book.
// @Description	Update .
// @Tags		Books
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"Bearer token"
// @Param		id				path		string	false 	"65c1de91056ae9755c64ffba"
// @Param  		Body 			body		BookPartial 	true	"body"
// @Failure		400
// @Failure		401
// @Failure		422
// @Failure		500
// @Router	/api/v1/books/{id} [patch]
func (h *Handler) UpdatePartial(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(BookPartial)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(http.StatusBadRequest, common.NewResponse(res, nil))
	}

	cerr := h.bookUsecase.BookUpdate(ic, request.Entity())
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	return c.JSON(
		http.StatusOK,
		common.NewResponse(common.ResultMap[core.OK], nil),
	)
}

// Delete 	godoc
// @Summary		Delete book.
// @Description	Delete .
// @Tags		Books
// @Accept		json
// @Produce		json
// @Param		Authorization	header		string	true	"Bearer token"
// @Param		id				path		string	false 	"65c1de91056ae9755c64ffba"
// @Failure		400
// @Failure		401
// @Failure		422
// @Failure		500
// @Router	/api/v1/books/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(Delete)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(http.StatusBadRequest, common.NewResponse(res, nil))
	}

	cerr := h.bookUsecase.BookDelete(ic, request.Id)
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	return c.JSON(
		http.StatusOK,
		common.NewResponse(common.ResultMap[core.OK], nil),
	)
}

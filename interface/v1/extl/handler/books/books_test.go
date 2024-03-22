package books

import (
	"context"
	_deps "dot/app/v1/deps"
	"dot/core"
	"dot/interface/v1/extl/handler/auth"
	"dot/interface/v1/general/common"
	appMiddleware "dot/interface/v1/general/middleware"
	"dot/pkg/helper"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm/module/apmechov4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	e            *echo.Echo
	feker        = faker.New()
	userId       = ""
	id           = ""
	email        = "test@email.com"
	password     = "password"
	accessToken  = ""
	refreshToken = ""
	deps         = _deps.BuildDependency()
)

func TestMain(t *testing.M) {
	e = echo.New()
	e.HideBanner = false
	e.HidePort = false
	e.Validator = common.NewValidator()

	e.Use(apmechov4.Middleware())
	e.Use(appMiddleware.ServiceTrackerID)
	e.Use(appMiddleware.ServiceRequestTime)
	e.Use(echoMiddleware.RemoveTrailingSlash())
	e.Use(appMiddleware.Recover())
	e.Use(appMiddleware.CORS())
	e.Use(appMiddleware.InternalContext)

	t.Run()

	cerr := deps.GetRepositories().AccountRepository.DeleteAccount(core.NewInternalContext(uuid.NewString()), userId)
	println(cerr)
}

func TestHandler_Register(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().AuthUsecase

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(helper.DataToString(map[string]any{
		"email":    email,
		"password": password,
	})))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := auth.Handler{uc}
	auth := e.Group("/api/v1/auth")
	auth.POST("/register", h.Register)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.CREATED)
	assert.NotNil(t, res.Data)
	assert.NotEmpty(t, res.Data.(map[string]any)["id"])
	userId = res.Data.(map[string]any)["id"].(string)
}

func TestHandler_Login(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().AuthUsecase

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(helper.DataToString(map[string]any{
		"email":    email,
		"password": password,
	})))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	h := auth.Handler{uc}
	auth := e.Group("/api/v1/auth")
	auth.POST("/login", h.Login)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
	assert.NotNil(t, res.Data)
	assert.NotEmpty(t, res.Data.(map[string]any)["access_token"])
	assert.NotEmpty(t, res.Data.(map[string]any)["refresh_token"])
	accessToken = res.Data.(map[string]any)["access_token"].(string)
}

func TestHandler_Create(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().BookUsecase

	req := httptest.NewRequest(http.MethodPost, "/api/v1/books", strings.NewReader(helper.DataToString(map[string]any{
		"author": userId,
		"name":   "string",
	})))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	book := e.Group("/api/v1/books")
	book.POST("", h.Create)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.CREATED)
	assert.NotNil(t, res.Data)
	assert.NotEmpty(t, res.Data.(map[string]any)["id"])
	id = res.Data.(map[string]any)["id"].(string)
}

func TestHandler_Get(t *testing.T) {
	time.Sleep(5 * time.Second)
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().BookUsecase

	req := httptest.NewRequest(http.MethodGet, "/api/v1/books/"+id, nil)
	println(req.RequestURI)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	book := e.Group("/api/v1/books")
	book.GET("/:id", h.Get)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
	assert.NotNil(t, res.Data)
	assert.Equal(t, res.Data.(map[string]any)["id"], id)
}

func TestHandler_Update(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().BookUsecase

	req := httptest.NewRequest(http.MethodPut, "/api/v1/books/"+id, strings.NewReader(helper.DataToString(map[string]any{
		"author": userId,
		"name":   "update",
	})))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	book := e.Group("/api/v1/books")
	book.PUT("/:id", h.Update)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
}

func TestHandler_UpdatePartial(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().BookUsecase

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/books/"+id, strings.NewReader(helper.DataToString(map[string]any{
		"name": "string",
	})))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	book := e.Group("/api/v1/books")
	book.PATCH("/:id", h.UpdatePartial)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
}

func TestHandler_List(t *testing.T) {
	time.Sleep(2 * time.Second)
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().BookUsecase

	req := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	book := e.Group("/api/v1/books")
	book.GET("", h.List)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
	assert.NotNil(t, res.Data)
	assert.Equal(t, len(res.Data.([]any)), 1)
	assert.Equal(t, res.Data.([]any)[0].(map[string]any)["name"], "string")
}

func TestHandler_Delete(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().BookUsecase

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/books/"+id, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	book := e.Group("/api/v1/books")
	book.DELETE("/:id", h.Delete)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
}

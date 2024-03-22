package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm/module/apmechov4"

	"dot/config"
	"dot/core"
	"dot/interface/v1/general/common"
	"dot/pkg/helper"

	_deps "dot/app/v1/deps"
	appMiddleware "dot/interface/v1/general/middleware"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var (
	e            *echo.Echo
	feker        = faker.New()
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

	cerr := deps.GetRepositories().AccountRepository.DeleteAccount(core.NewInternalContext(uuid.NewString()), id)
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
	h := Handler{uc}
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
	id = res.Data.(map[string]any)["id"].(string)
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
	h := Handler{uc}
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
	refreshToken = res.Data.(map[string]any)["refresh_token"].(string)
}

func TestHandler_PasswordChange(t *testing.T) {
	time.Sleep(3 * time.Second)
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().AuthUsecase

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password/change", strings.NewReader(helper.DataToString(map[string]any{
		"old_password":     password,
		"password":         "123456",
		"password_confirm": "123456",
	})))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	auth := e.Group("/api/v1/auth")
	auth.Use(appMiddleware.JwtAccessTokenMiddleware(uc, config.GetConfig().JwtSecretKeys...), appMiddleware.JwtAccessTokenParseMiddleware)
	auth.POST("/password/change", h.ChangePassword)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
}

func TestHandler_Relogin(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().AuthUsecase

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/relogin", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+refreshToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	auth := e.Group("/api/v1/auth")
	auth.Use(appMiddleware.JwtRefreshTokenMiddleware(uc, config.GetConfig().JwtSecretKeys...))
	auth.POST("/relogin", h.Relogin)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
	assert.NotEmpty(t, res.Data.(map[string]any)["access_token"])
	assert.NotEmpty(t, res.Data.(map[string]any)["refresh_token"])
	accessToken = res.Data.(map[string]any)["access_token"].(string)
	refreshToken = res.Data.(map[string]any)["refresh_token"].(string)
}

func TestHandler_Validate(t *testing.T) {
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().AuthUsecase

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	auth := e.Group("/api/v1/auth")
	auth.Use(appMiddleware.JwtTokenMiddleware(uc, config.GetConfig().JwtSecretKeys...))
	auth.POST("/validate", h.Validate)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
}

func TestHandler_Logout(t *testing.T) {
	time.Sleep(2 * time.Second)
	defer e.Shutdown(context.Background())

	uc := deps.GetServices().AuthUsecase

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	h := Handler{uc}
	auth := e.Group("/api/v1/auth")
	auth.Use(appMiddleware.JwtAccessTokenMiddleware(uc, config.GetConfig().JwtSecretKeys...))
	auth.POST("/logout", h.Logout)
	e.ServeHTTP(rec, req)

	// Assertions
	res := new(common.Response)
	helper.StringToStruct(rec.Body.String(), res)
	println(rec.Body.String())
	assert.Equal(t, res.Result.Code, core.OK)
}

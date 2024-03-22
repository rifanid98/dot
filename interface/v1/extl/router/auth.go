package router

import (
	"dot/interface/v1/extl/handler/auth"
	"github.com/labstack/echo/v4"

	"dot/app/v1/deps"
	"dot/config"
	"dot/interface/v1/general/middleware"
)

func authRouter(
	e *echo.Group,
	deps deps.IDependency,
) {
	service := deps.GetServices()
	handler := auth.New(service.AuthUsecase)
	cfg := deps.GetBase().Cfg

	auth := e.Group("/auth")
	auth.Use(middleware.InternalContext)
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.POST("/relogin", handler.Relogin, middleware.JwtRefreshTokenMiddleware(service.AuthUsecase, cfg.JwtSecretKeys...))
	auth.POST("/logout", handler.Logout, middleware.JwtAccessTokenMiddleware(service.AuthUsecase, cfg.JwtSecretKeys...))
	auth.POST("/validate", handler.Validate, middleware.JwtTokenMiddleware(service.AuthUsecase, cfg.JwtSecretKeys...))
	auth.POST("/password/change", handler.ChangePassword, middleware.JwtAccessTokenMiddleware(service.AuthUsecase, config.GetConfig().JwtSecretKeys...), middleware.JwtAccessTokenParseMiddleware)
}

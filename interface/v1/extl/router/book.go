package router

import (
	"github.com/labstack/echo/v4"

	"dot/app/v1/deps"
	"dot/interface/v1/extl/handler/books"
	"dot/interface/v1/general/middleware"
)

func bookRouter(
	e *echo.Group,
	deps deps.IDependency,
) {
	service := deps.GetServices()
	handler := books.New(service)
	cfg := deps.GetBase().Cfg

	books := e.Group("/books")
	books.Use(middleware.InternalContext)
	books.Use(middleware.JwtAccessTokenMiddleware(service.AuthUsecase, cfg.JwtSecretKeys...))
	books.Use(middleware.JwtAccessTokenParseMiddleware)
	books.POST("", handler.Create)
	books.GET("", handler.List)
	books.GET("/:id", handler.Get)
	books.PUT("/:id", handler.Update)
	books.PATCH("/:id", handler.UpdatePartial)
	books.DELETE("/:id", handler.Delete)
}

package router

import (
	"github.com/labstack/echo/v4"

	"dot/app/v1/deps"

	external "dot/interface/v1/extl/router"
	internal "dot/interface/v1/intl/router"
)

// Register godoc
// @title Majoo Auth Management API
// @version 1.0
// @description This is a Majoo Auth Management API Documentation.
// @termsOfService http://www.majoo.id/terms/
// @contact.name PT Majoo Teknologi Indonesia
// @contact.url https://www.majoo.id/
// @contact.message halo@majoo.id
// @securityDefinitions.bearer  BearerAuth
func Register(e *echo.Echo, deps deps.IDependency) {
	api := e.Group("/api/v1")
	internal.Register(api, deps)
	external.Register(api, deps)
}

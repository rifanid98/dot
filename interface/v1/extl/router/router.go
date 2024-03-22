package router

import (
	"github.com/labstack/echo/v4"

	"dot/app/v1/deps"
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
func Register(e *echo.Group, deps deps.IDependency) {
	authRouter(e, deps)
	bookRouter(e, deps)
	healthRouter(e, deps)
}

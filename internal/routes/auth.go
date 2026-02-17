package routes

import (
	"github.com/gk-dev10/sheguard_backend/internal/controller"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Group) {
	e.POST("/login", controller.Login)
	e.POST("/signup",controller.Signup)
	e.POST("/logout",controller.Logout)
}

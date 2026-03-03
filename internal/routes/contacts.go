package routes

import (
	"github.com/gk-dev10/sheguard_backend/internal/controller"
	"github.com/gk-dev10/sheguard_backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

func ContactRoutes(e *echo.Group) {
	// All contact routes require authentication
	e.Use(middleware.AppwriteAuth)

	e.GET("", controller.GetContacts)
	e.POST("", controller.CreateContact)
	e.PUT("/:id", controller.UpdateContact)
	e.DELETE("/:id", controller.DeleteContact)
	e.PATCH("/:id/type", controller.ToggleContactType)
	e.PATCH("/:id/pin", controller.ToggleContactPin)
}

package routes

import (
	"github.com/gk-dev10/sheguard_backend/internal/controller"
	"github.com/gk-dev10/sheguard_backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group){
	e.GET("/me",controller.GetMe,middleware.SupabaseAuth)
}
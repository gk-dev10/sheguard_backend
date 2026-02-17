package controller

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/gk-dev10/sheguard_backend/internal/db"
)

func GetMe(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)

	var uid pgtype.UUID
	_ = uid.Scan(userIDStr)

	profile, err := db.Queries.GetProfile(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "profile not found",
		})
	}

	return c.JSON(http.StatusOK, profile)
}

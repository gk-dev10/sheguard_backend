package controller

import (
	"net/http"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/dto"
	"github.com/gk-dev10/sheguard_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
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

func UpdateMe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}
	var uid pgtype.UUID
	_ = uid.Scan(userID)
	params := sqlc.UpdateProfileParams{
		ID:              uid,
		Name:            req.Name,
		PhoneNumber:     req.PhoneNumber,
		ProfileImageUrl: req.ProfileImageURL,
		BloodGroup: req.BloodGroup,
		Allergies:       req.Allergies,
		Medications:     req.Medications,
	}

	profile, err := db.Queries.UpdateProfile(
		c.Request().Context(),
		params,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "update failed",
		})
	}
	return c.JSON(http.StatusOK, profile)
}

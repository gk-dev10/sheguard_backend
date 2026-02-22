package controller

import (
	"net/http"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/dto"
	"github.com/gk-dev10/sheguard_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func CreateContact(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	var uid pgtype.UUID
	_ = uid.Scan(userIDStr)

	var req dto.CreateContactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	params := sqlc.CreateContactParams{
		UserID:      uid,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		ImageUri:    req.ImageURI,
		Type:        req.Type,
		IsPinned:    req.IsPinned,
	}

	contact, err := db.Queries.CreateContact(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create contact"})
	}

	return c.JSON(http.StatusCreated, contact)
}

func GetContacts(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	var uid pgtype.UUID
	_ = uid.Scan(userIDStr)

	contacts, err := db.Queries.GetContacts(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch contacts"})
	}

	return c.JSON(http.StatusOK, contacts)
}

func UpdateContact(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	var uid pgtype.UUID
	_ = uid.Scan(userIDStr)

	contactIDStr := c.Param("id")
	var contactID pgtype.UUID
	_ = contactID.Scan(contactIDStr)

	var req dto.UpdateContactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	// For UpdateContact, we need to pass nullable parameters
	// sqlc generates interface{} for nullable params if not configured otherwise, or *string
	// SQL: COALESCE($3, name) -> implies we pass the value or null.
	// We matched the query params to the DTO fields which are pointers.

	params := sqlc.UpdateContactParams{
		ID:          contactID,
		UserID:      uid,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		ImageUri:    req.ImageURI,
		Type:        req.Type,
		IsPinned:    req.IsPinned,
	}

	contact, err := db.Queries.UpdateContact(c.Request().Context(), params)
	if err != nil {
		// Check if rows affected is 0? sqlc Update usually returns one.
		// If using Exec results, we check error.
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update contact"})
	}

	return c.JSON(http.StatusOK, contact)
}

func DeleteContact(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	var uid pgtype.UUID
	_ = uid.Scan(userIDStr)

	contactIDStr := c.Param("id")
	var contactID pgtype.UUID
	_ = contactID.Scan(contactIDStr)

	params := sqlc.DeleteContactParams{
		ID:     contactID,
		UserID: uid,
	}

	err := db.Queries.DeleteContact(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete contact"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "contact deleted"})
}

func ToggleContactType(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	var uid pgtype.UUID
	_ = uid.Scan(userIDStr)

	contactIDStr := c.Param("id")
	var contactID pgtype.UUID
	_ = contactID.Scan(contactIDStr)

	params := sqlc.ToggleContactTypeParams{
		ID:     contactID,
		UserID: uid,
	}

	contact, err := db.Queries.ToggleContactType(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to toggle type"})
	}

	return c.JSON(http.StatusOK, contact)
}

func ToggleContactPin(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	var uid pgtype.UUID
	_ = uid.Scan(userIDStr)

	contactIDStr := c.Param("id")
	var contactID pgtype.UUID
	_ = contactID.Scan(contactIDStr)

	params := sqlc.ToggleContactPinParams{
		ID:     contactID,
		UserID: uid,
	}

	contact, err := db.Queries.ToggleContactPin(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to toggle pin"})
	}

	return c.JSON(http.StatusOK, contact)
}

package controller

import (
	"net/http"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/dto"
	"github.com/labstack/echo/v4"

	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"
)

func CreateContact(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req dto.CreateContactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	data := map[string]interface{}{
		"user_id":      userID,
		"name":         req.Name,
		"phone_number": req.PhoneNumber,
	}
	if req.ImageURI != nil {
		data["image_uri"] = *req.ImageURI
	}
	if req.Type != nil {
		data["type"] = *req.Type
	}
	if req.IsPinned != nil {
		data["is_pinned"] = *req.IsPinned
	}

	doc, err := db.Databases.CreateDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		id.Unique(),
		data,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create contact"})
	}

	var result map[string]interface{}
	doc.Decode(&result)
	return c.JSON(http.StatusCreated, result)
}

func GetContacts(c echo.Context) error {
	userID := c.Get("user_id").(string)

	docs, err := db.Databases.ListDocuments(
		db.DatabaseID,
		db.ContactsCollectionID,
		db.Databases.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.OrderDesc("is_pinned"),
			query.OrderAsc("name"),
		}),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch contacts"})
	}

	var result map[string]interface{}
	docs.Decode(&result)
	return c.JSON(http.StatusOK, result)
}

func UpdateContact(c echo.Context) error {
	userID := c.Get("user_id").(string)
	contactID := c.Param("id")

	var req dto.UpdateContactRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	// Verify ownership first
	existing, err := db.Databases.GetDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "contact not found"})
	}

	var existingData map[string]interface{}
	existing.Decode(&existingData)
	if existingData["user_id"] != userID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "not your contact"})
	}

	// Build partial update
	data := make(map[string]interface{})
	if req.Name != nil {
		data["name"] = *req.Name
	}
	if req.PhoneNumber != nil {
		data["phone_number"] = *req.PhoneNumber
	}
	if req.ImageURI != nil {
		data["image_uri"] = *req.ImageURI
	}
	if req.Type != nil {
		data["type"] = *req.Type
	}
	if req.IsPinned != nil {
		data["is_pinned"] = *req.IsPinned
	}

	if len(data) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no fields to update"})
	}

	doc, err := db.Databases.UpdateDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
		db.Databases.WithUpdateDocumentData(data),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update contact"})
	}

	var result map[string]interface{}
	doc.Decode(&result)
	return c.JSON(http.StatusOK, result)
}

func DeleteContact(c echo.Context) error {
	userID := c.Get("user_id").(string)
	contactID := c.Param("id")

	// Verify ownership
	existing, err := db.Databases.GetDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "contact not found"})
	}

	var existingData map[string]interface{}
	existing.Decode(&existingData)
	if existingData["user_id"] != userID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "not your contact"})
	}

	_, err = db.Databases.DeleteDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete contact"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "contact deleted"})
}

func ToggleContactType(c echo.Context) error {
	userID := c.Get("user_id").(string)
	contactID := c.Param("id")

	existing, err := db.Databases.GetDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "contact not found"})
	}

	var existingData map[string]interface{}
	existing.Decode(&existingData)
	if existingData["user_id"] != userID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "not your contact"})
	}

	currentType, _ := existingData["type"].(string)
	newType := "Trusted"
	if currentType == "Trusted" {
		newType = "Casual"
	}

	doc, err := db.Databases.UpdateDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
		db.Databases.WithUpdateDocumentData(map[string]interface{}{
			"type": newType,
		}),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to toggle type"})
	}

	var result map[string]interface{}
	doc.Decode(&result)
	return c.JSON(http.StatusOK, result)
}

func ToggleContactPin(c echo.Context) error {
	userID := c.Get("user_id").(string)
	contactID := c.Param("id")

	existing, err := db.Databases.GetDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "contact not found"})
	}

	var existingData map[string]interface{}
	existing.Decode(&existingData)
	if existingData["user_id"] != userID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "not your contact"})
	}

	currentPinned, _ := existingData["is_pinned"].(bool)

	doc, err := db.Databases.UpdateDocument(
		db.DatabaseID,
		db.ContactsCollectionID,
		contactID,
		db.Databases.WithUpdateDocumentData(map[string]interface{}{
			"is_pinned": !currentPinned,
		}),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to toggle pin"})
	}

	var result map[string]interface{}
	doc.Decode(&result)
	return c.JSON(http.StatusOK, result)
}

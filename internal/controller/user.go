package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/dto"
	"github.com/labstack/echo/v4"
)

func GetMe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	doc, err := db.Databases.GetDocument(
		db.DatabaseID,
		db.ProfilesCollectionID,
		userID,
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
	}

	// Decode the document into a generic map for the response
	var profile map[string]interface{}
	if decodeErr := doc.Decode(&profile); decodeErr != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to decode profile"})
	}

	return c.JSON(http.StatusOK, profile)
}

func UpdateMe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	// Build update data — only include non-nil fields
	data := make(map[string]interface{})
	if req.Name != nil {
		data["name"] = *req.Name
	}
	if req.PhoneNumber != nil {
		data["phone_number"] = *req.PhoneNumber
	}
	if req.ProfileImageURL != nil {
		data["profile_image_url"] = *req.ProfileImageURL
	}
	if req.BloodGroup != nil {
		data["blood_group"] = *req.BloodGroup
	}
	if req.Allergies != nil {
		data["allergies"] = *req.Allergies
	}
	if req.Medications != nil {
		data["medications"] = *req.Medications
	}

	if len(data) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no fields to update"})
	}

	doc, err := db.Databases.UpdateDocument(
		db.DatabaseID,
		db.ProfilesCollectionID,
		userID,
		db.Databases.WithUpdateDocumentData(data),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	var result map[string]interface{}
	if decodeErr := doc.Decode(&result); decodeErr != nil {
		// Fallback: return raw fields
		return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
	}

	return c.JSON(http.StatusOK, result)
}

// decodeDocument decodes an Appwrite document into a generic map.
func decodeDocument(data []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

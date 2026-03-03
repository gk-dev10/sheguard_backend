package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/dto"
	"github.com/labstack/echo/v4"

	"github.com/appwrite/sdk-for-go/id"
)

func Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Proxy login to Appwrite REST API to create an email/password session.
	loginURL := fmt.Sprintf("%s/account/sessions/email", db.Client.Endpoint)
	body, _ := json.Marshal(map[string]string{
		"email":    req.Email,
		"password": req.Password,
	})

	httpReq, _ := http.NewRequest("POST", loginURL, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Appwrite-Project", db.Client.Headers["X-Appwrite-Project"])

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Update last login in profile
		if userID, ok := result["userId"].(string); ok {
			go func() {
				db.Databases.UpdateDocument(
					db.DatabaseID,
					db.ProfilesCollectionID,
					userID,
					db.Databases.WithUpdateDocumentData(map[string]interface{}{
						"last_login_at": time.Now().UTC().Format(time.RFC3339),
					}),
				)
			}()

			// Create a JWT for the user to use with our backend
			jwt, jwtErr := db.Users.CreateJWT(
				userID,
				db.Users.WithCreateJWTDuration(3600),
			)
			if jwtErr == nil {
				result["jwt"] = jwt.Jwt
			}
		}
	}

	return c.JSON(resp.StatusCode, result)
}

func Logout(c echo.Context) error {
	// With Appwrite, logout is handled client-side by deleting the session.
	return c.JSON(http.StatusOK, echo.Map{"message": "logged out"})
}

func Signup(c echo.Context) error {
	var req dto.SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Create user via server SDK (Users API)
	user, err := db.Users.Create(
		id.Unique(),
		db.Users.WithCreateEmail(req.Email),
		db.Users.WithCreatePassword(req.Password),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "signup failed: " + err.Error()})
	}

	// Create a profile document using the user's ID as the document ID
	_, profileErr := db.Databases.CreateDocument(
		db.DatabaseID,
		db.ProfilesCollectionID,
		user.Id,
		map[string]interface{}{},
	)
	if profileErr != nil {
		fmt.Printf("Warning: failed to create profile for user %s: %v\n", user.Id, profileErr)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"user_id": user.Id,
		"email":   user.Email,
		"message": "signup successful",
	})
}

func RefreshToken(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authentication required"})
	}

	jwt, err := db.Users.CreateJWT(
		userID.(string),
		db.Users.WithCreateJWTDuration(3600),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"jwt": jwt.Jwt,
	})
}

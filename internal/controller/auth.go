package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgtype"

)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password len:4 num"`
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}
	if !utils.IsValidPIN(req.Password) {
	return c.JSON(http.StatusBadRequest, echo.Map{
		"error": "PIN must be exactly 4 digits",
	})
}

	url := os.Getenv("SUPABASE_URL") + "/auth/v1/token?grant_type=password"
	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	var uid pgtype.UUID
	if userRaw, ok := result["user"]; ok {
		if userMap, ok := userRaw.(map[string]interface{}); ok {
			if id, ok := userMap["id"].(string); ok {
				uid.Scan(id)
				db.Queries.UpdateLastLogin(c.Request().Context(), uid)
			}
		}
	}

	return c.JSON(resp.StatusCode, result)
}

func Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "missing token",
		})
	}

	url := os.Getenv("SUPABASE_URL") + "/auth/v1/logout"

	httpReq, _ := http.NewRequest("POST", url, nil)
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()
	return c.JSON(http.StatusOK, echo.Map{
		"message": "logged out",
	})
}

func Signup(c echo.Context) error {
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}
		if !utils.IsValidPIN(req.Password) {
	return c.JSON(http.StatusBadRequest, echo.Map{
		"error": "PIN must be exactly 4 digits",
	})
}
	url := os.Getenv("SUPABASE_URL") + "/auth/v1/signup"
	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	var uid pgtype.UUID
	if userRaw, ok := result["user"]; ok {
		if userMap, ok := userRaw.(map[string]interface{}); ok {
			if id, ok := userMap["id"].(string); ok {
				uid.Scan(id)
				db.Queries.CreateProfile(c.Request().Context(), uid)
			}
		}
	}
	return c.JSON(resp.StatusCode, result)
}
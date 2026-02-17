package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	url := os.Getenv("SUPABASE_URL") + "/auth/v1/token?grant_type=password"

	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return c.JSON(resp.StatusCode, result)
}

func Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing token",
		})
	}

	url := os.Getenv("SUPABASE_URL") + "/auth/v1/logout"

	httpReq, _ := http.NewRequest("POST", url, nil)
	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	return c.NoContent(http.StatusNoContent)
}
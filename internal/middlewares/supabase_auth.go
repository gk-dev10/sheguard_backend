package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func SupabaseAuth(next echo.HandlerFunc) echo.HandlerFunc {
	jwksURL := os.Getenv("SUPABASE_URL") + "/auth/v1/keys"
	jwks, _ := keyfunc.Get(jwksURL, keyfunc.Options{})

	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "missing token",
			})
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
		}

		return next(c)
	}
}
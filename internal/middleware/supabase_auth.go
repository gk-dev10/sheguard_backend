package middleware

import (
	"net/http"
	"strings"

	"github.com/gk-dev10/sheguard_backend/internal/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func SupabaseAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, auth.JWKS.Keyfunc)
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["sub"].(string)

		c.Set("user_id", userID)
		// fmt.Printf("userID: %v",userID)
		return next(c)
	}
}
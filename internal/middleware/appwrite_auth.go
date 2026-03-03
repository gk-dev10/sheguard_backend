package middleware

import (
	"net/http"
	"strings"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/labstack/echo/v4"
)

// AppwriteAuth validates the Appwrite JWT by calling Appwrite's account.get()
// with the JWT set on a per-request client. If valid, extracts the user ID.
func AppwriteAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing token"})
		}

		jwt := strings.TrimPrefix(authHeader, "Bearer ")

		// Create a per-request client with the user's JWT
		userClient := appwrite.NewClient(
			appwrite.WithEndpoint(db.Client.Endpoint),
			appwrite.WithProject(db.Client.Headers["X-Appwrite-Project"]),
			appwrite.WithJWT(jwt),
		)

		accountSvc := appwrite.NewAccount(userClient)
		user, err := accountSvc.Get()
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		c.Set("user_id", user.Id)
		return next(c)
	}
}

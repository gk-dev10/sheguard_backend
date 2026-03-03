package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/routes"
	"github.com/gk-dev10/sheguard_backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Test(c echo.Context) error {
	fmt.Println("Server is running , healthy ! ")
	return c.String(http.StatusOK, "Hello whomever!, This is sheguard's backend, Back up!!")
}

func main() {
	e := echo.New()
	e.Validator = &utils.CustomValidator{
		Validator: validator.New(),
	}
	godotenv.Load()

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	//Test route
	e.GET("/", Test)

	api := e.Group("/api")

	routes.AuthRoutes(api.Group("/auth"))
	routes.UserRoutes(api.Group(""))
	routes.ContactRoutes(api.Group("/contacts"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err)
	}
}

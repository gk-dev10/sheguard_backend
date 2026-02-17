package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gk-dev10/sheguard_backend/internal/auth"
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
	ctx := context.Background()
	err := db.Init(ctx)
	if err!= nil{
		log.Fatal(err) //used log.fatal to end the program if db is not connected(no panic or no fmt.print)
	}
		if err := auth.InitJWKS(); err != nil {
		log.Fatal("JWKS init failed:", err)
	}

	//Test route
	e.GET("/", Test)

	api := e.Group("/api")

	routes.AuthRoutes(api.Group("/auth"))
	routes.UserRoutes(api.Group(""))

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}

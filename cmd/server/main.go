package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gk-dev10/sheguard_backend/internal/db"
	"github.com/gk-dev10/sheguard_backend/internal/routes"
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
	godotenv.Load()
	ctx := context.Background()
	err := db.Init(ctx)
	if err!= nil{
		log.Fatal(err) //to end the program if db is not connected
	}

	//Test route
	e.GET("/", Test)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	routes.AuthRoutes(v1.Group("/auth"))

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}

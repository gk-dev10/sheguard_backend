package main

import (
	"fmt"
	"net/http"

	"github.com/gk-dev10/sheguard_backend/internal/routes"
	"github.com/labstack/echo/v4"
)

func Test(c echo.Context) error {
	fmt.Println("Server is running , healthy ! ")
	return c.String(http.StatusOK, "Hello whomever!, This is sheguard's backend, Back up!!")
}

func main() {
	e := echo.New()

	//Test route
	e.GET("/", Test)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	routes.AuthRoutes(v1.Group("/auth"))

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}

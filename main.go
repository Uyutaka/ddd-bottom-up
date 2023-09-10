package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"uyutaka.com/ddd-bottom-up/model"
)

func main() {

	var userApplicationService model.UserApplicationService

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		fmt.Println(userApplicationService)
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

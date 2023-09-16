package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"uyutaka.com/ddd-bottom-up/model"
)

func main() {
	repo := model.NewSliceUserRepository("test")
	userService := model.NewUserService(&repo)
	userFactory := model.UserFactory{}
	userRepository := &repo
	userApplicationService := model.NewUserApplicationService(userService, &userFactory, userRepository)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		result, err := userApplicationService.GetAll()
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "error in GetAll()")
		}
		var output string
		if result != nil {

			for _, user := range result.Users {
				output += user.ToString() + "\n"
			}
		}

		return c.String(http.StatusOK, output)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

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
	userFactory := model.NewUserFactory(repo.Storage)
	userRepository := &repo
	userApplicationService := model.NewUserApplicationService(userService, &userFactory, userRepository)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		// curl localhost:1323
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

	e.GET("/:id", func(c echo.Context) error {
		// curl localhost:1323/1
		id := c.Param("id")
		command := model.UserGetCommand{UserId: id}

		result, err := userApplicationService.Get(command)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		response := model.NewUserResponseModel(result.User)
		if response == nil {
			return c.String(http.StatusOK, "error")
		}
		output := response.Id + " " + response.Name
		return c.String(http.StatusOK, output)
	})

	e.POST("/", func(c echo.Context) error {
		// curl -X POST --data-urlencode 'name=xxxx' localhost:1323
		command := model.UserRegisterCommand{Name: c.FormValue("name")}

		result, err := userApplicationService.Register(command)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.String(http.StatusOK, "userId: "+result.Id+" created!")

	})

	e.PUT("/:id", func(c echo.Context) error {
		// curl -X PUT --data-urlencode 'name=updated!' localhost:1323/1
		id := c.Param("id")
		command := model.UserUpdateCommand{Id: id, Name: c.FormValue("name")}
		err := userApplicationService.Update(command)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.String(http.StatusOK, "userId: "+id+" updated!")
	})

	e.DELETE("/:id", func(c echo.Context) error {
		// curl -X DELETE localhost:1323/1
		id := c.Param("id")
		command := model.UserDeleteCommand{Id: id}
		err := userApplicationService.Delete(command)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.String(http.StatusOK, "userId: "+id+" deleted!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

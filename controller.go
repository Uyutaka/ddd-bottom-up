package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"uyutaka.com/ddd-bottom-up/model"
)

func getUsers(c echo.Context) error {
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
}

func getUser(c echo.Context) error {
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
}

func createUser(c echo.Context) error {
	command := model.UserRegisterCommand{Name: c.FormValue("name")}

	result, err := userApplicationService.Register(command)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, "userId: "+result.Id+" created!")
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	command := model.UserUpdateCommand{Id: id, Name: c.FormValue("name")}
	err := userApplicationService.Update(command)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, "userId: "+id+" updated!")
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	command := model.UserDeleteCommand{Id: id}
	err := userApplicationService.Delete(command)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, "userId: "+id+" deleted!")
}

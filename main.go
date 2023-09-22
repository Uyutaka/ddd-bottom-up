package main

import (
	"github.com/labstack/echo"
	"uyutaka.com/ddd-bottom-up/application"
	"uyutaka.com/ddd-bottom-up/model"
)

var (
	userApplicationService application.UserApplicationService
)

func main() {
	repo := model.NewSliceUserRepository("test")
	userService := model.NewUserService(&repo)
	userFactory := model.NewUserFactory(repo.Storage)
	userRepository := &repo
	userApplicationService = application.NewUserApplicationService(userService, &userFactory, userRepository)

	e := echo.New()

	// curl localhost:1323
	e.GET("/", getUsers)

	// curl localhost:1323/1
	e.GET("/:id", getUser)

	// curl -X POST --data-urlencode 'name=xxxx' localhost:1323
	e.POST("/", createUser)

	// curl -X PUT --data-urlencode 'name=updated!' localhost:1323/1
	e.PUT("/:id", updateUser)

	// curl -X DELETE localhost:1323/1
	e.DELETE("/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

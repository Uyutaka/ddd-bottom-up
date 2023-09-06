package main

import (
	"fmt"

	"uyutaka.com/ddd-bottom-up/model"
)

func main() {
	var u2 model.User
	var c1 model.Circle
	userId, _ := model.NewUserId("1")
	userName, _ := model.NewUserName("test user name")

	u2, _ = model.NewUser(userId, userName)

	fmt.Println(u2)
	fmt.Println(c1)
}

package application

import "uyutaka.com/ddd-bottom-up/model"

type (
	UserRegisterCommand struct {
		Name string
	}

	UserRegisterResult struct {
		Id string
	}

	UserGetCommand struct {
		UserId string
	}

	UserGetResult struct {
		User model.User
	}

	UserGetAllResult struct {
		Users []model.User
	}

	UserUpdateCommand struct {
		Id   string
		Name string
	}

	UserDeleteCommand struct {
		Id string
	}
)

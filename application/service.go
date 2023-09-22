package application

import (
	"errors"

	"uyutaka.com/ddd-bottom-up/model"
)

type (
	UserApplicationService struct {
		UserService    model.UserService
		UserFactory    model.IUserFactory
		UserRepository model.IUserRepository
	}
)

func NewUserApplicationService(userService model.UserService, userFactory model.IUserFactory, userRepository model.IUserRepository) UserApplicationService {
	return UserApplicationService{UserService: userService, UserFactory: userFactory, UserRepository: userRepository}
}

func (uas *UserApplicationService) Get(command UserGetCommand) (*UserGetResult, error) {
	id, _ := model.NewUserId(command.UserId)
	user, _ := uas.UserRepository.FindById(&id)
	if user == nil {
		return nil, errors.New("user not found")
	}
	result := UserGetResult{User: *user}
	return &result, nil
}

func (uas *UserApplicationService) GetAll() (*UserGetAllResult, error) {
	users, errors := uas.UserRepository.FindAll()
	if errors != nil {
		return nil, errors
	}
	result := UserGetAllResult{Users: *users}
	return &result, nil
}

func (uas *UserApplicationService) Register(command UserRegisterCommand) (*UserRegisterResult, error) {

	// starts tx
	userName, _ := model.NewUserName(command.Name)
	user, _ := uas.UserFactory.Create(&userName)
	if uas.UserRepository.Exists(*user) {
		return nil, errors.New("user already exists")
	}

	uas.UserRepository.Save(*user)
	// ends tx

	return &UserRegisterResult{Id: user.Id.V}, nil
}

func (uas *UserApplicationService) Update(command UserUpdateCommand) error {
	// starts tx
	id, _ := model.NewUserId(command.Id)
	user, _ := uas.UserRepository.FindById(&id)
	if user == nil {
		return errors.New("user not found")
	}

	if len(command.Name) != 0 {
		name, _ := model.NewUserName(command.Name)
		user.ChangeName(&name)
	} else {
		return errors.New("name is empty")
	}
	uas.UserRepository.Save(*user)

	// ends tx
	return nil
}

func (uas *UserApplicationService) Delete(command UserDeleteCommand) error {
	// starts tx
	id, _ := model.NewUserId(command.Id)
	user, _ := uas.UserRepository.FindById(&id)
	if user == nil {
		return errors.New("user not found")
	}

	err := uas.UserRepository.Delete(*user)
	if err != nil {
		return errors.New("could not delete user")
	}
	// ends tx

	return nil
}

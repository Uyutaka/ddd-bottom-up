package model

import (
	"errors"
	"math/rand"
)

type (
	UserName struct {
		V string
	}

	UserId struct {
		V string
	}

	// Aggregate Root
	User struct {
		Id   UserId
		Name UserName
	}

	UserService struct {
		userRepository IUserRepository
	}

	UserApplicationService struct {
		userService    UserService
		userFactory    IUserFactory
		userRepository IUserRepository
	}

	UserRegisterCommand struct {
		name string
	}
	UserRegisterResult struct {
		id UserId
	}

	UserFactory struct{}

	UserGetCommand struct {
		userId string
	}

	UserGetResult struct {
		user User
	}

	UserGetAllResult struct {
		users []User
	}
	UserUpdateCommand struct {
		id   string
		name string
	}

	UserDeleteCommand struct {
		id string
	}

	IUserRepository interface {
		Save(user User) error
		Find(id *UserId) (*User, error)
		FindAll() (*[]User, error)
		Exists(user User) bool
		Delete(user User) error
	}

	IUserFactory interface {
		Create(name *UserName) (*User, error)
	}
)

func NewUserName(v string) (UserName, bool) {
	if len(v) == 0 {
		return UserName{}, false
	}
	if len(v) < 3 {
		return UserName{}, false
	}

	return UserName{V: v}, true
}

func NewUserId(v string) (UserId, bool) {
	if len(v) == 0 {
		return UserId{}, false
	}

	return UserId{V: v}, true
}

func NewUser(id UserId, name UserName) (User, bool) {
	if len(id.V) == 0 {
		return User{}, false
	}

	return User{Id: id}, true
}

func NewCircleId(v string) (CircleId, bool) {
	if len(v) == 0 {
		return CircleId{}, false
	}
	return CircleId{V: v}, true
}

func NewCircleName(v string) (CircleName, bool) {
	if len(v) == 0 {
		return CircleName{}, false
	}
	if len(v) < 3 {
		return CircleName{}, false
	}
	if len(v) > 20 {
		return CircleName{}, false
	}

	return CircleName{V: v}, true
}

func NewUserService(userRepository IUserRepository) UserService {
	return UserService{userRepository: userRepository}
}

func NewUserApplicationService(userService UserService, userFactory IUserFactory, userRepository IUserRepository) UserApplicationService {
	return UserApplicationService{userService: userService, userFactory: userFactory, userRepository: userRepository}
}

func (u *User) ChangeUserName(name string) bool {

	if len(name) == 0 {
		return false
	}
	if len(name) < 3 {
		return false
	}

	userName, _ := NewUserName(name)
	u.Name = userName
	return true
}

func (u *User) IsPremium() bool {
	return rand.Float32() < 0.5
}

func (uf *UserFactory) Create(id *UserId, name *UserName) (*User, error) {
	return &User{Name: *name}, nil
}

func (uas *UserApplicationService) Get(command UserGetCommand) (*UserGetResult, error) {
	id, _ := NewUserId(command.userId)
	user, _ := uas.userRepository.Find(&id)
	if user == nil {
		return nil, errors.New("user not found")
	}
	result := UserGetResult{user: *user}
	return &result, nil
}

func (uas *UserApplicationService) GetAll() (*UserGetAllResult, error) {
	users, errors := uas.userRepository.FindAll()
	if errors != nil {
		return nil, errors
	}
	result := UserGetAllResult{users: *users}
	return &result, nil
}

func (uas *UserApplicationService) Register(command UserRegisterCommand) (*UserRegisterResult, error) {

	// starts tx
	userName, _ := NewUserName(command.name)
	user, _ := uas.userFactory.Create(&userName)
	if uas.userRepository.Exists(*user) {
		return nil, errors.New("user already exists")
	}

	uas.userRepository.Save(*user)
	// ends tx

	return &UserRegisterResult{id: user.Id}, nil
}

func (uas *UserApplicationService) Update(command UserUpdateCommand) error {
	// starts tx
	id, _ := NewUserId(command.id)
	user, _ := uas.userRepository.Find(&id)
	if user == nil {
		return errors.New("user not found")
	}

	if len(command.name) != 0 {
		name, _ := NewUserName(command.name)
		user.ChangeUserName(name.V)
		if uas.userService.userRepository.Exists(*user) {
			return errors.New("user already exists")
		}
	}
	uas.userRepository.Save(*user)

	// ends tx
	return nil
}

func (uas *UserApplicationService) Delete(command UserDeleteCommand) error {
	// starts tx
	id, _ := NewUserId(command.id)
	user, _ := uas.userRepository.Find(&id)
	if user == nil {
		return errors.New("user not found")
	}

	err := uas.userRepository.Delete(*user)
	if err != nil {
		return errors.New("could not delete user")
	}
	// ends tx

	return nil
}

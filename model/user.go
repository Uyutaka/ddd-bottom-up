package model

import (
	"errors"
)

var (
	USER_TYPE_PREMIUM = UserType{V: "premium"}
	USER_TYPE_NORMAL  = UserType{V: "normal"}
)

type (
	UserName struct {
		V string
	}

	UserId struct {
		V string
	}

	UserType struct {
		V string
	}

	// Aggregate Root
	User struct {
		Id    UserId
		Name  UserName
		UType UserType
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
		FindById(id *UserId) (*User, error)
		FindByName(name *UserName) (*User, error)
		FindAll() (*[]User, error)
		Exists(user User) bool
		Delete(user User) error
	}

	IUserFactory interface {
		Create(name *UserName) (*User, error)
	}

	SliceUserRepository struct {
		connectionInfo string
		tmpUserStorage []User
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

func NewUserType(v string) (UserType, bool) {
	return UserType{V: v}, true
}

func NewUser(id UserId, name UserName, uType UserType) (User, bool) {
	if len(id.V) == 0 {
		return User{}, false
	}
	if len(name.V) == 0 {
		return User{}, false
	}

	return User{Id: id, Name: name, UType: uType}, true
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

func (u *User) ChangeName(name *UserName) bool {
	if name == nil {
		return false
	}
	u.Name = *name
	return true
}

func (u *User) Upgrade() {
	u.UType = USER_TYPE_PREMIUM
}

func (u *User) DownGrade() {
	u.UType = USER_TYPE_NORMAL
}

func (u *User) IsPremium() bool {
	return u.UType.V == "premium"
}

func (u *User) ToString() string {
	return u.Id.V + " " + u.Name.V + " " + u.UType.V
}

func (us *UserService) Exists(user *User) bool {
	duplicatedUser, _ := us.userRepository.FindByName(&user.Name)
	return duplicatedUser != nil
}

func (uf *UserFactory) Create(id *UserId, name *UserName) (*User, error) {
	return &User{Name: *name}, nil
}

func (uas *UserApplicationService) Get(command UserGetCommand) (*UserGetResult, error) {
	id, _ := NewUserId(command.userId)
	user, _ := uas.userRepository.FindById(&id)
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
	user, _ := uas.userRepository.FindById(&id)
	if user == nil {
		return errors.New("user not found")
	}

	if len(command.name) != 0 {
		name, _ := NewUserName(command.name)
		user.ChangeName(&name)
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
	user, _ := uas.userRepository.FindById(&id)
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

func NewSliceUserRepository(connectionInfo string) SliceUserRepository {
	return SliceUserRepository{connectionInfo: connectionInfo, tmpUserStorage: []User{}}
}

func (sur *SliceUserRepository) Save(user User) error {
	sur.tmpUserStorage = append(sur.tmpUserStorage, user)
	return nil
}

func (sur *SliceUserRepository) FindById(id *UserId) (*User, error) {
	for _, user := range sur.tmpUserStorage {
		if user.Id.V == id.V {
			return &user, nil
		}
	}
	return nil, nil
}

func (sur *SliceUserRepository) FindByName(name *UserName) (*User, error) {
	for _, user := range sur.tmpUserStorage {
		if user.Name.V == name.V {
			return &user, nil
		}
	}
	return nil, nil
}

func (sur *SliceUserRepository) FindAll() (*[]User, error) {
	return &sur.tmpUserStorage, nil
}

func (sur *SliceUserRepository) Exists(user User) bool {
	for _, u := range sur.tmpUserStorage {
		if u.Id.V == user.Id.V {
			return true
		}
	}
	return false
}

func (sur *SliceUserRepository) Delete(user User) error {
	for i, u := range sur.tmpUserStorage {
		if u.Id.V == user.Id.V {
			sur.tmpUserStorage = append(sur.tmpUserStorage[:i], sur.tmpUserStorage[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

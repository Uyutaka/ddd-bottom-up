package model

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

	UserResponseModel struct {
		Id   string
		Name string
	}

	UserPostRequestModel struct {
		Name string
	}

	UserService struct {
		userRepository IUserRepository
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

func NewUserResponseModel(user User) *UserResponseModel {
	return &UserResponseModel{Id: user.Id.V, Name: user.Name.V}
}

func NewUserService(userRepository IUserRepository) UserService {
	return UserService{userRepository: userRepository}
}

func (us *UserService) Exists(user *User) bool {
	duplicatedUser, _ := us.userRepository.FindByName(&user.Name)
	return duplicatedUser != nil
}

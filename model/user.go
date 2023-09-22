package model

import (
	"errors"
	"strconv"
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

	UserFactory struct {
		storage *TmpUserStorage
	}

	SliceUserRepository struct {
		connectionInfo string
		Storage        *TmpUserStorage
	}

	UserResponseModel struct {
		Id   string
		Name string
	}

	UserPostRequestModel struct {
		Name string
	}

	TmpUserStorage struct {
		data []User
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

func NewUserFactory(storage *TmpUserStorage) UserFactory {
	return UserFactory{storage: storage}
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

func (uf *UserFactory) Create(name *UserName) (*User, error) {
	id := uf.assignId()
	userId, _ := NewUserId(id)
	return &User{Name: *name, Id: userId, UType: USER_TYPE_NORMAL}, nil
}

func (uf *UserFactory) assignId() string {
	max := 0
	for _, user := range uf.storage.data {
		intId, err := strconv.Atoi(user.Id.V)
		if err != nil {
			break
		}
		if max < intId {
			max = intId
		}
	}
	return strconv.Itoa(max + 1)
}

func NewSliceUserRepository(connectionInfo string) SliceUserRepository {
	storage := TmpUserStorage{data: []User{
		User{Id: UserId{V: "1"}, Name: UserName{V: "user1"}, UType: USER_TYPE_NORMAL},
		User{Id: UserId{V: "2"}, Name: UserName{V: "user2"}, UType: USER_TYPE_PREMIUM},
	}}
	return SliceUserRepository{connectionInfo: connectionInfo, Storage: &storage}
}

func (sur *SliceUserRepository) Save(user User) error {
	if sur.Exists(user) {
		sur.Storage.Update(user)
	} else {
		sur.Storage.Insert(user)
	}
	return nil
}

func (sur *SliceUserRepository) FindById(id *UserId) (*User, error) {
	for _, user := range sur.Storage.data {
		if user.Id.V == id.V {
			return &user, nil
		}
	}
	return nil, nil
}

func (sur *SliceUserRepository) FindByName(name *UserName) (*User, error) {
	for _, user := range sur.Storage.data {
		if user.Name.V == name.V {
			return &user, nil
		}
	}
	return nil, nil
}

func (sur *SliceUserRepository) FindAll() (*[]User, error) {
	return &sur.Storage.data, nil
}

func (sur *SliceUserRepository) Exists(user User) bool {
	for _, u := range sur.Storage.data {
		if u.Id.V == user.Id.V {
			return true
		}
	}
	return false
}

func (sur *SliceUserRepository) Delete(user User) error {
	for i, u := range sur.Storage.data {
		if u.Id.V == user.Id.V {
			sur.Storage.data = append(sur.Storage.data[:i], sur.Storage.data[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func NewUserResponseModel(user User) *UserResponseModel {
	return &UserResponseModel{Id: user.Id.V, Name: user.Name.V}
}

func (tus *TmpUserStorage) Insert(user User) {
	tus.data = append(tus.data, user)
}

func (tus *TmpUserStorage) Update(user User) {
	for i, u := range tus.data {
		if u.Id.V == user.Id.V {
			tus.data[i] = user
			return
		}
	}
}

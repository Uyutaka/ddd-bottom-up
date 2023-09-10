package model

import (
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

	IUserRepository interface {
		Save(user User) error
		Find(id *UserId) (*User, error)
		Exists(user User) (bool, error)
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

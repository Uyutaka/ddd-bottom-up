package inMemoryInfrastructure

import (
	"errors"

	"uyutaka.com/ddd-bottom-up/model"
)

type (
	TmpUserStorage struct {
		data []model.User
	}
	SliceUserRepository struct {
		connectionInfo string
		Storage        *TmpUserStorage
	}
)

func NewSliceUserRepository(connectionInfo string) SliceUserRepository {
	storage := TmpUserStorage{data: []model.User{
		{Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL},
		{Id: model.UserId{V: "2"}, Name: model.UserName{V: "user2"}, UType: model.USER_TYPE_PREMIUM},
	}}
	return SliceUserRepository{connectionInfo: connectionInfo, Storage: &storage}
}

func (sur *SliceUserRepository) Save(user model.User) error {
	if sur.Exists(user) {
		sur.Storage.Update(user)
	} else {
		sur.Storage.Insert(user)
	}
	return nil
}

func (sur *SliceUserRepository) FindById(id *model.UserId) (*model.User, error) {
	for _, user := range sur.Storage.data {
		if user.Id.V == id.V {
			return &user, nil
		}
	}
	return nil, nil
}

func (sur *SliceUserRepository) FindByName(name *model.UserName) (*model.User, error) {
	for _, user := range sur.Storage.data {
		if user.Name.V == name.V {
			return &user, nil
		}
	}
	return nil, nil
}

func (sur *SliceUserRepository) FindAll() (*[]model.User, error) {
	return &sur.Storage.data, nil
}

func (sur *SliceUserRepository) Exists(user model.User) bool {
	for _, u := range sur.Storage.data {
		if u.Id.V == user.Id.V {
			return true
		}
	}
	return false
}

func (sur *SliceUserRepository) Delete(user model.User) error {
	for i, u := range sur.Storage.data {
		if u.Id.V == user.Id.V {
			sur.Storage.data = append(sur.Storage.data[:i], sur.Storage.data[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func (tus *TmpUserStorage) Insert(user model.User) {
	tus.data = append(tus.data, user)
}

func (tus *TmpUserStorage) Update(user model.User) {
	for i, u := range tus.data {
		if u.Id.V == user.Id.V {
			tus.data[i] = user
			return
		}
	}
}

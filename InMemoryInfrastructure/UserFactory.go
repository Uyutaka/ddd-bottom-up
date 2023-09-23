package inMemoryInfrastructure

import (
	"strconv"

	"uyutaka.com/ddd-bottom-up/model"
)

type (
	UserFactory struct {
		storage *TmpUserStorage
	}
)

func NewUserFactory(storage *TmpUserStorage) UserFactory {
	return UserFactory{storage: storage}
}

func (uf *UserFactory) Create(name *model.UserName) (*model.User, error) {
	id := uf.assignId()
	userId, _ := model.NewUserId(id)
	return &model.User{Name: *name, Id: userId, UType: model.USER_TYPE_NORMAL}, nil
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

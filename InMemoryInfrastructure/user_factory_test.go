package inMemoryInfrastructure

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"uyutaka.com/ddd-bottom-up/model"
)

func TestUserFactory_Create(t *testing.T) {
	type fields struct {
		storage *TmpUserStorage
	}
	type args struct {
		name *model.UserName
	}
	type wants struct {
		user  model.User
		error error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  wants
	}{
		{
			name: "normal",
			fields: fields{
				storage: &TmpUserStorage{data: []model.User{
					{Id: model.UserId{V: "1"}, Name: model.UserName{V: "test_user1"}, UType: model.USER_TYPE_NORMAL},
					{Id: model.UserId{V: "2"}, Name: model.UserName{V: "test_user2"}, UType: model.USER_TYPE_NORMAL},
				}},
			},
			args: args{name: &model.UserName{V: "test_user3"}},
			wants: wants{
				user:  model.User{Id: model.UserId{V: "3"}, Name: model.UserName{V: "test_user3"}, UType: model.USER_TYPE_NORMAL},
				error: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf := &UserFactory{
				storage: tt.fields.storage,
			}
			user, err := uf.Create(tt.args.name)
			assert.Equal(t, tt.wants.error, err,
				fmt.Sprintf("UserFactory.Create() error = %v, wantErr %v", err, tt.wants.error))

			assert.Equal(t, true, reflect.DeepEqual(user, &tt.wants.user),
				fmt.Sprintf("UserFactory.Create() = %v, want %v", user, tt.wants.user))
		})
	}
}

func TestUserFactory_assignId(t *testing.T) {
	type fields struct {
		storage *TmpUserStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				storage: &TmpUserStorage{data: []model.User{
					{Id: model.UserId{V: "1"}, Name: model.UserName{V: "test_user1"}, UType: model.USER_TYPE_NORMAL},
					{Id: model.UserId{V: "2"}, Name: model.UserName{V: "test_user2"}, UType: model.USER_TYPE_NORMAL},
				}},
			},
			want: "3",
		},
		{
			name: "empty data",
			fields: fields{
				storage: &TmpUserStorage{data: []model.User{}},
			},
			want: "1",
		},
		{
			name: "not continuously incrementing data",
			fields: fields{
				storage: &TmpUserStorage{data: []model.User{
					{Id: model.UserId{V: "1"}, Name: model.UserName{V: "test_user1"}, UType: model.USER_TYPE_NORMAL},
					{Id: model.UserId{V: "4"}, Name: model.UserName{V: "test_user2"}, UType: model.USER_TYPE_NORMAL},
				}},
			},
			want: "5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf := &UserFactory{
				storage: tt.fields.storage,
			}
			if got := uf.assignId(); got != tt.want {
				t.Errorf("UserFactory.assignId() = %v, want %v", got, tt.want)
			}
		})
	}
}

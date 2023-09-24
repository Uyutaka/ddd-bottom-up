package inMemoryInfrastructure

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"uyutaka.com/ddd-bottom-up/model"
)

func TestSliceUserRepository_Exists(t *testing.T) {
	type fields struct {
		Storage *TmpUserStorage
	}
	type args struct {
		user model.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "normal - exists",
			fields: fields{
				Storage: &TmpUserStorage{data: []model.User{
					{Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL},
				}},
			},
			args: args{
				user: model.User{
					Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL,
				},
			},
			want: true,
		},
		{
			name: "normal - not exists",
			fields: fields{
				Storage: &TmpUserStorage{data: []model.User{}},
			},
			args: args{
				user: model.User{
					Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL,
				},
			},
			want: false,
		},
		{
			name: "normal - not exists",
			fields: fields{
				Storage: &TmpUserStorage{data: []model.User{
					{Id: model.UserId{V: "2"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL},
				}},
			},
			args: args{
				user: model.User{
					Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sur := &SliceUserRepository{
				Storage: tt.fields.Storage,
			}
			got := sur.Exists(tt.args.user)
			assert.Equal(t, tt.want, got,
				fmt.Sprintf("SliceUserRepository.Exists() = %v, want %v", got, tt.want))
		})
	}
}

func TestTmpUserStorage_Update(t *testing.T) {
	type fields struct {
		data []model.User
	}
	type args struct {
		user model.User
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		expectedUpdatedFields []model.User
	}{
		{
			name: "successfully updated",
			fields: fields{
				data: []model.User{
					{Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL},
				},
			},
			args: args{
				user: model.User{Id: model.UserId{V: "1"}, Name: model.UserName{V: "updated_user1"}, UType: model.USER_TYPE_PREMIUM},
			},
			expectedUpdatedFields: []model.User{
				{Id: model.UserId{V: "1"}, Name: model.UserName{V: "updated_user1"}, UType: model.USER_TYPE_PREMIUM},
			},
		},
		{
			name: "failed to find target user",
			fields: fields{
				data: []model.User{
					{Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL},
				},
			},
			args: args{
				user: model.User{Id: model.UserId{V: "3"}, Name: model.UserName{V: "updated_user3"}, UType: model.USER_TYPE_PREMIUM},
			},
			expectedUpdatedFields: []model.User{
				{Id: model.UserId{V: "1"}, Name: model.UserName{V: "user1"}, UType: model.USER_TYPE_NORMAL},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tus := &TmpUserStorage{
				data: tt.fields.data,
			}
			tus.Update(tt.args.user)
			assert.Equal(t, true, reflect.DeepEqual(tus.data, tt.expectedUpdatedFields),
				fmt.Sprintf("tus.data = %v, expectedUpdatedFields = %v", tus.data, tt.expectedUpdatedFields))
		})
	}
}

func TestTmpUserStorage_Insert(t *testing.T) {
	type fields struct {
		data []model.User
	}
	type args struct {
		user model.User
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		expectedUpdatedFields []model.User
	}{
		{
			name: "normal",
			fields: fields{data: []model.User{
				{Id: model.UserId{V: "1"}, Name: model.UserName{V: "test_user1"}, UType: model.USER_TYPE_NORMAL},
			}},
			args: args{
				user: model.User{Id: model.UserId{V: "2"}, Name: model.UserName{V: "test_user2"}, UType: model.USER_TYPE_PREMIUM},
			},
			expectedUpdatedFields: []model.User{
				{Id: model.UserId{V: "1"}, Name: model.UserName{V: "test_user1"}, UType: model.USER_TYPE_NORMAL},
				{Id: model.UserId{V: "2"}, Name: model.UserName{V: "test_user2"}, UType: model.USER_TYPE_PREMIUM},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tus := &TmpUserStorage{
				data: tt.fields.data,
			}
			tus.Insert(tt.args.user)
			assert.Equal(t, true, reflect.DeepEqual(tus.data, tt.expectedUpdatedFields),
				fmt.Sprintf("tus.data = %v, expectedUpdatedFields = %v", tus.data, tt.expectedUpdatedFields))

		})
	}
}

// TODO FindById()

package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserName(t *testing.T) {
	type args struct {
		v string
	}
	type expected struct {
		userName UserName
		ok       bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "normal",
			args: args{v: "username"},
			expected: expected{
				userName: UserName{V: "username"},
				ok:       true,
			},
		},
		{
			name: "empty",
			args: args{v: ""},
			expected: expected{
				userName: UserName{V: ""},
				ok:       false,
			},
		},
		{
			name: "3 string",
			args: args{v: "aiu"},
			expected: expected{
				userName: UserName{V: "aiu"},
				ok:       true,
			},
		},
		{
			name: "2 string",
			args: args{v: "ai"},
			expected: expected{
				userName: UserName{V: ""},
				ok:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := NewUserName(tt.args.v)
			assert.Equal(t, true, reflect.DeepEqual(got, tt.expected.userName),
				fmt.Sprintf("NewUserName() got = %v, want %v", got, tt.expected.userName))

			assert.Equal(t, tt.expected.ok, ok,
				fmt.Sprintf("NewUserName() got1 = %v, want %v", ok, tt.expected.ok))
		})
	}
}

func TestNewUserId(t *testing.T) {
	type args struct {
		v string
	}
	type wants struct {
		userId UserId
		ok     bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "normal",
			args:  args{v: "1"},
			wants: wants{userId: UserId{V: "1"}, ok: true},
		},
		{
			name:  "empty",
			args:  args{v: ""},
			wants: wants{userId: UserId{V: ""}, ok: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId, ok := NewUserId(tt.args.v)
			assert.Equal(t, true, reflect.DeepEqual(userId, tt.wants.userId),
				fmt.Sprintf("NewUserId() got = %v, want %v", userId, tt.wants.userId))

			assert.Equal(t, tt.wants.ok, ok,
				fmt.Sprintf("NewUserId() got1 = %v, want %v", ok, tt.wants.ok))
		})
	}
}

func TestNewUserType(t *testing.T) {
	type args struct {
		v string
	}
	type wants struct {
		userType UserType
		ok       bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name:  "normal",
			args:  args{v: "test_type"},
			wants: wants{userType: UserType{V: "test_type"}, ok: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userType, ok := NewUserType(tt.args.v)
			assert.Equal(t, true, reflect.DeepEqual(userType, tt.wants.userType),
				fmt.Sprintf("NewUserType() got = %v, want %v", userType, tt.wants.userType))

			assert.Equal(t, tt.wants.ok, ok,
				fmt.Sprintf("NewUserType() got1 = %v, want %v", ok, tt.wants.ok))
		})
	}
}

func TestNewUser(t *testing.T) {
	type args struct {
		id    UserId
		name  UserName
		uType UserType
	}
	type wants struct {
		user User
		ok   bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "normal",
			args: args{
				id:    UserId{V: "1"},
				name:  UserName{V: "test_name"},
				uType: UserType{V: "test_type"}},
			wants: wants{
				user: User{Id: UserId{V: "1"},
					Name:  UserName{V: "test_name"},
					UType: UserType{V: "test_type"}},
				ok: true},
		},
		{
			name: "empty UserId",
			args: args{
				id:    UserId{V: ""},
				name:  UserName{V: "test_name"},
				uType: UserType{V: "test_type"}},
			wants: wants{
				user: User{Id: UserId{V: ""},
					Name:  UserName{V: ""},
					UType: UserType{V: ""}},
				ok: false},
		},
		{
			name: "empty UserName",
			args: args{
				id:    UserId{V: "1"},
				name:  UserName{V: ""},
				uType: UserType{V: "test_type"}},
			wants: wants{
				user: User{Id: UserId{V: ""},
					Name:  UserName{V: ""},
					UType: UserType{V: ""}},
				ok: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, ok := NewUser(tt.args.id, tt.args.name, tt.args.uType)
			assert.Equal(t, true, reflect.DeepEqual(user, tt.wants.user),
				fmt.Sprintf("NewUser() got = %v, want %v", user, tt.wants.user))

			assert.Equal(t, tt.wants.ok, ok,
				fmt.Sprintf("NewUser() got1 = %v, want %v", ok, tt.wants.ok))
		})
	}
}

func TestUser_ChangeName(t *testing.T) {
	type fields struct {
		Id    UserId
		Name  UserName
		UType UserType
	}
	type args struct {
		name *UserName
	}
	type wants struct {
		ok          bool
		updatedName UserName
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
				Id:    UserId{"1"},
				Name:  UserName{"test_name"},
				UType: UserType{"test_type"},
			},
			args: args{
				name: &UserName{"updated_name"},
			},
			wants: wants{
				ok:          true,
				updatedName: UserName{"updated_name"},
			},
		},
		{
			name: "name is nil",
			fields: fields{
				Id:    UserId{"1"},
				Name:  UserName{"test_name"},
				UType: UserType{"test_type"},
			},
			args: args{
				name: nil,
			},
			wants: wants{
				ok:          false,
				updatedName: UserName{"test_name"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Id:    tt.fields.Id,
				Name:  tt.fields.Name,
				UType: tt.fields.UType,
			}
			ok := u.ChangeName(tt.args.name)

			assert.Equal(t, tt.wants.ok, ok,
				fmt.Sprintf("User.ChangeName() = %v, want %v", ok, tt.wants.ok))

			assert.Equal(t, true, reflect.DeepEqual(u.Name, tt.wants.updatedName),
				fmt.Sprintf("NewUser() got = %v, want %v", u.Name, tt.wants.updatedName))
		})
	}
}

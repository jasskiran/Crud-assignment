package repository

import (
	"testing"

	"Assignment/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_userRepository_Add(t *testing.T) {
	mock, db := NewMock()
	type args struct {
		uow *UnitOfWork
		out *models.User
	}
	defer mock.Close()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"createUserWithNameAndPassword",
			args{
				uow: &UnitOfWork{Db: mock},
				out: &models.User{
					Id:          1,
					Name:        "user1",
					Password:    "password",
					GitUsername: nil,
				},
			},
			false,
		},
		{
			"createUserWithNameAndPasswordAndGitUserName",
			args{
				uow: &UnitOfWork{Db: mock},
				out: &models.User{
					Id:          2,
					Name:        "user2",
					Password:    "password",
					GitUsername: "user2",
				},
			},
			false,
		},
		{
			"createUserWithNameAndPasswordAndGitUserNameWithSameUserName",
			args{
				uow: &UnitOfWork{Db: mock},
				out: &models.User{
					Id:          1,
					Name:        "user1",
					Password:    "password",
					GitUsername: "user1",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

			u := userRepository{}
			query := "INSERT INTO user (id, name, password, git_username) VALUES (?, ?, ?, ?)"
			prep := db.ExpectPrepare(query)
			prep.ExpectExec().WithArgs(tt.args.out.Id, tt.args.out.Name, tt.args.out.Password, tt.args.out.GitUsername).WillReturnResult(sqlmock.NewResult(0, 1))
			err := u.Add(tt.args.uow, tt.args.out)
			assert.NoError(t, err)
		})
	}
}

func Test_userRepository_Update(t *testing.T) {
	mock, db := NewMock()
	type args struct {
		uow  *UnitOfWork
		user *models.User
		Id   int
	}
	defer mock.Close()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"UpdateName",
			args{
				uow:  &UnitOfWork{Db: mock},
				user: &models.User{
					Name:  "user1",
				},
				Id:   1,
			},
			false,
		},
		{
			"UpdateNameAndPassword",
			args{
				uow:  &UnitOfWork{Db: mock},
				user: &models.User{
					Name:  "user1",
				},
				Id:   1,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := userRepository{}
			query := "UPDATE user SET (name) VALUES (?) where id = ?"
			prep := db.ExpectPrepare(query)
			prep.ExpectExec().WithArgs(tt.args.user.Name, tt.args.Id).WillReturnResult(sqlmock.NewResult(0, 1))
			err := u.Add(tt.args.uow, tt.args.user)
			assert.NoError(t, err)
		})
	}

}

func Test_userRepository_GetLoggedInUser(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()
	type args struct {
		uow    *UnitOfWork
		userId int
	}
	tests := []struct {
		name     string
		args     args
		wantUser *models.User
		wantErr  bool
	}{
		{
			"GetUserWIthLoggedInDetails",
			args{
				uow:  &UnitOfWork{Db: mock},
				userId: 1,
			},
			&models.User{
				Id:          1,
				Name:        "user1",
				Password:    "password",
				GitUsername: "user1",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := userRepository{}
			query := "SELECT (id, name, password, git_username) from user where id = ?"
			rows := db.NewRows([]string{"id", "name", "password", "git_username"}).
				AddRow(tt.wantUser.Id, tt.wantUser.Name, tt.wantUser.Password, tt.wantUser.GitUsername)
			db.ExpectQuery(query).WithArgs(tt.args.userId).WillReturnRows(rows)
			user, err := u.GetLoggedInUser(tt.args.uow, tt.args.userId)
			assert.NotNil(t, user)
			assert.NoError(t, err)
		})
		}
}

func Test_userRepository_Login(t *testing.T) {
	mock, db := NewMock()

	type args struct {
		uow  *UnitOfWork
		name string
	}

	defer mock.Close()

	tests := []struct {
		name     string
		args     args
		wantUser *models.User
		wantErr  bool
	}{
		{
			"LoginUser",
			args{
				uow:  &UnitOfWork{Db: mock},
				name: "user1",
			},
			&models.User{
				Id:          1,
				Name:        "user1",
				Password:    "password",
				GitUsername: "user1",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := userRepository{}
			query := "SELECT (id, name, password, git_username) from user where id = ?"
			rows := db.NewRows([]string{"id", "name", "password", "git_username"}).
				AddRow(tt.wantUser.Id, tt.wantUser.Name, tt.wantUser.Password, tt.wantUser.GitUsername)

			db.ExpectQuery(query).WithArgs(tt.args.name).WillReturnRows(rows)

			user, err := u.Login(tt.args.uow, tt.args.name)
			assert.NotNil(t, user)
			assert.NoError(t, err)
		})
	}
}

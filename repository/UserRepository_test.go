package repository

import (
	"fmt"
	"testing"

	"Assignment/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_userRepository_Add(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()

	type args struct {
		uow *UnitOfWork
		out *models.User
	}

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
					GitUsername: "nil",
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
			query := "INSERT INTO user (.+) VALUES (.+)"
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
				uow: &UnitOfWork{Db: mock},
				user: &models.User{
					Name: "user",
				},
				Id: 1,
			},
			false,
		},
		{
			"UpdateNameAndPassword",
			args{
				uow: &UnitOfWork{Db: mock},
				user: &models.User{
					Name: "user1",
				},
				Id: 1,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := userRepository{}
			query := fmt.Sprintf(`UPDATE 
										user 
									SET 
										name = \? 
									WHERE 
										id = \?`,
			)
			prep := db.ExpectPrepare(query)
			prep.ExpectExec().WithArgs(tt.args.user.Name, tt.args.Id).WillReturnResult(sqlmock.NewResult(0, 1))
			err := u.Update(tt.args.uow, tt.args.user, tt.args.Id)
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
				uow:    &UnitOfWork{Db: mock},
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
			query := "SELECT id, name, password, git_username FROM user WHERE id = ?"
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
	defer mock.Close()

	type args struct {
		uow  *UnitOfWork
		name string
	}

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
			query := "SELECT id, name, password, git_username from user where name = ?"
			rows := db.NewRows([]string{"id", "name", "password", "git_username"}).
				AddRow(tt.wantUser.Id, tt.wantUser.Name, tt.wantUser.Password, tt.wantUser.GitUsername)

			db.ExpectQuery(query).WithArgs(tt.args.name).WillReturnRows(rows)

			user, err := u.Login(tt.args.uow, tt.args.name)
			assert.NotNil(t, user)
			assert.NoError(t, err)
		})
	}
}

func Test_userRepository_DeleteToken(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()

	type args struct {
		uow    *UnitOfWork
		userId int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"SetTokenToFalse",
			args{
				uow:    &UnitOfWork{Db: mock},
				userId: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := userRepository{}
			query := `Update authentication
							set
								active = \?
							where
								user_id = \?`
			prep := db.ExpectPrepare(query)
			prep.ExpectExec().WithArgs(false, tt.args.userId).WillReturnResult(sqlmock.NewResult(0, 1))
			err := u.DeleteToken(tt.args.uow, tt.args.userId)
			assert.NoError(t, err)
		})
	}
}

func Test_userRepository_GetToken(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()

	type args struct {
		uow    *UnitOfWork
		userId int
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Auth
		wantErr bool
	}{
		{
			"GetToken",
			args{
				uow:    &UnitOfWork{Db: mock},
				userId: 1,
			},
			&models.Auth{
				Id:     1,
				UserId: 1,
				Token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk4MjY0NDQsInVzZXJfaWQiOjd9.1zCj3ArZ_-159I6WLis4XCi5sC6qAO9NMymJLQTDKwE",
				Active: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := userRepository{}
			query := "SELECT id, user_id, token, active from authentication where active = 1 and user_id = ?"
			rows := db.NewRows([]string{"id", "user_id", "token", "active"}).
				AddRow(tt.want.Id, tt.want.UserId, tt.want.Token, tt.want.Active)
			db.ExpectQuery(query).WithArgs(tt.args.userId).WillReturnRows(rows)
			token, err := u.GetToken(tt.args.uow, tt.args.userId)
			assert.NotNil(t, token)
			assert.NoError(t, err)
		})
	}
}

func Test_userRepository_AddToken(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()

	type args struct {
		uow *UnitOfWork
		out *models.Auth
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{
				uow: &UnitOfWork{Db: mock},
				out: &models.Auth{
					Id:     1,
					UserId: 1,
					Token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk4MjY0NDQsInVzZXJfaWQiOjd9.1zCj3ArZ_-159I6WLis4XCi5sC6qAO9NMymJLQTDKwE",
					Active: true,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := userRepository{}
			query := "INSERT INTO authentication (.+) VALUES (.+)"
			prep := db.ExpectPrepare(query)
			prep.ExpectExec().WithArgs(tt.args.out.Id, tt.args.out.UserId, tt.args.out.Token, tt.args.out.Active).WillReturnResult(sqlmock.NewResult(0, 1))
			err := u.AddToken(tt.args.uow, tt.args.out)
			assert.NoError(t, err)
		})
	}
}

package repository

import (
	"testing"

	"Assignment/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_gitUserRepository_CreateGitUser(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()
	type args struct {
		uow *UnitOfWork
		out *models.Github
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
				out: &models.Github{
					Id:       1,
					UserName: "user1",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := gitUserRepository{}
			query := "INSERT INTO github (.+) VALUES (.+)"
			prep := db.ExpectPrepare(query)
			prep.ExpectExec().WithArgs(tt.args.out.Id, tt.args.out.UserName).WillReturnResult(sqlmock.NewResult(0, 1))
			err := u.CreateGitUser(tt.args.uow, tt.args.out)
			assert.NoError(t, err)
		})
	}
}

func Test_gitUserRepository_GetGitUser(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()

	type args struct {
		uow    *UnitOfWork
		userId int
	}
	tests := []struct {
		name        string
		args        args
		wantGitUser *models.Github
		wantErr     bool
	}{
		{
			"",
			args{
				uow:    &UnitOfWork{Db: mock},
				userId: 1,
			},
			&models.Github{
				Id:       1,
				UserName: "user1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := gitUserRepository{}

			query := "SELECT github.id, github.username FROM github LEFT OUTER JOIN user ON github.username = user.git_username WHERE user.id = ?"

			rows := db.NewRows([]string{"id", "username"}).
				AddRow(tt.wantGitUser.Id, tt.wantGitUser.UserName)
			db.ExpectQuery(query).WithArgs(tt.args.userId).WillReturnRows(rows)
			user, err := u.GetGitUser(tt.args.uow, tt.args.userId)
			assert.NotNil(t, user)
			assert.NoError(t, err)
		})
	}
}

/*
SELECT github.id, github.username FROM github LEFT OUTER JOIN user ON github.username = user.git_username WHERE user.id = ?
*/

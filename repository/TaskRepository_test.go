package repository

import (
	"testing"

	"Assignment/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_taskRepository_Create(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()
	type args struct {
		uow    *UnitOfWork
		userId int
		out    *models.Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{
				uow:    &UnitOfWork{Db: mock},
				userId: 1,
				out: &models.Task{
					Id:          1,
					Name:        "task1",
					Description: "description",
					StartDate:   "",
					EndDate:     "",
					ZoomLink:    "zoom1",
					MeetLink:    "nil",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			task := taskRepository{}
			query := "INSERT INTO task (.+) VALUES (.+)"
			prep := db.ExpectPrepare(query)
			prep.ExpectExec().WithArgs(tt.args.out.Id, tt.args.userId, tt.args.out.Name, tt.args.out.Description, tt.args.out.StartDate, tt.args.out.EndDate, tt.args.out.ZoomLink, tt.args.out.MeetLink).
				WillReturnResult(sqlmock.NewResult(0, 1))
			err := task.Create(tt.args.uow, tt.args.userId, tt.args.out)
			assert.NoError(t, err)
		})
	}
}

func Test_taskRepository_GetTasks(t *testing.T) {
	mock, db := NewMock()
	defer mock.Close()

	type args struct {
		uow       *UnitOfWork
		startDate string
		endDate   string
		userId    int
	}
	tests := []struct {
		name     string
		args     args
		wantTask *models.Task
		wantErr  bool
	}{
		{
			"",
			args{
				uow:       &UnitOfWork{Db: mock},
				startDate: "2020-12-29",
				endDate:   "2021-12-29",
				userId:    1,
			},
			&models.Task{
				Id:          1,
				Name:        "Task1",
				Description: "This is first task",
				ZoomLink:    "abcd",
				MeetLink:    "abcd",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			task := taskRepository{}
			query := `SELECT id, name, description, zoom_link, meet_link from task where start_date >= \? and end_date <= \? and user_id = \?`

			rows := db.NewRows([]string{"id", "name", "description", "zoom_link", "meet_link"}).
				AddRow(tt.wantTask.Id, tt.wantTask.Name, tt.wantTask.StartDate, tt.wantTask.ZoomLink, tt.wantTask.MeetLink)
			db.ExpectQuery(query).WithArgs(tt.args.startDate, tt.args.endDate, tt.args.userId).WillReturnRows(rows)
			tsk, err := task.GetTasks(tt.args.uow, tt.args.startDate, tt.args.endDate, tt.args.userId)
			assert.NotNil(t, tsk)
			assert.NoError(t, err)
		})
	}
}

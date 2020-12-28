package repository

import (
	"fmt"

	"Assignment/models"
)

type taskRepository struct {
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

type TaskRepository interface {
	Create(uow *UnitOfWork, userId int, out *models.Task) error
	GetTasks(uow *UnitOfWork, startDate string, endDate string, userId int) (task *models.Task, err error)
}

func (t *taskRepository) Create(uow *UnitOfWork, userId int, out *models.Task) error {

	query := fmt.Sprintf(`
			INSERT INTO task
			(
				id,
				user_id,
				name,
				description,
				start_date,
				end_date,
				zoom_link,
				meet_link
			)
			VALUES
			(
				?,
				?,
				?,
				?,
				?,
				?,
				?,
				?
			)`,
	)

	stmt, err := uow.Db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(&out.Id, &userId, &out.Name, &out.Description, &out.StartDate, &out.EndDate, &out.ZoomLink, &out.MeetLink)
	if err != nil {
		return err
	}
	return nil
}

func (t taskRepository) GetTasks(uow *UnitOfWork, startDate string, endDate string, userId int) (task *models.Task, err error) {
	query := fmt.Sprintf(`
									SELECT
										id,
										name,
										description,
										zoom_link,
										meet_link
									from
										task
									where
										start_date >= ? and end_date <= ? and user_id = ?`,
	)
	err = uow.Db.QueryRow(query, startDate, endDate, userId).Scan(&task.Id, &task.Name, &task.Description, &task.ZoomLink, task.MeetLink)
	if err != nil {
		return nil, err
	}
	return
}

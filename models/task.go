package models

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Task struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	ZoomLink    string `json:"zoom_link"`
	MeetLink    string `json:"meet_link"`
}

func NewTask(logger *logrus.Logger, name string, description string) (*Task, error) {
	if len(name) == 0 || len(description) == 0 {
		err := errors.New("name and password are required")
		logger.Error(err)
		return nil, err
	}

	task := &Task{
		Name:        name,
		Description: description,
	}
	return task, nil
}

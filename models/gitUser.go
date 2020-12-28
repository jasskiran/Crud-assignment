package models

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Github struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
}

func NewGitUser(logger *logrus.Logger, name string) (*Github, error){
	if len(name) == 0 {
		err := errors.New("name are required")
		logger.Error(err)
		return nil, err
	}

	user := &Github{
		UserName:        name,
	}

	return user, nil
}
package models

import (
	"errors"
	"log"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// user struct containing user fields
type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	GitUsername string `json:"git_username"`
}

func NewUser(logger *logrus.Logger, name string, password string, gitUser string) (*User, error) {
	if len(name) == 0 || len(password) == 0 {
		err := errors.New("name and password are required")
		logger.Error("name and password are required")
		return nil, err
	}

	// convert string password to hash
	hash := hashPassword([]byte(password))
	user := &User{
		Name:        name,
		Password:    hash,
		GitUsername: gitUser,
	}
	return user, nil
}

// converts string password to hash
func hashPassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func Authenticate(logger *logrus.Logger, name string, password string) error{
	if len(name) == 0 || len(password) == 0 {
		err := errors.New("name and password are required")
		logger.Error(err)
		return  err
	}
	return nil
}

func (user *User) Update(logger *logrus.Logger, name string) error {
	return nil
}

// compare hash password and string password
func CompareHashAndPassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

type Auth struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	AuthUUID string `json:"auth_uuid"`
}

//
type AuthDetails struct {
	AuthUUID string
	UserId   int64
}

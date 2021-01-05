package models

import (
	"errors"
	"log"
	"strings"

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

func Authenticate(logger *logrus.Logger, name string, password string) error {
	if len(name) == 0 || len(password) == 0 {
		err := errors.New("name and password are required")
		logger.Error(err)
		return err
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

func CheckTokenValidation(auth string, token string) (bool, error) {

	splitted := strings.Split(auth, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		err := errors.New("name and password are required")
		return false, err
	}
	tokenPart := splitted[1] //Grab the token part, what we are truly interested in

	if tokenPart == token {
		return true, nil
	}
	return false, nil
}

type Auth struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
	Active bool   `json:"active"`
}

//
type AuthDetails struct {
	AuthUUID string
	UserId   int64
}

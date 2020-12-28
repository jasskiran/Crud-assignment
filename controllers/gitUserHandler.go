package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Assignment/models"
	"Assignment/repository"
	"github.com/sirupsen/logrus"
)

type GitUserController struct {
	uow  *repository.UnitOfWork
	gitUserRepo repository.GitUserRepository
	Logger *logrus.Logger
}


func NewGitUserController(uow  *repository.UnitOfWork, gitUserRepo repository.GitUserRepository, Logger *logrus.Logger) *GitUserController{
	return &GitUserController{
		uow:      uow,
		gitUserRepo: gitUserRepo,
		Logger:   Logger,
	}
}

type gitUserDTO struct{
	Username string `json:"username"`
}

// CreateUserGithub creates a new user on github
func (controller *GitUserController) CreateGithubUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var requestDto gitUserDTO
	err := json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := models.NewGitUser(controller.Logger, requestDto.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.gitUserRepo.CreateGitUser(controller.uow, user)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}


// get github details of the current logged in user
func (controller *GitUserController)GetUserGithubDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	fmt.Println("userId", userId)
	gituser, err := controller.gitUserRepo.GetGitUser(controller.uow, userId)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(gituser)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
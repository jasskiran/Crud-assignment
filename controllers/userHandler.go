package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Assignment/models"
	"Assignment/repository"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	uow      *repository.UnitOfWork
	userRepo repository.UserRepository
	Logger   *logrus.Logger
}

func NewUserController(uow *repository.UnitOfWork, userRepo repository.UserRepository, Logger *logrus.Logger) *UserController {
	return &UserController{
		uow:      uow,
		userRepo: userRepo,
		Logger:   Logger,
	}
}

type userDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	GitUser  string `json:"git_user"`
}

// RegisterUser creates a new user
func (controller *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var requestDto userDTO
	err := json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := models.NewUser(controller.Logger, requestDto.Name, requestDto.Password, requestDto.GitUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.userRepo.Add(controller.uow, user)
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

// signIn
func (controller *UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var requestDTO struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestDTO)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = models.Authenticate(controller.Logger, requestDTO.Name, requestDTO.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := controller.userRepo.Login(controller.uow, requestDTO.Name)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v", u)
	fmt.Printf("%+v", requestDTO.Password)

	// compare password with hash string
	//err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(use.Password))
	err = models.CompareHashAndPassword(u.Password, requestDTO.Password)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// generate token
	tok, err := GenerateJWT(u.Id)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var response struct {
		Token string `json:"token"`
	}
	response.Token = tok

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// update user's profile
func (controller *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = controller.userRepo.Update(controller.uow, &user, userId)
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

//get details of the current logged in user
func (controller *UserController) GetCurrentUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	user, err := controller.userRepo.GetLoggedInUser(controller.uow, userId)
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

//SignOut - log out the current user
/*func (controller *UserController) SignOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId


	delErr := DeleteAuth(svc.DbSvc, token.UserId, token.AuthUUID)
	if delErr != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.StatusText(http.StatusOK)
}*/

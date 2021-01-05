package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Assignment/models"
	"Assignment/repository"
	"github.com/sirupsen/logrus"
)

// user wrapper
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

// request struct for register user
type userDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	GitUser  string `json:"git_user"`
}

/*
RegisterUser creates a new user
POST: http://localhost:8000/signin
Content-Type: application/json
Request:
{
  "name":  "user",
  "password":  "password"
	"git_username": "user"
}
*/
func (controller *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var requestDto userDTO
	err := json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// validating details for creating new user
	user, err := models.NewUser(controller.Logger, requestDto.Name, requestDto.Password, requestDto.GitUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// adding new user to database
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

/*
SignIn function is used to sign in with name and password
POST: http://localhost:8000/signin
Content-Type: application/json
Request:
{
  "name":  "user",
  "password":  "password"
}
Response:
{
token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk3NTgxODcsInVzZXJfaWQiOjd9.yUCEeSqxa5WKQZXxtAul_bmi1G0gw5BWWPCaw4QdzV4"
}
*/
func (controller *UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var requestDTO struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// decoding request body
	err := json.NewDecoder(r.Body).Decode(&requestDTO)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// authenticating request body to check if name and password are not nil
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

	// compare password with hash string
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

	//fmt.Println("config token",viper.Get("token.jwt_token"))

	var authentictaion models.Auth
	authentictaion.UserId = u.Id
	authentictaion.Token = tok

	err = controller.userRepo.AddToken(controller.uow, &authentictaion)
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

/*
UpdateProfile function update user's profile
PUT: http://localhost:8000/profile
Content-Type: application/json
Request:
{
  "name":  "user"
}

Response:
{
"id": 1,
"name": "user",
"password": "vgdxwnwmxswZPQLZ"
"git_username": "user1"
}
*/
func (controller *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	auth := r.Header.Get("Authorization")

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	tk, err := controller.userRepo.GetToken(controller.uow, userId)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	exists, err := models.CheckTokenValidation(auth, tk.Token)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
		return
	}

	if exists {
		var user models.User
		// decoding request body
		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			controller.Logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// updating user name in user table
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
	w.WriteHeader(http.StatusBadRequest)
}

/*
GetCurrentUserDetails function get details of the current logged in user
GET: http://localhost:8000/profile

Response:
{
"id": 1,
"name": "user",
"password": "vgdxwnwmxswZPQLZ"
"git_username": "user1"
}
*/
func (controller *UserController) GetCurrentUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	auth := r.Header.Get("Authorization")
	fmt.Println("auth", auth)
	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	fmt.Println(token)

	tk, err := controller.userRepo.GetToken(controller.uow, userId)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	exists, err := models.CheckTokenValidation(auth, tk.Token)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
		return
	}

	if exists {
		// get user details from database
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
	w.WriteHeader(http.StatusBadRequest)
}

/*
SignOut function is used to sign out the user by destroying its token
POST: http://localhost:8000/signout
*/
func (controller *UserController) SignOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	auth := r.Header.Get("Authorization")

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	tk, err := controller.userRepo.GetToken(controller.uow, userId)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	exists, err := models.CheckTokenValidation(auth, tk.Token)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
		return
	}

	if exists {
		// updating user name in user table
		err := controller.userRepo.DeleteToken(controller.uow, userId)
		if err != nil {
			controller.Logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.StatusText(http.StatusOK)
	}
	w.WriteHeader(http.StatusBadRequest)
}

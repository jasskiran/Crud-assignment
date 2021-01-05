package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Assignment/models"
	"Assignment/repository"
	"github.com/sirupsen/logrus"
)

const TimeLayout = "2006-01-02T15:04:05.000Z"

// task wrapper
type TaskController struct {
	uow      *repository.UnitOfWork
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
	Logger   *logrus.Logger
}

func NewTaskController(uow *repository.UnitOfWork, taskRepo repository.TaskRepository, userRepo repository.UserRepository, Logger *logrus.Logger) *TaskController {
	return &TaskController{
		uow:      uow,
		taskRepo: taskRepo,
		userRepo: userRepo,
		Logger:   Logger,
	}
}

// CreateTask creates a new task
func (controller *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
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
		var requestDto models.Task
		err := json.NewDecoder(r.Body).Decode(&requestDto)
		if err != nil {
			controller.Logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse(TimeLayout, requestDto.StartDate)
		if err != nil {
			panic(err)
		}
		endDate, err := time.Parse(TimeLayout, requestDto.EndDate)
		if err != nil {
			panic(err)
		}

		// In YY-MM-DD
		sdate := startDate.Format("2006-01-02")
		edate := endDate.Format("2006-01-02")

		task, err := models.NewTask(controller.Logger, requestDto.Name, requestDto.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		task.StartDate = sdate
		task.EndDate = edate
		err = controller.taskRepo.Create(controller.uow, userId, task)
		if err != nil {
			controller.Logger.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			controller.Logger.Error("err.Error()")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

// get all tasks of the user between start date and end date
func (controller *TaskController) GetTasks(w http.ResponseWriter, r *http.Request) {
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
		startDate := r.URL.Query().Get("start_date")
		endDate := r.URL.Query().Get("end_date")
		//startDate := r.FormValue("start_date")
		//endDate := r.FormValue("end_date")
		fmt.Printf("%+v", startDate)
		fmt.Printf("%+v", endDate)

		//sd, _ := time.Parse(LayoutISO, startDate)
		//ed, _ := time.Parse(LayoutISO, endDate)

		if len(startDate) == 0 || len(endDate) == 0 {
			http.Error(w, http.StatusText(http.StatusLengthRequired), http.StatusLengthRequired)
			return
		}

		task, err := controller.taskRepo.GetTasks(controller.uow, startDate, endDate, userId)
		if err != nil {
			controller.Logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			controller.Logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

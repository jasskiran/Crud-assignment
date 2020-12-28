package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"Assignment/models"
	"Assignment/repository"
	"github.com/sirupsen/logrus"
)

const (
	LayoutISO = "2006-01-02"
)

type TaskController struct {
	uow  *repository.UnitOfWork
	taskRepo repository.TaskRepository
	Logger *logrus.Logger
}

func NewTaskController(uow  *repository.UnitOfWork, taskRepo repository.TaskRepository, Logger *logrus.Logger) *TaskController{
	return &TaskController{
		uow:      uow,
		taskRepo: taskRepo,
		Logger:   Logger,
	}
}

type taskDTO struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	ZoomLink    *string   `json:"zoom_link"`
	MeetLink    *string   `json:"meet_link"`
}
// CreateTask creates a new task
func (controller *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	var requestDto taskDTO
	err := json.NewDecoder(r.Body).Decode(&requestDto)
	if err != nil {
		controller.Logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	task, err := models.NewTask(controller.Logger, requestDto.Name, requestDto.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	task.StartDate = time.Now()
	task.EndDate = time.Now()

	err = controller.taskRepo.Create(controller.uow, userId, task)
	if err != nil {
		controller.Logger.Error("query error")
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

// get all tasks of the user between start date and end date
func (controller *TaskController) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//get authenticated from the request context
	token := r.Context().Value("userId").(Token)
	userId := token.UserId

	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

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

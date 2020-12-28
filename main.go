package main

import (
	"database/sql"
	"net/http"
	"os"

	"Assignment/controllers"
	"Assignment/repository"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main(){
	lvl, _ := logrus.ParseLevel("trace")
	logger := newLogger(lvl)

	db, err := sql.Open("mysql", "root:password@tcp(localhost:3307)/test")
	if err != nil{
		logger.Fatal(err)
	}
	defer db.Close()
	router := initialiseRoutes(db, logger)

	logger.Infof("starting server")
	logger.Fatal(http.ListenAndServe(":8000", router))

}

func initialiseRoutes(db *sql.DB, logger *logrus.Logger ) *mux.Router{

	router := mux.NewRouter().StrictSlash(true)

	uow := repository.NewUnitOfWork(db)

	userRepository := repository.NewUserRepository()
	userController := controllers.NewUserController(uow, userRepository, logger)

	router.HandleFunc("/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/signin", userController.SignIn).Methods("POST")
	router.HandleFunc("/profile", controllers.AuthRequired(userController.GetCurrentUserDetails)).Methods("GET")
	router.HandleFunc("/profile", controllers.AuthRequired(userController.UpdateProfile)).Methods("PUT")
	//router.HandleFunc("/signout", controllers.AuthRequired(svc.SignOut)).Methods("POST")

	gitUserRepository := repository.NewGitUserRepository()
	gitUserController := controllers.NewGitUserController(uow, gitUserRepository, logger)

	router.HandleFunc("/github", gitUserController.CreateGithubUser).Methods("POST")
	router.HandleFunc("/github", controllers.AuthRequired(gitUserController.GetUserGithubDetails)).Methods("GET")

	taskRepository := repository.NewTaskRepository()
	taskController := controllers.NewTaskController(uow, taskRepository, logger)

	router.HandleFunc("/task", controllers.AuthRequired(taskController.CreateTask)).Methods("POST")
	router.HandleFunc("/task", controllers.AuthRequired(taskController.GetTasks)).Methods("GET")

	return router
}


func newLogger(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.Level = level
	logger.Out = os.Stdout
	return logger
}
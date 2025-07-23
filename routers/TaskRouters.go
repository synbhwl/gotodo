package routers 

import (
	"gotodo/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handlers.Hello).Methods("GET")
	r.HandleFunc("/tasks/new", handlers.AddTodo).Methods("POST")
	r.HandleFunc("/tasks/all", handlers.SeeTasks).Methods("GET")
	r.HandleFunc("/tasks/search", handlers.SearchTask).Methods("GET")
	r.HandleFunc("/tasks/delete", handlers.DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/edit", handlers.EditTask).Methods("PATCH")
	r.HandleFunc("/tasks/filter", handlers.FilterByStatus).Methods("GET")
	return r
}

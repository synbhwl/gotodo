package main

import (
	"gotodo/handlers"
	"net/http"
	"github.com/gorilla/mux"
)

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/hello", handlers.Hello).Methods("GET")
	r.HandleFunc("/tasks/new", handlers.AddTodo).Methods("POST")
	r.HandleFunc("/tasks/all", handlers.SeeTasks).Methods("GET")
	r.HandleFunc("/tasks/search", handlers.SearchTask).Methods("GET")
	r.HandleFunc("/tasks/delete", handlers.DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/edit", handlers.EditTask).Methods("PATCH")
	r.HandleFunc("/tasks/filter", handlers.FilterByStatus).Methods("GET")
	http.ListenAndServe(":3000", r)
}
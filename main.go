package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Task struct { // making a type Task that'll hold all those fields
	Title string `json:"title"`
	Id string `json:"id"`
	Completed bool `json:"completed"`
	Timestamp string `json:"timestamp"`
}

func hello(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	msg := map[string]string{"message":"the server says hello"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func addTodo(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{ //checking method
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	var newtask Task //new task with type Task
	err := json.NewDecoder(r.Body).Decode(&newtask) //itll parse the json title into the title field of newtask
	if err != nil{
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return 
	}

	 //making a new task of type Task that holds all those fields
	newtask.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	newtask.Completed = false
	newtask.Timestamp = time.Now().Format(time.RFC3339)


	file, ferr := os.OpenFile("tasks.json", os.O_RDWR | os.O_CREATE, 0644) //opening a file (making it a reader)
	if ferr != nil{
		http.Error(w, "couldnt open task file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	filebytes, rerr := io.ReadAll(file) // reading from the file and getting a slice of bytes
	if rerr != nil{
		http.Error(w, "couldnt read task file", http.StatusInternalServerError)
		return
	}

	var tasks []Task //building a slice of type Task that'll hold all the tasks
	jsonerr := json.Unmarshal(filebytes, &tasks) //parsing the bytes into the tasks slice 
	if jsonerr != nil{
		http.Error(w, "couldnt parse the json file", http.StatusInternalServerError)
		return
	}
	
	tasks = append(tasks, newtask) //appending new task to the tasks slice 

	updated, mrshlerr := json.MarshalIndent(tasks, "", " ") //parsing the go variable into json
	if mrshlerr != nil{
		http.Error(w, "couldnt parse the json file", http.StatusInternalServerError)
		return
	}

	werr := os.WriteFile("tasks.json", updated, 0644) //overwriting the file with json
	if werr != nil{
		http.Error(w, "couldnt parse the json file", http.StatusInternalServerError)
		return 
	}

	msg := map[string]string{"message":"task added successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)

}

func seeTasks(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{ 
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	file, ferr := os.Open("tasks.json") 
	if ferr != nil{
		http.Error(w, "couldnt open task file", http.StatusInternalServerError)
		return
	}

	filebytes, rerr := io.ReadAll(file)
	if rerr != nil{
		http.Error(w, "couldnt read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(filebytes)

}

func searchTask(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{ 
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	title := r.URL.Query().Get("title")

	file, ferr := os.Open("tasks.json") 
	if ferr != nil{
		http.Error(w, "couldnt open task file", http.StatusInternalServerError)
		return
	}

	filebytes, rerr := io.ReadAll(file)
	if rerr != nil{
		http.Error(w, "couldnt read file", http.StatusInternalServerError)
		return
	}

	var tasks []Task
	jmerr := json.Unmarshal(filebytes, &tasks)
	if jmerr != nil{
		http.Error(w, "couldnt marshall file bytes", http.StatusInternalServerError)
		return 
	}

	var matched []Task
	for _, task := range tasks{
		if strings.ToLower(task.Title) == strings.ToLower(title){
			matched = append(matched, task)
		}
	}

	matchedbyte, jumerr := json.MarshalIndent(matched, "", " ")
	if jumerr != nil{
		http.Error(w, "couldnt unmarshall data", http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(matchedbyte)
}

func main(){
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/tasks/new", addTodo)
	http.HandleFunc("/tasks/all", seeTasks)
	http.HandleFunc("/tasks/search", searchTask)
	http.ListenAndServe(":3000", nil)
}
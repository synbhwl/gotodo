package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type reqbody struct{
	Title string `json:"task"`
}

type Task struct{
	id string
	title string
	completed bool
	timestamp string
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

	type Data struct { //making a struct that'll hold the title from req body
		Title string `json:"title"` //check title
	}

	type Task struct { // making a type Task that'll hold all those fields
		Title string
		Id string
		Completed bool
		Timestamp string
	}

	var newData Data //new object of type Data
	err := json.NewDecoder(r.Body).Decode(&newData) //feeding it into new data so that it takes the title
	if err != nil{
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return 
	}

	newtask := Task{ //making a new task of type Task that holds all those fields
		Title : newData.Title,
		Id : fmt.Sprintf("%d", time.Now().UnixNano()),
		Completed : false,
		Timestamp: time.Now().Format(time.RFC3339),
	}


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

func main(){
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/tasks/new", addTodo)
	http.ListenAndServe(":3000", nil)
}
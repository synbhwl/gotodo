package handlers 

import (
	"net/http"
	"encoding/json"
	"os"
	"fmt"
	"time"
	"io"
	"strings"
)

type Task struct { // making a type Task that'll hold all those fields
	Title string `json:"title"`
	Id string `json:"id"`
	Completed bool `json:"completed"`
	Timestamp string `json:"timestamp"`
}

func Hello(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	msg := map[string]string{"message":"the server says hello"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// adding to do 
func AddTodo(w http.ResponseWriter, r *http.Request){
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

// seeing all tasks handler 
func SeeTasks(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{ 
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	file, ferr := os.Open("tasks.json") 
	if ferr != nil{
		http.Error(w, "couldnt open task file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	filebytes, rerr := io.ReadAll(file)
	if rerr != nil{
		http.Error(w, "couldnt read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(filebytes)

}

// seeing one task handler 
func SearchTask(w http.ResponseWriter, r *http.Request){
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

	defer file.Close()

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

// delete a task 
func DeleteTask(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete{ 
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	title := r.URL.Query().Get("title")


	file, ferr := os.OpenFile("tasks.json", os.O_RDWR, 0644)
	if ferr != nil{
		http.Error(w, "couldnt read file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	filebytes, rerr := io.ReadAll(file)
	if rerr != nil{
		http.Error(w, "couldnt read file bytes", http.StatusInternalServerError)
	}

	var tasks []Task
	jumerr := json.Unmarshal(filebytes, &tasks)
	if jumerr != nil{
		http.Error(w, "couldnt unmarshall", http.StatusInternalServerError)
		return
	}

	var matched []Task
	for _, task := range tasks{
		if strings.ToLower(task.Title) != strings.ToLower(title){
			matched = append(matched, task)
		}
	}

	matchedbytes, jmerr := json.MarshalIndent(matched, "", " ")
	if jmerr != nil{
		http.Error(w, "couldnt unmarshall", http.StatusInternalServerError)
		return
	}

	werr := os.WriteFile("tasks.json", matchedbytes, 0644)
	if werr != nil{
		http.Error(w, "couldnt update file", http.StatusInternalServerError)
		return
	}

	msg := map[string]string{"message":"task deleted successfully"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// edit task
func EditTask(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPatch{ 
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	oldtitle := r.URL.Query().Get("old_title")
	newtitle := r.URL.Query().Get("new_title")
	completed := r.URL.Query().Get("completed")

	if oldtitle == ""{
		http.Error(w, "empty title param", http.StatusBadRequest)
		return
	}

	file, ferr := os.OpenFile("tasks.json", os.O_RDWR, 0644)
	if ferr != nil{
		http.Error(w, "couldnt open file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	
	filebytes, rerr := io.ReadAll(file)
	if rerr != nil{
		http.Error(w, "couldnt read file", http.StatusInternalServerError)
		return
	}

	var tasks []*Task
	jumerr := json.Unmarshal(filebytes, &tasks)
	if jumerr != nil{
		http.Error(w, "couldnt unmarshal bytes", http.StatusInternalServerError)
		return
	}


	for _, task := range tasks{
		if strings.ToLower(oldtitle) == strings.ToLower(task.Title){
			if newtitle != ""{
				task.Title = newtitle
			}
			if completed != ""{
				status := completed == "true"
				task.Completed = status
			}
		} 
	}

	updated, jmerr := json.MarshalIndent(tasks, "", " ")
	if jmerr != nil{
		http.Error(w, "couldnt marshal into json bytes", http.StatusInternalServerError)
		return
	}

	werr := os.WriteFile("tasks.json", updated, 0644)
	if werr != nil{
		http.Error(w, "couldnt write file", http.StatusInternalServerError)
		return
	}

	msg := map[string]string{"message":"task updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// filter by status
func FilterByStatus(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("completed") == "true"

	file, errFileOpen := os.OpenFile("tasks.json", os.O_RDWR, 0644)
	if errFileOpen != nil{
		http.Error(w, "cannot open file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	filebytes, errFileRead := io.ReadAll(file)
	if errFileRead != nil{
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	var tasks []Task
	errUnmarshal := json.Unmarshal(filebytes, &tasks)
	if errUnmarshal != nil{
		http.Error(w, "cannot unmarshal file bytes", http.StatusInternalServerError)
		return
	}

	var matched []Task
	for _, task := range tasks{
		if task.Completed == status{
			matched = append(matched, task)
		}
	}

	matchedbytes, errMarshal := json.MarshalIndent(matched, "", " ")
	if errMarshal != nil{
		http.Error(w, "cannot unmarshal file bytes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(matchedbytes)
}
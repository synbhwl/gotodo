package handlers

import (
	"encoding/json"
	"fmt"
	"gotodo/helpers"
	"net/http"
	"strings"
	"time"
)

type Task struct { 
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
	var newtask *helpers.Task
	err := json.NewDecoder(r.Body).Decode(&newtask) 
	if err != nil{
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return 
	}

	newtask.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	newtask.Completed = false
	newtask.Timestamp = time.Now().Format(time.RFC3339)

	tasks, errFileRead := helpers.ReadTasks("tasks.json")
	if errFileRead != nil{
		http.Error(w, "file related issue", http.StatusInternalServerError)
		return
	}
	
	tasks = append(tasks, newtask) 

	errFileWrite := helpers.WriteTasks("tasks.json", tasks)
	if errFileWrite != nil{
		http.Error(w, "file related issue", http.StatusInternalServerError)
		return
	}

	msg := map[string]string{"message":"task added successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)

}

// seeing all tasks handler 
func SeeTasks(w http.ResponseWriter, r *http.Request){
	tasks, errFile := helpers.ReadTasks("tasks.json")
	if errFile != nil{
		http.Error(w, "issue with file", http.StatusInternalServerError)
		return
	}

	updated, errMarshal := json.MarshalIndent(tasks, "", " ")
	if errMarshal != nil{
		http.Error(w, "issue marshalling", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(updated)

}

// seeing one task handler 
func SearchTask(w http.ResponseWriter, r *http.Request){
	title := r.URL.Query().Get("title")

	tasks, errFile := helpers.ReadTasks("tasks.json")
	if errFile != nil{
		http.Error(w, "issue with file", http.StatusInternalServerError)
		return
	}


	var matched []*helpers.Task
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
	title := r.URL.Query().Get("title")

	tasks, errFile := helpers.ReadTasks("tasks.json")
	if errFile != nil{
		http.Error(w, "issue with file", http.StatusInternalServerError)
		return
	}


	var matched []*helpers.Task
	for _, task := range tasks{
		if strings.ToLower(task.Title) != strings.ToLower(title){
			matched = append(matched, task)
		}
	}

	errFileWrite := helpers.WriteTasks("tasks.json", matched)
	if errFileWrite != nil{
		http.Error(w, "issue with updating file", http.StatusInternalServerError)
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

	tasks, errFileRead := helpers.ReadTasks("tasks.json")
	if errFileRead != nil{
		http.Error(w, "issues with file", http.StatusInternalServerError)
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

	errWriteFile := helpers.WriteTasks("tasks.json", tasks)
	if errWriteFile != nil{
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

	tasks, errFileRead := helpers.ReadTasks("tasks.json")
	if errFileRead != nil{
		http.Error(w, "issues with file", http.StatusInternalServerError)
		return
	}

	var matched []*helpers.Task
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
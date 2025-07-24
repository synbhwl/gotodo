package handlers

import (
	"encoding/json"
	"fmt"
	"gotodo/helpers"
	"log"
	"net/http"
	"strings"
	"time"
)

// Hello sends a starting message to the user, can also be used as to check if the router is working fine
func Hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	log.Println("hello was sent")

	msg := map[string]string{"message": "the server says hello"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// AddTodo adds the user's POST request's todo title along with other details that the server decides like ID, timestamps, and status
func AddTodo(w http.ResponseWriter, r *http.Request) {
	var newtask *helpers.Task
	err := json.NewDecoder(r.Body).Decode(&newtask)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	newtask.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	newtask.Completed = false
	newtask.Timestamp = time.Now().Format(time.RFC3339)

	tasks, errFileRead := helpers.ReadTasks("tasks.json")
	if errFileRead != nil {
		http.Error(w, "file related issue", http.StatusInternalServerError)
		return
	}

	tasks = append(tasks, newtask)

	errFileWrite := helpers.WriteTasks("tasks.json", tasks)
	if errFileWrite != nil {
		http.Error(w, "file related issue", http.StatusInternalServerError)
		return
	}

	log.Println("new todo was added to the list")
	msg := map[string]string{"message": "task added successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)

}

// SeeTasks lets the user see all the tasks in the tasks.json file (problem: it doesnt have auth, so everybody has access to all the tasks)
func SeeTasks(w http.ResponseWriter, r *http.Request) {
	tasks, errFile := helpers.ReadTasks("tasks.json")
	if errFile != nil {
		http.Error(w, "issue with file", http.StatusInternalServerError)
		return
	}

	updated, errMarshal := json.MarshalIndent(tasks, "", " ")
	if errMarshal != nil {
		http.Error(w, "issue marshalling", http.StatusInternalServerError)
		return
	}

	log.Println("list of all tasks sent")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(updated)
}

// SearchTask lets the user search for a single task by title in isolation. the user gets a json object with all the details in the specific task like id, time and status
func SearchTask(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "title query param missing", http.StatusBadRequest)
		return
	}

	tasks, errFile := helpers.ReadTasks("tasks.json")
	if errFile != nil {
		http.Error(w, "issue with file", http.StatusInternalServerError)
		return
	}

	var matched []*helpers.Task
	for _, task := range tasks {
		if strings.ToLower(task.Title) == strings.ToLower(title) {
			matched = append(matched, task)
		}
	}

	matchedbyte, jumerr := json.MarshalIndent(matched, "", " ")
	if jumerr != nil {
		http.Error(w, "couldnt unmarshall data", http.StatusInternalServerError)
		return
	}

	log.Println("specific task was sent to user")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(matchedbyte)
}

// DeleteTask lets the user delete a specific task by title (it may also delete more than one task if there are duplicate titles)
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "title query param missing", http.StatusBadRequest)
		return
	}

	tasks, errFile := helpers.ReadTasks("tasks.json")
	if errFile != nil {
		http.Error(w, "issue with file", http.StatusInternalServerError)
		return
	}

	var matched []*helpers.Task
	for _, task := range tasks {
		if strings.ToLower(task.Title) != strings.ToLower(title) {
			matched = append(matched, task)
		}
	}

	errFileWrite := helpers.WriteTasks("tasks.json", matched)
	if errFileWrite != nil {
		http.Error(w, "issue with updating file", http.StatusInternalServerError)
		return
	}

	log.Println("specific task was deleted as per request")
	msg := map[string]string{"message": "task deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// EditTask lets the user edit a specific task in terms of two things, title and status. old_status is mendatory, but the other two query params i.e new_title and completed are optional depending on the preference of the user. if either of the optional params are missing, the field stays same. in case of two missing, nothing changes, but doesnt return an error
func EditTask(w http.ResponseWriter, r *http.Request) {
	oldtitle := r.URL.Query().Get("old_title")
	newtitle := r.URL.Query().Get("new_title")
	completed := r.URL.Query().Get("completed")

	if oldtitle == "" {
		http.Error(w, "empty title param", http.StatusBadRequest)
		return
	}

	tasks, errFileRead := helpers.ReadTasks("tasks.json")
	if errFileRead != nil {
		http.Error(w, "issues with file", http.StatusInternalServerError)
		return
	}

	for _, task := range tasks {
		if strings.ToLower(oldtitle) == strings.ToLower(task.Title) {
			if newtitle != "" {
				task.Title = newtitle
			}
			if completed != "" {
				status := completed == "true"
				task.Completed = status
			}
		}
	}

	errWriteFile := helpers.WriteTasks("tasks.json", tasks)
	if errWriteFile != nil {
		http.Error(w, "couldnt write file", http.StatusInternalServerError)
		return
	}

	log.Println("specific task was edited as per request")
	msg := map[string]string{"message": "task updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// FilterByStatus lets the user filter tasks on the basis of completed field, the user has two choices, true or false, the function lets the user see the tasks that matches their filter boolean
func FilterByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("completed") == "true"

	tasks, errFileRead := helpers.ReadTasks("tasks.json")
	if errFileRead != nil {
		http.Error(w, "issues with file", http.StatusInternalServerError)
		return
	}

	var matched []*helpers.Task
	for _, task := range tasks {
		if task.Completed == status {
			matched = append(matched, task)
		}
	}

	matchedbytes, errMarshal := json.MarshalIndent(matched, "", " ")
	if errMarshal != nil {
		http.Error(w, "cannot unmarshal file bytes", http.StatusInternalServerError)
		return
	}

	log.Println("filtered list of tasks by status was sent to user")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(matchedbytes)
}

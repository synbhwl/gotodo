package helpers

import (
	"encoding/json"
	"io"
	"os"
)

// Task struct lays out the essential fields/structure of a task object, all tasks are therefore objects of this struct. since only one struct is there, it lives in the file helpers instead of another helper file
type Task struct {
	Title     string `json:"title"`
	Id        string `json:"id"`
	Completed bool   `json:"completed"`
	Timestamp string `json:"timestamp"`
}

// ReadTasks is a file helper that opens the file, reads bytes, unmarshal it into a slice of pointers to the struct Task and returns the slice thereby
func ReadTasks(filename string) ([]*Task, error) {
	file, errFileOpen := os.OpenFile(filename, os.O_RDWR, 0644)
	if errFileOpen != nil {
		return nil, errFileOpen
	}

	defer file.Close()

	filebytes, errFileRead := io.ReadAll(file)
	if errFileRead != nil {
		return nil, errFileRead
	}

	var tasks []*Task
	errUnmarshal := json.Unmarshal(filebytes, &tasks)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return tasks, nil
}

// WriteTasks opens a file, marshals the slice of pointers into json bytes and thereby writes it into the file 
func WriteTasks(filename string, tasks []*Task) error {
	file, errFileOpen := os.OpenFile(filename, os.O_RDWR, 0644)
	if errFileOpen != nil {
		return errFileOpen
	}

	defer file.Close()

	updated, errMarshal := json.MarshalIndent(tasks, "", " ")
	if errMarshal != nil {
		return errMarshal
	}

	errFileWrite := os.WriteFile(filename, updated, 0644)
	if errFileWrite != nil {
		return errFileWrite
	}

	return nil
}

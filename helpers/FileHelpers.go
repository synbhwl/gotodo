package helpers

import (
	"encoding/json"
	"io"
	"os"
)

type Task struct { 
	Title string `json:"title"`
	Id string `json:"id"`
	Completed bool `json:"completed"`
	Timestamp string `json:"timestamp"`
}

func ReadTasks(filename string)([]*Task, error){
	file, errFileOpen := os.OpenFile(filename, os.O_RDWR, 0644)
	if errFileOpen != nil{
		return nil, errFileOpen
	}

	defer file.Close()

	filebytes, errFileRead := io.ReadAll(file)
	if errFileRead != nil{
		return nil, errFileRead
	}

	var tasks []*Task
	errUnmarshal := json.Unmarshal(filebytes, &tasks)
	if errUnmarshal != nil{
		return nil, errUnmarshal
	}

	return tasks, nil
}

func WriteTasks(filename string, tasks []*Task) error{
	file, errFileOpen := os.OpenFile(filename, os.O_RDWR, 0644)
	if errFileOpen != nil{
		return errFileOpen
	}

	defer file.Close()

	updated, errMarshal := json.MarshalIndent(tasks, "", " ")
	if errMarshal != nil{
		return errMarshal
	}

	errFileWrite := os.WriteFile(filename, updated, 0644)
	if errFileWrite != nil{
		return errFileWrite
	}

	return nil
}
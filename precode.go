package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTasks(writer http.ResponseWriter, request *http.Request) {
	tasks, err := json.Marshal(tasks)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write(tasks); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func addTask(writer http.ResponseWriter, request *http.Request) {
	var newTask Task

	err := json.NewDecoder(request.Body).Decode(&newTask)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[newTask.ID] = newTask
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func getTask(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	task, ok := tasks[id]
	if !ok {
		errorMsg := fmt.Sprintf("Task with id: %s not found", id)
		http.Error(writer, errorMsg, http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	err := json.NewEncoder(writer).Encode(task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteTask(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	_, ok := tasks[id]
	if !ok {
		errorMsg := fmt.Sprintf("Task with id: %s not found", id)
		http.Error(writer, errorMsg, http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	writer.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

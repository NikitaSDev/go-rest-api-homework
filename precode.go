package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// Ниже напишите обработчики для каждого эндпоинта
func getTasks(w http.ResponseWriter, r *http.Request) {

	resp, err := json.Marshal(tasks)
	if err != nil {
		fmt.Printf("ошибка сериализации: %s\n", err.Error())
		http.Error(w, "getting tasks error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Printf("ошибка отправки задачи: %s\n", err.Error())
		http.Error(w, "sending task error", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		fmt.Printf("ошибка десереализации: %s\n", err.Error())
		http.Error(w, "readibg task error", http.StatusBadRequest)
		return
	}

	if _, ok := tasks[task.ID]; ok {
		fmt.Printf("ошибка создания: задача уже существует")
		http.Error(w, "task already exist", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		fmt.Printf("задача с id %s не найдена\n", id)
		http.Error(w, "task not found", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		fmt.Printf("ошибка сериализации: %s\n", err.Error())
		http.Error(w, "getting task error", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	_, ok := tasks[id]
	if !ok {
		fmt.Printf("задача с id %s не найдена\n", id)
		http.Error(w, "task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/task", addTask)
	r.Get("/task/{id}", getTask)
	r.Delete("/task/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

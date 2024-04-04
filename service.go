package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TaskService struct {
	storage *Storage
}

func NewTaskService(storage *Storage) TaskService {
	return TaskService{storage: storage}
}

func (t TaskService) TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	tasksFromDB, _ := t.storage.SelectTasks()
	dtos := TasksToDto(tasksFromDB)
	response := ResponseTasks{Dtos: dtos}
	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}
	log.Println("[Info] Success: tasks from DB are given")
	w.Write(responseBody)
}

func (t TaskService) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		t.addTask(w, r)
	case http.MethodGet:
		t.getTask(w, r)
	case http.MethodPut:
		t.editTask(w, r)
	case http.MethodDelete:
		//removeTask(w, r)
	}
}

func (t TaskService) editTask(w http.ResponseWriter, r *http.Request) {
	var inDTO TaskInputDTO

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&inDTO); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("json Decoder:", err)

		return
	}

	task, err := DtoToTask(inDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("DtoToTask:", err)

		return
	}

	if task.ID == 0 {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)

		return
	}

	err = t.storage.UpdateTask(task)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("UpdateTask:", err)

		return
	}

	w.Write([]byte(`{}`))

}

func (t TaskService) getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("[WARN] Failed convertation:", err)
		return
	}

	log.Println("[By 1  Id] до перехода в сторидж ")

	task, err := t.storage.SelectById(int(id))
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}
	dto := TaskToDto(task)

	log.Println("[By 5  Id] дто отработало ")

	responseBody, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}

	w.Write(responseBody)
}

func (t TaskService) addTask(w http.ResponseWriter, r *http.Request) {
	var inDTO TaskInputDTO

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&inDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed json decoding:", err)
		return
	}

	task, err := DtoToTask(inDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed Dto-to-Task convertation:", err)
		return
	}

	log.Println("[DTO ] : " + inDTO.Date + inDTO.Title + inDTO.Comment + inDTO.Repeat)
	log.Println("[TASK] : " + task.Date + task.Title + task.Comment + task.Repeat)

	id, err := t.storage.InsertTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("[WARN] Failed to add a task:", err)
		return
	}

	log.Println("[Info] Success: Task added with id = " + strconv.Itoa(id))
	w.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, id)))
}

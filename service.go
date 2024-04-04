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

func (t TaskService) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		t.addTask(w, r)

	case http.MethodGet:
		//getTask(w, r)
	case http.MethodPut:
		//updateTask(w, r)
	case http.MethodDelete:
		//removeTask(w, r)
	}
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

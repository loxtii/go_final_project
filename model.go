package main

import (
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      int
	Date    string
	Title   string
	Comment string
	Repeat  string
}

// func (t Task) String() string {
// 	return fmt.Sprintf("ID: %d, Date: %s, Title: %s, Comment: %s, Repeat: %s", t.ID, t.Date, t.Title, t.Comment, t.Repeat)
// }

type TaskInputDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// type ResponseTasks struct {
// 	Tasks []TaskDTO `json:"tasks"`
// }

func DtoToTask(dto TaskInputDTO) (Task, error) {
	if dto.Title == "" {
		return Task{}, fmt.Errorf("empty title")
	}

	today := makeDate(time.Now())
	date := today

	var err error

	if dto.Date != "" {
		date, err = time.Parse("20060102", dto.Date)
		if err != nil {
			return Task{}, fmt.Errorf("invalid date format")
		}
	}

	if date.Before(today) {
		if dto.Repeat == "" {
			date = today
		} else {
			date, err = CalculateNextDate(date, today, dto.Repeat)
			if err != nil {
				return Task{}, fmt.Errorf("can't get next date: %w", err)
			}
		}
	}

	id := 0
	if dto.ID != "" {
		id, _ = strconv.Atoi(dto.ID)
	}

	return Task{
		ID:      int(id),
		Date:    date.Format("20060102"),
		Title:   dto.Title,
		Comment: dto.Comment,
		Repeat:  dto.Repeat,
	}, nil
}

func makeDate(datetime time.Time) time.Time {
	y, m, d := datetime.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

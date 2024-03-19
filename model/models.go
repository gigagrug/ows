package model

import (
	"time"
)

type User struct {
	ID          string       `json:"id"`
	Email       string       `json:"email"`
	Password    string       `json:"password"`
	CreatedAt   time.Time    `json:"createdAt"`
	Projects    []Project    `json:"projects"`
	LogServices []LogService `json:"logServices"`
}

type Project struct {
	ID         string       `json:"id"`
	Name       string       `json:"name" validate:"required"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
	LogService []LogService `json:"logServices"`
	UserID     string       `json:"userId"`
}

type LogService struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Logs      []Log     `json:"logs"`
	UserID    string    `json:"userId"`
	ProjectID string    `json:"projectId"`
}

type Log struct {
	ID           string    `json:"id"`
	Severity     string    `json:"severity"`
	Info         string    `json:"info"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
	LogServiceID string    `json:"logServiceId"`
	UserID       string    `json:"userId"`
}

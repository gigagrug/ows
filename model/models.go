package model

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	Projects  []Project `json:"projects"`
}

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    string    `json:"user_id"`
}

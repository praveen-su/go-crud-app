package main

import "time"

type User struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Task struct {
	ID     string `json:"id" gorm:"primaryKey"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type AuditLog struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Entity    string
	EntityID  string
	Action    string
	CreatedAt time.Time
}

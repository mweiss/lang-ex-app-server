package models

import "time"

type UserLanguage struct {
	Id         int64
	UserId     int64
	Language   string
	Level      string
	IsLearning bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

package models

import "time"

type UserLanguage struct {
	Id         int64
	UserId     int64
	language   string
	level      string
	isLearning bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

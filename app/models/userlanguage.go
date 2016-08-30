package models

import "time"

type UserLanguage struct {
	Id         int64      `json:"id"`
	UserId     int64      `json:"userId"`
	Language   string     `json:"language"`
	Level      string     `json:"level"`
	IsLearning bool       `json:"isLearning"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `json:"-"`
}

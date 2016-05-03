package models

import "time" // if you need/want

type Badge struct {
	Id        int64
	UserId    int64
	BadgeName string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time // for soft delete
}

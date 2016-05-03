package models

import "time" // if you need/want

type Post struct {
	Id        int64
	AuthorId  int64
	Title     string
	Body      string `sql:"size:4000"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

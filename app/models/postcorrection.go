package models

import "time" // if you need/want

type PostCorrection struct {
	Id        int64
	PostId    int64 `sql:"index"`
	AuthorId  int64
	Comment   string `sql:"size:4000"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	PostEdits []PostEdit
}

package models

import "time" // if you need/want

type PostEdit struct {
	Id               int64
	PostCorrectionId int64 `sql:"index"`
	Section          string
	StartIndex       int64
	Length           int64
	NewText          string
	Comment          string `sql:"size:4000"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

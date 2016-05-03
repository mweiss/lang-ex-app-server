package models

import "time"

type TestEntity struct {
	Id                int64
	Name              string
	EncryptedPassword []byte
	Password          string `sql:"-"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}

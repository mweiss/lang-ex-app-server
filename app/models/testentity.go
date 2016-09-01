package models

import "time"

type TestEntity struct {
	Id                uint
	Name              string
	EncryptedPassword []byte
	Password          string `sql:"-"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

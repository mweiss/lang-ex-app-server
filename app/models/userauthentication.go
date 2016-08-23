package models

import "time" // if you need/want

type UserAuthentication struct {
	Id            int64
	UserId        int64
	FacebookToken string
	LoginToken    string `sql:"unique_index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

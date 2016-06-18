package models

import "time" // if you need/want

type User struct {
	Id                     int64
	Name                   string
	SelfIntroductionPostId int64
	ImageURL               string
	EncryptedPassword      []byte
	Password               string `sql:"-"`
	KarmaScore             int64
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              time.Time // for soft delete

	UserLanguages []UserLanguage
}

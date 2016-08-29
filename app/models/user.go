package models

import "time" // if you need/want

type User struct {
	Id                        int64
	FacebookId                string `sql:"index"`
	Name                      string
	FirstName                 *string
	LastName                  *string
	Email                     string
	SelfIntroductionPostId    int64
	AllowNonNativeCorrections bool
	ImageURL                  string
	EncryptedPassword         []byte
	Password                  string `sql:"-"`
	KarmaScore                int64
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 *time.Time // for soft delete

	UserLanguages []UserLanguage
}

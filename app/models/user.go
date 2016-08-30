package models

import "time" // if you need/want

type User struct {
	Id                        int64      `json:"id"`
	FacebookId                string     `json:"facebookId" sql:"index"`
	Name                      string     `json:"name"`
	FirstName                 *string    `json:"firstName"`
	LastName                  *string    `json:"lastName"`
	Email                     string     `json:"email"`
	SelfIntroductionPostId    int64      `json:"selfIntroductionPostId"`
	AllowNonNativeCorrections bool       `json:"allowNonNativeCorrections"`
	ImageURL                  string     `json:"imageURL"`
	EncryptedPassword         []byte     `json:"-"`
	Password                  string     `json:"-" sql:"-"`
	KarmaScore                int64      `json:"id"`
	CreatedAt                 time.Time  `json:"-"`
	UpdatedAt                 time.Time  `json:"-"`
	DeletedAt                 *time.Time `json:"-"` // for soft delete

	UserLanguages []UserLanguage `json:"userLanguages"`
}

package models

type User struct {
	Model
	FacebookId                string  `json:"facebookId" sql:"index"`
	Name                      string  `json:"name"`
	FirstName                 *string `json:"firstName"`
	LastName                  *string `json:"lastName"`
	Email                     string  `json:"email"`
	SelfIntroductionPostId    uint    `json:"selfIntroductionPostId"`
	AllowNonNativeCorrections bool    `json:"allowNonNativeCorrections"`
	ImageURL                  string  `json:"imageURL"`
	EncryptedPassword         []byte  `json:"-"`
	Password                  string  `json:"-" sql:"-"`
	KarmaScore                uint    `json:"id"`

	UserLanguages []UserLanguage `json:"userLanguages"`
}

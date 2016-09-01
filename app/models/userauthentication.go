package models

type UserAuthentication struct {
	Model
	UserId        uint
	FacebookToken string
	LoginToken    string `sql:"unique_index"`
}

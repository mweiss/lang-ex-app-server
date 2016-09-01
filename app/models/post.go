package models

type Post struct {
	Model
	AuthorId uint   `sql:"index" json:"authorId"`
	Title    string `json:"title"`
	Body     string `sql:"size:4000" json:"body"`
	Language string `json:"language"`

	PostCorrections []PostCorrection `json:"postCorrections"`
	User            User             `json:"user"`
}

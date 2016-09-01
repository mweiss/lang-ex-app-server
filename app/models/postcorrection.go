package models

type PostCorrection struct {
	Model
	PostId   uint   `sql:"index" json:"postId"`
	AuthorId uint   `json:"authorId"`
	Comment  string `sql:"size:4000" json:"comment"`

	PostEdits []PostEdit `json:"postEdits"`
}

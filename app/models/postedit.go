package models

type PostEdit struct {
	Model
	PostCorrectionId uint   `sql:"index" json:"postCorrectionId"`
	Section          string `json:"section"`
	StartIndex       uint   `json:"startIndex"`
	Length           uint   `json:"length"`
	NewText          string `json:"newText"`
	Comment          string `sql:"size:4000" json:"comment"`
}

package models

type UserLanguage struct {
	Model
	UserId     uint   `json:"userId"`
	Language   string `json:"language"`
	Level      string `json:"level"`
	IsLearning bool   `json:"isLearning"`
}

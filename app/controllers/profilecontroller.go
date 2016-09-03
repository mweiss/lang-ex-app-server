package controllers

import (
	"encoding/json"
	"github.com/mweiss/lang-ex-app-server/app/models"
	"github.com/revel/revel"
	"log"
	"net/http"
)

type ProfileController struct {
	App
}

type UpdateLanguageJson struct {
	Language   string `json:"language"`
	Level      string `json:"level"`
	IsLearning bool   `json:"isLearning"`
}

type UpdateProfileJson struct {
	Languages                      *[]UpdateLanguageJson `json:"languages"`
	DisplayName                    *string               `json:"displayName"`
	AllowNonNativeSpeakersToUpdate *bool                 `json:"allowNonNativeSpeakersToUpdate"`
}

func (c ProfileController) GetUser() models.User {
	var user models.User
	c.Txn.Where("id = ? AND deleted_at is null", c.UserId).First(&user)
	return user
}

func (c ProfileController) GetProfile() revel.Result {

	var user models.User = c.GetUser()

	// From the user id, fetch the self introduction post
	if user.Id != 0 {
		c.Txn.Model(&user).Related(&user.UserLanguages)
		log.Print(user.Id)
		return c.RenderJson(user)
	} else {
		c.Response.Status = http.StatusUnauthorized
		return c.Render()
	}
}

func (c ProfileController) UpdateLanguages(languages []UpdateLanguageJson) {
	// Delete the user's old languages
	c.Txn.Exec("UPDATE user_languages SET deleted_at = now() WHERE user_id = ?", c.UserId)

	// insert new languges for the user
	for _, language := range languages {
		var userLanguage models.UserLanguage
		userLanguage.Language = language.Language
		userLanguage.Level = language.Level
		userLanguage.IsLearning = language.IsLearning
		userLanguage.UserId = c.UserId
		if c.Txn.NewRecord(userLanguage) {
			c.Txn.Create(&userLanguage)
		}
	}
}

func (c ProfileController) UpdateUser(updateProfileJson UpdateProfileJson) {
	var user models.User = c.GetUser()
	if user.Id != 0 {
		if updateProfileJson.DisplayName != nil {
			user.Name = *updateProfileJson.DisplayName
		}
		if updateProfileJson.AllowNonNativeSpeakersToUpdate != nil {
			user.AllowNonNativeCorrections = *updateProfileJson.AllowNonNativeSpeakersToUpdate
		}
		c.Txn.Save(&user)
	}
}

func LanguageIsSupported(language string) bool {
	switch language {
	case
		"en",
		"de",
		"zh",
		"ja":
		return true
	}
	return false
}

func ValidLanguageLevel(level string) bool {
	switch level {
	case
		"Native",
		"Fluent",
		"Advanced",
		"Intermediate",
		"Beginner":
		return true
	}
	return false
}

func (c ProfileController) ValidateLanguages(languages *[]UpdateLanguageJson) {
	if languages != nil {
		for _, language := range *languages {
			if !LanguageIsSupported(language.Language) {
				c.Validation.Error("Language is not supported: " + language.Language)
			}
			if !ValidLanguageLevel(language.Level) {
				c.Validation.Error("Language level is not supported: " + language.Level)
			}
		}
	}
}

func (c ProfileController) ValidateDisplayName(displayName *string) {
	if displayName != nil {
		c.Validation.MaxSize(*displayName, 255)
		c.Validation.MinSize(*displayName, 1)
	}
}

func (c ProfileController) UpdateProfile() revel.Result {
	var updateProfileJson UpdateProfileJson

	log.Print(c.UserId)

	// Check to make sure we have a valid user
	if c.UserId == 0 {
		c.Response.Status = http.StatusUnauthorized
	} else {
		// Read in the json data
		err := json.NewDecoder(c.Request.Body).Decode(&updateProfileJson)

		if err != nil {
			log.Print("JSON decode error: ", err)
		} else {

			log.Print(updateProfileJson.Languages)
			// Validate the input parameters
			c.ValidateLanguages(updateProfileJson.Languages)
			c.ValidateDisplayName(updateProfileJson.DisplayName)

			if !c.Validation.HasErrors() {
				// Update the user's profile and user languages
				if updateProfileJson.Languages != nil {
					c.UpdateLanguages(*updateProfileJson.Languages)
				}
				if updateProfileJson.DisplayName != nil || updateProfileJson.AllowNonNativeSpeakersToUpdate != nil {
					c.UpdateUser(updateProfileJson)
				}
			}
		}

		if c.Validation.HasErrors() {
			log.Print("")
			log.Print(c.Validation.Errors)
			c.Response.Status = http.StatusBadRequest
		} else {
			c.Response.Status = http.StatusOK
		}

		defer c.Request.Body.Close()
	}

	return c.Render()
}

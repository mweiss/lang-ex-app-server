package controllers

import (
	"github.com/mweiss/lang-ex-app-server/app/models"
	"github.com/revel/revel"
)

type ProfileController struct {
	App
}

type ProfileJson struct {
	User             *models.User
	SelfIntroduction *models.Post
}

func (c ProfileController) GetProfile() revel.Result {
	userid := c.GetUserId()

	var profileJson ProfileJson
	var user models.User

	c.Txn.Where("id = ? AND deleted_at is null", userid).First(&user)

	// From the user id, fetch the self introduction post
	if user.Id != 0 {
		profileJson.User = &user

		var selfIntroduction models.Post
		c.Txn.Where("id = ? AND deleted_at is null ", user.SelfIntroductionPostId).First(&selfIntroduction)
		if selfIntroduction.Id != 0 {
			profileJson.SelfIntroduction = &selfIntroduction
		}
	}

	return c.RenderJson(profileJson)
}

// TODO: implement this correctly.  For now, just return a constant.  This should maybe go
// a level above so that all controllers have access to this method.
func (c ProfileController) GetUserId() int64 {
	return 1
}

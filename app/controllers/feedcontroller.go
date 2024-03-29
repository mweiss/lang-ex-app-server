package controllers

import (
	"encoding/json"
	"github.com/mweiss/lang-ex-app-server/app/models"
	"github.com/revel/revel"
	"log"
)

type FeedController struct {
	App
}

type PostResponseJson struct {
	Posts []models.Post `json:posts`
}

type PostCorrectionsJson struct {
	PostCorrections []models.PostCorrection `json:postCorrections`
}

type PostIdJson struct {
	Id uint `json:id`
}

type PostCorrectIdJson struct {
	Id      uint   `json:id`
	EditIds []uint `json:editIds`
}

func (c FeedController) CreateCorrection(postId uint) revel.Result {
	var postCorrection models.PostCorrection
	err := json.NewDecoder(c.Request.Body).Decode(&postCorrection)

	if err != nil {
		log.Print("JSON decode error: ", err)
	} else {
		// Verify that the post exists
		var post models.Post = c.FetchPostById(postId)

		if post.Id != 0 {
			postCorrection.PostId = postId
			postCorrection.AuthorId = c.UserId
			if c.Txn.NewRecord(postCorrection) {
				c.Txn.Create(&postCorrection)
			}
		}
	}

	defer c.Request.Body.Close()

	log.Print(postCorrection)

	if postCorrection.Id != 0 {
		editIds := make([]uint, len(postCorrection.PostEdits))
		for i, v := range postCorrection.PostEdits {
			editIds[i] = v.Id
		}
		log.Print(editIds)
		return c.RenderJson(PostCorrectIdJson{Id: postCorrection.Id, EditIds: editIds})
	} else {
		c.Response.Status = 400
		return c.RenderText("")
	}
}

func (c FeedController) CreatePost() revel.Result {
	var post models.Post
	err := json.NewDecoder(c.Request.Body).Decode(&post)

	if err != nil {
		log.Print("JSON decode error: ", err)
	} else {
		if c.Txn.NewRecord(post) {
			post.AuthorId = c.UserId
			c.Txn.Create(&post)
		}
	}

	defer c.Request.Body.Close()
	if post.Id != 0 {
		return c.RenderJson(PostIdJson{post.Id})
	} else {
		c.Response.Status = 400
		return c.RenderText("")
	}
}

func (c FeedController) GetCorrections(id uint) revel.Result {
	var post models.Post = c.FetchPostById(id)

	// TODO: I should consolidate how I do this into one helper method.
	if post.Id == 0 {
		c.Response.Status = 404
		return c.RenderText("")
	} else {
		return c.RenderJson(PostCorrectionsJson{post.PostCorrections})
	}
}

func (c FeedController) FetchPostById(id uint) models.Post {
	var post models.Post
	c.Txn.Where("id = ? AND deleted_at is null", id).First(&post)

	if post.Id != 0 {
		var posts []models.Post = []models.Post{post}
		c.FillPosts(posts)
	}
	return post
}

func (c FeedController) GetPostById(id uint) revel.Result {

	var post models.Post = c.FetchPostById(id)

	if post.Id == 0 {
		c.Response.Status = 404
		return c.RenderText("") // TODO, there might be a cleaner way of rendering an empty response
	} else {
		return c.RenderJson(post)
	}
}

// TODO: Need to implement the 'skipped' parameter
func (c FeedController) GetPostByUser() revel.Result {

	// Fetch the request parameters.
	var user uint
	var correctedByUser uint

	c.Params.Bind(&user, "user")
	c.Params.Bind(&correctedByUser, "correctedByUser")

	var posts []models.Post = []models.Post{}

	log.Print("hey")
	log.Print(user)
	if user != 0 {
		c.Txn.Where("author_id = ? AND deleted_at is null", user).Find(&posts)
		log.Print("fetched some posts")
		log.Print(posts)
	} else if correctedByUser != 0 {
		c.Txn.Where("deleted_at is null AND 1 <= (SELECT count(*) FROM post_corrections pc "+
			"WHERE pc.author_id = ? and pc.id = posts.id)", correctedByUser).Find(&posts)
	}
	c.FillPosts(posts)

	return c.RenderJson(PostResponseJson{posts})
}

func (c FeedController) Feed() revel.Result {
	// For now, let's just return every post in the languages they want
	userLanguage := c.LanguagesToLearn()

	// Fetch all posts
	var posts []models.Post
	c.Txn.Where("language in (?) and deleted_at is null", userLanguage).Find(&posts)
	c.FillPosts(posts)
	return c.RenderJson(PostResponseJson{posts})
}

// TODO:
// This is extremely inefficient, but I'm going to do a fetch for each
// sub entity because it doesn't look like GORM supports one to many relationships,
// in a way that would query efficiently (besides preloading, which seems equally shitty)
// It's just prototype code, so fuck it.
func (c FeedController) FillPosts(posts []models.Post) {
	for i := range posts {
		c.Txn.Model(posts[i]).Related(&posts[i].PostCorrections)
		c.Txn.Where("id = ?", posts[i].AuthorId).First(&posts[i].User)
		for j := range posts[i].PostCorrections {
			c.Txn.Model(posts[i].PostCorrections[j]).Related(&posts[i].PostCorrections[j].PostEdits)
		}
	}
}

// TODO: Fix this trash code to determine what languages a user is learning
func (c FeedController) LanguagesToLearn() []string {

	// Default the user id
	var languages []string

	rows, err := c.Txn.Raw("SELECT language FROM user_languages WHERE user_id = ? and is_learning = 1 and deleted_at is null", c.UserId).Rows()
	defer rows.Close()

	// TODO: probably not the right way to handle a sql error, maybe it's better do a panic statement
	if err != nil {
		for rows.Next() {
			var s *string
			rows.Scan(&s)
			languages = append(languages, *s)
		}
	}

	if len(languages) == 0 {
		languages = []string{"en"}
	}
	return languages
}

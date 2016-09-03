package controllers

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mweiss/lang-ex-app-server/app/models"
	"github.com/revel/revel"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

type LoginController struct {
	App
}

type LoginResponseJson struct {
	Token string `json:"token"`
}

type FacebookJson struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Picture   PictureDataJson `json:"picture"`
	Email     string          `json"email"`
}

type PictureDataJson struct {
	Data PictureJson `json:"data"`
}

type PictureJson struct {
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	IsSilhouette bool   `json:"is_silhouette"`
	Url          string `json:"url"`
}

func (c LoginController) FetchFacebookToken() string {
	var fbToken string
	c.Params.Bind(&fbToken, "fbToken")
	return fbToken
}

func (c LoginController) FetchFacebookData() (FacebookJson, error) {
	fbToken := c.FetchFacebookToken()
	fbJson := FacebookJson{}

	// Create the URL to fetch the user's data from facebook
	var Url *url.URL
	Url, err := url.Parse("https://graph.facebook.com/me")

	if err != nil {
		return fbJson, err
	}

	parameters := url.Values{}
	parameters.Add("fields", "id,name,first_name,last_name,picture.type(square).width(200).height(200),email")
	parameters.Add("access_token", fbToken)
	Url.RawQuery = parameters.Encode()

	// Make the request to facebook
	res, err := http.Get(Url.String())

	if err != nil {
		return fbJson, err
	}

	// Decode the json we receive
	err = json.NewDecoder(res.Body).Decode(&fbJson)

	if err != nil {
		return fbJson, err
	}

	return fbJson, nil
}

func (c LoginController) FetchOrCreateUser(fbData FacebookJson) models.User {
	// Check to see if we have a user already
	var user models.User

	c.Txn.Where("facebook_id = ? AND deleted_at is null", fbData.Id).First(&user)
	if user.Id == 0 {
		user.FacebookId = fbData.Id
		user.Name = fbData.Name
		user.FirstName = &fbData.FirstName
		user.LastName = &fbData.LastName
		user.Email = fbData.Email
		user.ImageURL = fbData.Picture.Data.Url
		if c.Txn.NewRecord(user) {
			c.Txn.Create(&user)
		}
	}

	return user
}

func (c LoginController) CreateAuthenticationRow(user models.User, fbToken string) string {
	var userAuthentication models.UserAuthentication
	userAuthentication.UserId = user.Id
	userAuthentication.FacebookToken = c.FetchFacebookToken()
	userAuthentication.LoginToken = createLoginToken()
	if c.Txn.NewRecord(userAuthentication) {
		c.Txn.Create(&userAuthentication)
	}
	return userAuthentication.LoginToken
}

// Creates a randomized login token seeded by the hostname, process id, current time, and two random ints that
// are then hashed to sha512
func createLoginToken() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Print(err)
	}

	s := fmt.Sprintf("%v_%v_%v_%v_%v_apjd7xc", hostname, os.Getpid(), time.Now(), rand.Int63(), rand.Int63())
	hash := sha512.Sum512([]byte(s))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (c LoginController) Login() revel.Result {
	fbData, err := c.FetchFacebookData()
	if err != nil {
		log.Print(err)
		// Return no access error depending on the error we get back
	}

	user := c.FetchOrCreateUser(fbData)
	token := c.CreateAuthenticationRow(user, c.FetchFacebookToken())

	return c.RenderJson(LoginResponseJson{token})
}

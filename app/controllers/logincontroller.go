package controllers

import (
	"github.com/revel/revel"
)

type LoginController struct {
	App
}

type LoginResponseJson struct {
	Token string `json:"token"`
}

func (c LoginController) Login() revel.Result {
	return c.RenderJson(LoginResponseJson{"1234"})
}

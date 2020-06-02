package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	fmt.Println(c.Ctx.Request.RemoteAddr)
	c.TplName ="index.html"
	_ = c.Render()
}
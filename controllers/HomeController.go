package controllers

import (
	"mbook/models"

	"github.com/astaxie/beego"
)

type HomeController struct {
	BaseController
}
/*
	展示首页的函数，主要是获取全部的分类信息，并返回到前端页面进行渲染展示
*/
func (c *HomeController) Index() {
	if cates, err := new(models.Category).GetCates(-1, 1); err == nil {
		c.Data["Cates"] = cates
	} else {
		beego.Error(err.Error())
	}

	c.TplName = "home/list.html"
}

func (c *HomeController) Index2() {
	c.TplName = "home/list.html"
}

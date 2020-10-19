package controllers

import (
	"math"
	"strconv"
	"mbook/models"
	"mbook/utils"

	"github.com/astaxie/beego"
)

type ExploreController struct {
	BaseController
}

func (c *ExploreController) Index() {
	var (
		cid       int //分类id
		cate      models.Category
		urlPrefix = beego.URLFor("ExploreController.Index")
	)

	if cid, _ = c.GetInt("cid"); cid > 0 { // 从前端接收
		cateModel := new(models.Category)
		cate = cateModel.Find(cid)
		c.Data["Cate"] = cate
	}

	c.Data["Cid"] = cid
	c.TplName = "explore/index.html"

	pageIndex, _ := c.GetInt("page", 1)
	pageSize := 24

	books, totalCount, err := models.NewBook().HomeData(pageIndex, pageSize, cid)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}

	if totalCount > 0 { // 如果找到的结果大于0，则得进行分页处理，显示页码
		urlSuffix := ""
		if cid > 0 {
			urlSuffix = urlSuffix + "&cid=" + strconv.Itoa(cid)
		}
		html := utils.NewPaginations(4, totalCount, pageSize, pageIndex, urlPrefix, urlSuffix)
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}

	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	c.Data["Lists"] = books
}

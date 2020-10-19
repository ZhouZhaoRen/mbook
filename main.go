package main

import (
	_ "mbook/routers"
	"github.com/astaxie/beego"
	_ "mbook/sysinit" // 引入这个包，才能一开始就初始化数据库
)

func main() {
	beego.Run()
}


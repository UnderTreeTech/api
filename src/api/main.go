package main

import (
	//如果不使用单例模式，可用此方法初始化mysql,redis的连接对象，但要都放在init函数中
	_ "api/dao"
	_ "api/docs"
	_ "api/routers"
	_ "api/util/xss"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	_ "xinda.im/boss/common"
	ctl "xinda.im/boss/controllers"
	_ "xinda.im/boss/routers"
)

func initFilter() {
	//	beego.InsertFilter("/*", beego.BeforeExec, FilterUser)
}

func initStatic() {
	//默认路径是从main.go所在的文件夹开始，这里是boss
	beego.SetStaticPath("/static", "static")
	beego.SetStaticPath("/views", "views")
	beego.SetStaticPath("/fonts", "views/fonts")
	beego.SetStaticPath("/images", "views/images")
	beego.SetStaticPath("/css", "views/css")
	beego.SetStaticPath("/js", "views/js")
}

func init() {
	ctl.Init()
	initFilter()
	initStatic()

	beego.Info("执行完初始化")
}

func main() {
	beego.Info("运行main")
	beego.SetLogger("file", `{"filename":"log/boss.log"}`)
	beego.ErrorController(&ctl.ErrorController{})

	beego.Info("运行Run")
	beego.Run()
}

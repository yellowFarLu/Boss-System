package controllers

import (
	"github.com/astaxie/beego"
	"xinda.im/boss/models"
)

func Init() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	models.Init()
}

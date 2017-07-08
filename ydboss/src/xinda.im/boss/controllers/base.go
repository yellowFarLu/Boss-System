package controllers

import (
	"xinda.im/boss/models"
)

var (
	loginModal    *models.LoginModel
	accountModal  *models.AccountModel
	customerModel *models.CustomerModel
	agentModel    *models.AgentModel
	bizModal      *models.BizModel
)

////////////////////////////////权限///////////////////////////////////////
const (
	NotAuthority = "404.html"
)

const (
	CODE_SUCC = 200
	CODE_AUTH = 302
)

const (
	SESSION_USER = "userInfo"
)

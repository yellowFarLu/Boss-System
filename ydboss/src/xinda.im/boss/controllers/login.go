package controllers

import (
	"github.com/astaxie/beego"
	. "xinda.im/boss/common"
	_ "xinda.im/boss/models"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

// 单点登录
func (this *LoginController) YDLogin() {
	userName, rsp := loginModal.YDLogin(this.GetString("token"))
	var role int
	if rsp.IsOK() {
		account := &Account{AccountId: userName}
		this.LoginSetSession(account, &role)

		if role == ROLE_XIN_ADMIN {
			this.Redirect("agent.html", 302)
			return
		} else if role == ROLE_AGENT {
			this.Redirect("agent_employee.html", 302)
			return
		} else if role == ROLE_EMPLOYEE {
			this.Redirect("sale_customer.html", 302)
			return
		}

	} else {
		this.Redirect("login.html", CODE_AUTH)
	}

}

// 请求登录
func (this *LoginController) Login() {

	account := &Account{}
	account.AccountId = this.GetString("name")
	account.Pwd = this.GetString("password")
	rsp := NewResponse(SYSTEMSUCCESS)
	status := loginModal.Login(account)
	var role int
	if status == LOGINSTATUSOK {
		this.LoginSetSession(account, &role)
	}

	// fLogin := loginModal.FirstLogin(account.AccountId)
	// if fLogin {
	// 	role = -1 //跳转重密码界面
	// }

	rsp["status"] = status
	rsp["href"] = role
	this.Data["json"] = &rsp
	this.ServeJSON()
}

func (this *LoginController) LoginSetSession(account *Account, role *int) {
	newAcc := accountModal.GetUserInfo(account.AccountId)
	newAcc.AccountId = account.AccountId
	newAcc.Pwd = account.Pwd
	(*role) = newAcc.RoleId

	this.SetSession(SESSION_USER, newAcc)
}

func (this *LoginController) LoginOut() {
	this.DestroySession()
	this.Redirect("login.html", 302)
}

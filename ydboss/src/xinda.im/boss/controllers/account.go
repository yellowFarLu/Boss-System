package controllers

import (
	. "xinda.im/boss/common"

	"strconv"

	"github.com/astaxie/beego"
)

type AccountController struct {
	beego.Controller
}

func (c *AccountController) Get() {
	c.TplName = "sale_customer.html"
}

func (this *AccountController) RenderEmployeeInfo() {
	this.TplName = "employee_info.html"
}

func (this *AccountController) FilterAccount() {
	g := this.Input()
	sortKey := this.GetString("sortKey", "timex")
	sortIndex, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	tags := g["tags[]"]
	provinces := g["provinces[]"]

	account := this.GetSession(SESSION_USER).(*Account)

	filter := &FilterParam{
		G:         g,
		SortKey:   sortKey,
		SortIndex: sortIndex,
		Page:      page,
		Tags:      tags,
		Provinces: provinces,
	}

	list := accountModal.FilterAccountList(*filter, account)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	record, page_total := accountModal.FilterAccountCount(*filter, account)
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AccountController) KeyWordAccount() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)

	page, _ := this.GetInt("page", 1)

	if account.RoleId == ROLE_EMPLOYEE {
		rsp = NewResponse(SYSTEMERROR)
	} else {
		keyW := this.GetString("keyWord")
		accounts := accountModal.ShowKeyWorldAccount(keyW, account, page)
		rsp = NewResponse(SYSTEMSUCCESS)
		rsp["list"] = accounts

		var page_total, record int
		record, page_total = accountModal.GetCountInfo(strconv.Itoa(account.AgentId), keyW)
		rsp["page_info"] = NewPageInfo(page, page_total, record)
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (c *AccountController) GetAccountList() {
	GetAccountInfo(c)
}

func (c *AccountController) LoginOut() {

}

func (c *AccountController) GetLoginList() {
	c.TplName = "login.html"
}

func GetAccountInfo(this *AccountController) {
	page, _ := this.GetInt("page", 1)

	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE {
		rsp = NewResponse(SYSTEMERROR)
	} else {
		accounts := accountModal.ShowKeyWorldAccount("", account, page)
		this.Data["list"] = accounts

		rsp = NewResponse(SYSTEMSUCCESS)
		rsp["list"] = accounts

		var page_total, record int
		record, page_total = accountModal.GetCountInfo(strconv.Itoa(account.AgentId), "")
		rsp["page_info"] = NewPageInfo(page, page_total, record)
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AccountController) UpdateAccountPsw() {
	id := this.GetSession(SESSION_USER).(*Account).AccountId
	new_psw := this.GetString("new_psw")
	code := SYSTEMSUCCESS
	if err := accountModal.UpdateAccountPsw(id, new_psw); err != nil {
		code = SYSTEMERROR
	}
	rsp := NewResponse(code)
	this.Data["json"] = rsp
	this.ServeJSON()
}
func (this *AccountController) ResetAccountPsw() {
	id := this.GetString("id")
	new_psw := this.GetString("new_psw")
	user := this.GetSession(SESSION_USER).(*Account)
	code := SYSTEMSUCCESS

	if ok := accountModal.CheckRightByAccount(user, id); !ok {
		code = SYSTEMERROR
	} else {
		if err := accountModal.UpdateAccountPsw(id, new_psw); err != nil {
			code = SYSTEMERROR
		}
	}
	rsp := NewResponse(code)
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AccountController) AlertAccount() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE {
		rsp = NewResponse(SYSTEMERROR)
		this.Data["json"] = rsp
		this.ServeJSON()
	} else {
		emp := &Account{}
		emp.Name = this.GetString("Name")
		emp.Mobile = this.GetString("Mobile")
		emp.Gender, _ = this.GetInt("Gender")
		emp.AccountId = this.GetString("AccountId")
		ok := accountModal.AlertAccountInfo(emp)

		this.JSonCode(ok)
	}
}

func (this *AccountController) DelAccount() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE {
		rsp = NewResponse(SYSTEMERROR)
		this.Data["json"] = rsp
		this.ServeJSON()
	} else {
		m := this.Input()
		DelArr := m["AccountIds[]"]
		ok := accountModal.DelAccount(DelArr)

		this.JSonCode(ok)
	}

}

func (this *AccountController) AddAccount() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	var AccountId = this.GetString("AccountId")
	if account.RoleId == ROLE_EMPLOYEE {
		rsp = NewResponse(SYSTEMERROR)
	} else if accountModal.AcountIsExisted(AccountId) {
		rsp = NewResponse(ACCOUNTEXIST)
		this.Data["json"] = rsp
		this.ServeJSON()

	} else {
		emp := &Account{}
		emp.AccountId = AccountId
		emp.Pwd = this.GetString("Pwd")
		emp.Name = this.GetString("Name")
		emp.Gender, _ = this.GetInt("Gender")
		emp.Mobile = this.GetString("Mobile")

		emp.AgentId = account.AgentId

		// 插入数据库
		ok := accountModal.AddAccount(emp)
		this.JSonCode(ok)
	}
}

func (this *AccountController) GetCustomComments() {
	account := this.GetSession(SESSION_USER).(*Account)
	page, _ := this.GetInt("page", 1)
	id := this.GetString("id")
	var list []*Comment
	var err error
	list, err = customerModel.GetAccountComments(id, page, "custom", account)

	code := SYSTEMSUCCESS
	if err != nil {
		code = SYSTEMERROR
	}
	rsp := NewResponse(code)
	rsp["list"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AccountController) GetOpportComments() {
	account := this.GetSession(SESSION_USER).(*Account)
	id := this.GetString("id")
	page, _ := this.GetInt("page", 1)
	list, err := customerModel.GetAccountComments(id, page, "opport", account)
	code := SYSTEMSUCCESS
	if err != nil {
		code = SYSTEMERROR
	}
	rsp := NewResponse(code)
	rsp["list"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AccountController) GetAccountBySort() {
	key := this.GetString("sortKey")
	index, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	user := this.GetSession(SESSION_USER).(*Account)
	list, page_info := accountModal.GetAccountBySort(key, index, page, user)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AccountController) Prepare() {
	if this.GetSession(SESSION_USER) == nil {
		this.Redirect("login.html", 302)
		return
	}
	user := this.GetSession(SESSION_USER).(*Account)
	this.Data["user"] = user
}

func (c *AccountController) JSonCode(ok bool) {
	var code int
	if ok {
		code = DELOK
	} else {
		code = DELFAIL
	}
	var rsp = NewResponse(code)
	c.Data["json"] = rsp
	c.ServeJSON()
}

package controllers

import (
	"github.com/astaxie/beego"
	. "xinda.im/boss/common"
)

type AgentController struct {
	beego.Controller
}

func (this *AgentController) Get() {
	this.GetAuthority("agent.html")
}

func (this *AgentController) AgentEmployee() {
	this.GetAuthority("agent_employee.html")
}

func (this *AgentController) AgentCustomer() {
	this.GetAuthority("agent_customer.html")
}

func (this *AgentController) KeyWordAgent() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE || account.RoleId == ROLE_AGENT {
		rsp = NewResponse(SYSTEMERROR)
	} else { // 渠道或超级管理员以上
		rsp = NewResponse(SYSTEMSUCCESS)
		key := this.GetString("keyWord")
		page, _ := this.GetInt("page", 1)
		agents := agentModel.ShowKeyWorldAgent(key, page)

		rsp["list"] = agents
		var page_total, record int
		record, page_total = agentModel.GetCountInfo(key)
		rsp["page_info"] = NewPageInfo(page, page_total, record)
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AgentController) GetAgentList() {
	user := this.GetSession(SESSION_USER).(*Account)
	this.GetChannelInfo(user)
}

func (c *AgentController) GetChannelInfo(account *Account) {
	var rsp map[string]interface{}
	if account.RoleId == ROLE_EMPLOYEE || account.RoleId == ROLE_AGENT {
		rsp = NewResponse(SYSTEMERROR)
		return
	}
	page, _ := c.GetInt("page", 1)
	agents := agentModel.GetAgentList(account.AccountId, page)
	rsp = NewResponse(SYSTEMSUCCESS)
	rsp["list"] = agents

	var page_total, record int
	record, page_total = agentModel.GetCountInfo("")
	rsp["page_info"] = NewPageInfo(page, page_total, record)

	c.Data["json"] = rsp
	c.ServeJSON()
}

func (this *AgentController) GetAgentBySort() {
	key := this.GetString("sortKey")
	index, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	list, page_info := agentModel.GetAgentBySort(key, index, page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

// 筛选渠道
func (this *AgentController) FilterAgent() {

	g := this.Input()

	sortKey := this.GetString("sortKey", "timex")
	sortIndex, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	tags := g["tags[]"]
	provinces := g["provinces[]"]

	filter := &FilterParam{
		G:         g,
		SortKey:   sortKey,
		SortIndex: sortIndex,
		Page:      page,
		Tags:      tags,
		Provinces: provinces,
	}

	rsp := NewResponse(SYSTEMSUCCESS)
	list := agentModel.FilterAgentList(*filter)
	rsp["list"] = list
	record, page_total := agentModel.FilterAgentCount(*filter)
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *AgentController) AlertAgent() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE || account.RoleId == ROLE_AGENT {
		rsp = NewResponse(SYSTEMERROR)
		this.Data["json"] = rsp
		this.ServeJSON()
	} else {
		agent := &Agent{}
		agent.Name = this.GetString("AgentName")
		agent.Contacts = this.GetString("Contacts")
		agent.Mail = this.GetString("ContactsMail")
		agent.Mobile = this.GetString("ContactsPhone")
		agent.Note = this.GetString("Note")
		agent.AgentId, _ = this.GetInt("AgentId")
		ok := agentModel.AlertAgentInfo(agent)

		this.JSonCode(ok)
	}
}

func (this *AgentController) DelAgent() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE || account.RoleId == ROLE_AGENT {
		rsp = NewResponse(SYSTEMERROR)
		this.Data["json"] = rsp
		this.ServeJSON()
	} else {
		m := this.Input()
		DelArr := m["AgentIds[]"]
		ok := agentModel.DelAgent(DelArr)

		this.JSonCode(ok)
	}

}

func (this *AgentController) InsertAgent() {
	var rsp map[string]interface{}
	account := this.GetSession(SESSION_USER).(*Account)
	var accountId = this.GetString("AccountId")
	if account.RoleId == ROLE_EMPLOYEE || account.RoleId == ROLE_AGENT {
		rsp = NewResponse(SYSTEMERROR)

	} else if accountModal.AcountIsExisted(accountId) {

		rsp = NewResponse(ACCOUNTEXIST)
		this.Data["json"] = rsp
		this.ServeJSON()
	} else {
		agent := &Agent{}
		newAcc := &Account{}
		agent.Name = this.GetString("AgentName")
		agent.Contacts = this.GetString("Contacts")
		agent.Mail = this.GetString("ContactsMail")
		agent.Mobile = this.GetString("ContactsPhone")
		agent.Note = this.GetString("Note")

		newAcc.AccountId = accountId
		newAcc.Name = this.GetString("Name")
		newAcc.Pwd = this.GetString("Password")
		newAcc.Gender, _ = this.GetInt("Sex")
		newAcc.Mobile = this.GetString("ContactsPhone")
		ok := agentModel.AddAgent(agent, newAcc)

		this.JSonCode(ok)
	}
}

// get方法权限判断
func (this *AgentController) GetAuthority(aim string) {
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE { // 渠道的信息只能由超级管理员修改
		this.Redirect(NotAuthority, 302)
	} else {
		this.TplName = aim
	}
}

func (c *AgentController) JSonCode(ok bool) {
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

func (this *AgentController) Prepare() {
	if this.GetSession(SESSION_USER) == nil {
		this.Redirect("login.html", 302)
		return
	}

	user := this.GetSession(SESSION_USER).(*Account)
	this.Data["user"] = user
}

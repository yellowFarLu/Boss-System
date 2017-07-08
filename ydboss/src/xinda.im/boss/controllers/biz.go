package controllers

import (
	"github.com/astaxie/beego"
	. "xinda.im/boss/common"
)

type BizController struct {
	beego.Controller
}

func (this *BizController) Get() {
	this.TplName = "opportunities.html"
}

func (this *BizController) RenderAgentBiz() {
	this.TplName = "agent_biz.html"
}

func (this *BizController) RenderOpportSearch() {
	this.TplName = "opport_search.html"
}

// 获取全部商机
func (this *BizController) GetBizunitiesList() {
	account := this.GetSession(SESSION_USER).(*Account)
	this.GetOpcInfo(account)
}

// 渠道展示商机
func (this *BizController) AgentShowBizs() {
	account := this.GetSession(SESSION_USER).(*Account)
	page, _ := this.GetInt("page", 1)
	biz := bizModal.AgentShowBizs(account, "", page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = biz

	var page_total, record int
	record, page_total = bizModal.GetAgentBizCountInfo(account.AgentId, "")
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	this.Data["json"] = rsp
	this.ServeJSON()
}

// 渠道关键字查询
func (this *BizController) GetAgentOpportKeyword() {
	account := this.GetSession(SESSION_USER).(*Account)
	page, _ := this.GetInt("page", 1)
	key := this.GetString("keyWord")
	biz := bizModal.AgentShowBizs(account, key, page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = biz

	var page_total, record int
	record, page_total = bizModal.GetAgentBizCountInfo(account.AgentId, key)
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	this.Data["json"] = rsp
	this.ServeJSON()
}

// 销售根据关键字查询
func (c *BizController) KeyWordBizList() {
	account := c.GetSession(SESSION_USER).(*Account)
	key := c.GetString("keyWord")
	page, _ := c.GetInt("page", 1)
	biz := bizModal.ShowBizs(account, key, page)
	rsp := NewResponse(SYSTEMSUCCESS)

	rsp["list"] = biz
	var page_total, record int
	record, page_total = bizModal.GetCountInfo(account.AccountId, key)

	rsp["page_info"] = NewPageInfo(page, page_total, record)
	c.Data["json"] = rsp
	c.ServeJSON()
}

func (this *BizController) GetBizComments() {
	id, _ := this.GetInt("id")
	_type, _ := this.GetInt("type")
	page, _ := this.GetInt("page", 1)
	list, err := bizModal.GetBizComments(id, _type, page)
	code := SYSTEMSUCCESS
	if err != nil {
		code = SYSTEMERROR
	}
	rsp := NewResponse(code)
	rsp["Comments"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

// 渠道获取商机详细页
func (this *BizController) GetAgentBizDetail() {
	this.TplName = "agent_opport_info.html"
}

func (this *BizController) GetEmpBizDetail() {
	this.TplName = "opport_info.html"
}

// 获取详细页
func (this *BizController) EmpBizInfo() {
	rsp := NewResponse(SYSTEMSUCCESS)
	account := this.GetSession(SESSION_USER).(*Account)
	id := this.GetString("id")

	var biz *Biz
	biz = bizModal.EmpBizDetail(id, account)
	rsp["biz"] = *biz

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *BizController) FilterBizEmpList() {
	user := this.GetSession(SESSION_USER).(*Account)
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

	list := bizModal.FilterBizList(filter, user)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	record, page_total := bizModal.FilterBizCount(*filter)
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *BizController) FilterBizAgentList() {
	user := this.GetSession(SESSION_USER).(*Account)
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

	list := bizModal.FilterAgentBizList(filter, user)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	record, page_total := bizModal.FilterBizCount(*filter)
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (c *BizController) GetOpcInfo(account *Account) {
	page, _ := c.GetInt("page", 1)
	biz := bizModal.ShowBizs(account, "", page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = biz

	var page_total, record int
	record, page_total = bizModal.GetCountInfo(account.AccountId, "")
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	c.Data["json"] = rsp
	c.ServeJSON()
}

// AlertCommentsFromC
func (this *BizController) AlertCommentsFromC() {
	bizId, _ := this.GetInt("biz_id", 0)
	comments := this.GetString("comments")
	ok := bizModal.AlertCommentsFromC(bizId, comments)
	this.JSonCode(ok)
}

// 修改商机
func (this *BizController) AlertBiz() {
	g := this.Input()
	account := this.GetSession(SESSION_USER).(*Account)
	biz := &Biz{}
	biz.Customer = &Customer{}
	biz.BizId, _ = this.GetInt("opportunities_id")
	biz.Title = this.GetString("title")
	biz.Content = this.GetString("content")
	biz.Customer.CustomerId, _ = this.GetInt("customerId")

	biz.Amount, _ = this.GetFloat("quota")

	biz.Status, _ = this.GetInt("status")
	rt := this.GetString("real_time")
	biz.RealTime = &rt
	tags := this.Input()["tags[]"]
	tagsAlert, _ := this.GetInt("tagsAlert")
	ok := bizModal.AlertBiz(g, account, biz, tags, tagsAlert)
	rsp := NewResponse(SYSTEMSUCCESS)
	if !ok {
		rsp["code"] = SYSTEMERROR
	}
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *BizController) InsertCommentsFromC() {
	user := this.GetSession(SESSION_USER).(*Account)
	id, _ := this.GetInt("id")
	content := this.GetString("comment_content")
	code := SYSTEMSUCCESS
	if ok := bizModal.InsertCommentsFromC(content, id, user); !ok {
		code = SYSTEMERROR
	}
	rsp := NewResponse(code)
	this.Data["json"] = rsp
	this.ServeJSON()
}

// 删商机
func (this *BizController) DelBiz() {
	account := this.GetSession(SESSION_USER).(*Account)
	DelArr := this.Input()["OpportIds[]"]
	status, _ := this.GetInt("status")
	ok := bizModal.DeleteBizs(account, DelArr, status)
	rsp := NewResponse(SYSTEMSUCCESS)
	if !ok {
		rsp["code"] = SYSTEMERROR
	}
	this.Data["json"] = rsp
	this.ServeJSON()
}

//插入商机
func (this *BizController) InsertBiz() {
	account := this.GetSession(SESSION_USER).(*Account)
	g := this.Input()
	biz := &Biz{}
	biz.Customer = &Customer{}

	biz.AccountId = account.AccountId
	biz.Customer.CustomerId, _ = this.GetInt("customerId")
	biz.Content = this.GetString("content")
	biz.Title = this.GetString("title")
	biz.Amount, _ = this.GetFloat("quota")
	biz.Status, _ = this.GetInt("status")
	est := this.GetString("estimatetime")
	biz.EstimateTime = &est
	realTime := this.GetString("real_time")
	biz.RealTime = &realTime
	tags := this.Input()["tags[]"]

	ok := bizModal.InsertBiz(g, account, biz, tags)
	rsp := NewResponse(SYSTEMSUCCESS)
	if !ok {
		rsp["code"] = SYSTEMERROR
	}
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *BizController) GetAgentOpportBySort() {
	key := this.GetString("sortKey")
	index, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	user := this.GetSession(SESSION_USER).(*Account)
	list, page_info := bizModal.GetAgentOpportBySort(key, index, page, user)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *BizController) GetOpportBySort() {
	key := this.GetString("sortKey")
	index, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	user := this.GetSession(SESSION_USER).(*Account)
	list, page_info := bizModal.GetOpportBySort(key, index, page, user)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *BizController) Prepare() {
	if this.GetSession(SESSION_USER) == nil {
		this.Redirect("login.html", 302)
		return
	}
	user := this.GetSession(SESSION_USER).(*Account)
	this.Data["user"] = user
}

func (this *BizController) JSonCode(ok bool) {
	var code int
	if ok {
		code = DELOK
	} else {
		code = DELFAIL
	}
	var rsp = NewResponse(code)
	this.Data["json"] = rsp
	this.ServeJSON()
}

package controllers

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
	. "xinda.im/boss/common"
)

type CustomerController struct {
	beego.Controller
}

func (this *CustomerController) Get() {
	this.GetAuthority("all_customer.html")
}

func (this *CustomerController) RenderCustomerSearch() {
	this.TplName = "customer_search.html"
}

// 查看用户是否存在
func (this *CustomerController) CustomerIsExists() {
	name := this.GetString("name")
	ok := customerModel.CustomerIsExists(name)
	this.JSonCode(ok)
}

// 添加备注
func (this *CustomerController) AddComment() {
	account := this.GetSession(SESSION_USER).(*Account)
	comment := &Comment{}
	comment.CustomerId, _ = this.GetInt("customer_id")
	comment.Comments = this.GetString("comment_area")
	ok := customerModel.AddComment(account, comment)
	this.JSonCode(ok)
}

// 获取销售客户详细页
func (this *CustomerController) GetSaleCustomerDetail() {
	this.TplName = "sale_customer_info.html"
}

// 获取渠道客户详细页
func (this *CustomerController) GetAgentCustomerDetail() {
	this.TplName = "agent_customer_info.html"
}

// 获取超级管理员详细页
func (this *CustomerController) GetSuperCustomerDetail() {
	this.TplName = "all_customer_info.html"
}

func (this *CustomerController) GetEmpCustomInfo() {
	rsp := NewResponse(SYSTEMSUCCESS)
	account := this.GetSession(SESSION_USER).(*Account)
	id := this.GetString("customerId")

	// 先判断用户id对应的用户是不是这个账号的
	if !customerModel.CustomerBelongAccount(id, account) {
		rsp["code"] = SYSTEMERROR
	} else {
		var customer *Customer
		customer = customerModel.EmpCustomerDetail(id, account)
		rsp["customer"] = *customer
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetAgentCustomerInfo() {
	rsp := NewResponse(SYSTEMSUCCESS)
	account := this.GetSession(SESSION_USER).(*Account)
	id := this.GetString("customerId")

	// 先判断用户id对应的用户是不是这个渠道的
	if !customerModel.CustomerBelongAgent(id, account) {
		rsp["code"] = SYSTEMERROR
	} else {
		var customer *Customer
		customer = customerModel.AgentCustomerDetail(id, account)
		rsp["customer"] = *customer
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetSuperCustomerInfo() {
	rsp := NewResponse(SYSTEMSUCCESS)
	id := this.GetString("customerId")
	var customer *Customer
	customer = customerModel.SuperCustomerDetail(id)

	rsp["customer"] = *customer
	this.Data["json"] = rsp
	this.ServeJSON()
}

// 获取销售客户的备注
func (this *CustomerController) GetSaleComments() {
	id, _ := this.GetInt("customerId")
	_type, _ := this.GetInt("type", 0)
	page, _ := this.GetInt("page", 1)
	list := customerModel.GetComments(id, _type, page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["Comments"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetSaleCustomerTypeList() {
	_type := this.GetString("type")
	page, _ := this.GetInt("page", 1)
	user := this.GetSession(SESSION_USER).(*Account)
	list, page_info := customerModel.GetSaleCustomerTypeList(_type, page, user)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

//修改客户流水
func (this *CustomerController) AlertComment() {
	commentId, _ := this.GetInt("commentId")
	comments := this.GetString("comments")
	ok := customerModel.AlertComments(commentId, comments)
	this.JSonCode(ok)
}

/////////////////////////////////////////////////////////////////////
func (this *CustomerController) GetAllCustomerList() {
	account := this.GetSession(SESSION_USER).(*Account)
	GetAllCustomerInfo(account, this, 0)
}

func (this *CustomerController) GetAllCustomerBySort() {
	key := this.GetString("sortKey")
	index, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	list, page_info := customerModel.GetAllCustomerBySort(key, index, page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetAgentCustomerBySort() {
	key := this.GetString("sortKey")
	index, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	user := this.GetSession(SESSION_USER).(*Account)
	list, page_info := customerModel.GetAgentCustomerBySort(key, index, page, user)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetSaleCustomerBySort() {
	key := this.GetString("sortKey")
	index, _ := this.GetInt("sortIndex", 0)
	page, _ := this.GetInt("page", 1)
	user := this.GetSession(SESSION_USER).(*Account)
	list, page_info := customerModel.GetSaleCustomerBySort(key, index, page, user)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	rsp["page_info"] = page_info
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetCustomerListForSale() {
	account := this.GetSession(SESSION_USER).(*Account)
	GetCustomerInfo(account, this, 2)
}

func (this *CustomerController) GetEmpCustomerList() {
	user := this.GetSession(SESSION_USER).(*Account)
	page, _ := this.GetInt("page")
	customers := customerModel.ShowEmployeeKeyWordCustomer("", user, page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = customers

	var page_total, record int
	record, page_total = customerModel.GetEmpCountInfo(user.AccountId, "")
	rsp["page_info"] = NewPageInfo(page, page_total, record)

	this.Data["json"] = rsp
	this.ServeJSON()
}

// 筛选
func (this *CustomerController) FilterCustomerSuperList() {
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
	customers := customerModel.FilterAllCustomer(filter)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = customers

	var page_total, record int
	record, page_total = customerModel.FilterAllCustomerCount(filter)
	rsp["page_info"] = NewPageInfo(page, page_total, record)

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) FilterCustomerAgentList() {
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

	customers := customerModel.FilterAgentCustomer(user, filter)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = customers

	var page_total, record int
	record, page_total = customerModel.FilterAgentCustomerCount(user, filter)
	rsp["page_info"] = NewPageInfo(page, page_total, record)

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) FilterCustomerEmpList() {
	user := this.GetSession(SESSION_USER).(*Account)
	g := this.Input()
	sortKey := this.GetString("sortKey", "timex")
	page, _ := this.GetInt("page", 1)
	sortIndex, _ := this.GetInt("sortIndex", 0)
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

	customers := customerModel.FilterEmpCustomer(user, filter)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = customers

	var page_total, record int
	record, page_total = customerModel.FilterEmpCustomerCount(user, filter)
	rsp["page_info"] = NewPageInfo(page, page_total, record)

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (c *CustomerController) EmployeeKeyWordCustomer() {
	user := c.GetSession(SESSION_USER).(*Account)
	keyWord := c.GetString("keyWord")
	page, _ := c.GetInt("page")
	customers := customerModel.ShowEmployeeKeyWordCustomer(keyWord, user, page)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = customers

	var page_total, record int
	record, page_total = customerModel.GetEmpCountInfo(user.AccountId, keyWord)
	rsp["page_info"] = NewPageInfo(page, page_total, record)

	c.Data["json"] = rsp
	c.ServeJSON()
}

func (this *CustomerController) GetAgentCustomerList() {
	user := this.GetSession(SESSION_USER).(*Account)
	page, _ := this.GetInt("page")
	list := customerModel.AgentShowKeyWorldCustomer("", user, page)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list

	var page_total, record int
	record, page_total = customerModel.GetAdminCountInfo(strconv.Itoa(user.AgentId), "")
	rsp["page_info"] = NewPageInfo(page, page_total, record)

	this.Data["json"] = rsp
	this.ServeJSON()
}

func GetCustomerInfo(account *Account, c *CustomerController, role int) {
	page, _ := c.GetInt("page", 1)
	rsp := NewResponse(SYSTEMSUCCESS)
	customers := customerModel.AdminShowKeyWorldCustomer("", account, page)
	rsp["list"] = customers
	record, page_total := customerModel.GetCountInfo(strconv.Itoa(account.AgentId), "")
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	c.Data["json"] = rsp
	c.ServeJSON()
}
func GetAllCustomerInfo(account *Account, c *CustomerController, role int) {
	page, _ := c.GetInt("page", 1)
	rsp := NewResponse(SYSTEMSUCCESS)
	customers := customerModel.AdminShowKeyWorldCustomer("", account, page)
	rsp["list"] = customers
	record, page_total := customerModel.GetAdminCountInfo(strconv.Itoa(account.AgentId), "")
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	c.Data["json"] = rsp
	c.ServeJSON()
}

func (c *CustomerController) DistributionCustomer() {
	custom_id, _ := c.GetInt("customer")
	account_id, _ := c.GetInt("account_id")
	user := c.GetSession(SESSION_USER).(*Account)
	ok := customerModel.DistributionCustomer(user, custom_id, account_id)

	c.JSonCode(ok)
}

func (c *CustomerController) AlertCustomer() {
	user := c.GetSession(SESSION_USER).(*Account)
	g := c.Input()
	customer := &Customer{Province: &Province{}, City: &City{}}
	customer.RTXNum, _ = c.GetInt("buin")
	customer.CustomerId, _ = c.GetInt("customerId")
	customer.EntName = c.GetString("entName")
	customer.Contacts = c.GetString("contacts")
	customer.Phone = c.GetString("phone")
	customer.Mail = c.GetString("mail")
	customer.Mobile = c.GetString("mobile")
	customer.Province.Id, _ = c.GetInt("province")
	customer.City.Id, _ = c.GetInt("city")
	customer.Note = c.GetString("remarks")
	customer.EmpCount, _ = c.GetInt("emp_count")
	tags := c.Input()["tags[]"]
	customer.QQ = c.GetString("qq")
	last_follow_time := c.GetString("follow_time")
	customer.LastFollowTime = &last_follow_time
	tagsAlert, _ := c.GetInt("tagsAlert", 0)

	ok := customerModel.AlertCustomer(g, customer, user, tags, tagsAlert)

	c.JSonCode(ok)
}

func (c *CustomerController) AdminKeyWordCustomer() {
	user := c.GetSession(SESSION_USER).(*Account)
	keyWord := c.GetString("keyWord")
	page, _ := c.GetInt("page")
	customers := customerModel.AdminShowKeyWorldCustomer(keyWord, user, page)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = customers

	var page_total, record int
	record, page_total = customerModel.GetAdminCountInfo(strconv.Itoa(user.AgentId), keyWord)
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	c.Data["json"] = rsp
	c.ServeJSON()
}
func (c *CustomerController) AgentKeyWordCustomer() {
	user := c.GetSession(SESSION_USER).(*Account)
	keyWord := c.GetString("keyWord")
	page, _ := c.GetInt("page")
	customers := customerModel.AgentShowKeyWorldCustomer(keyWord, user, page)

	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = customers

	var page_total, record int
	record, page_total = customerModel.GetAgentCountInfo(user.AgentId, keyWord)
	rsp["page_info"] = NewPageInfo(page, page_total, record)
	c.Data["json"] = rsp
	c.ServeJSON()
}

func (c *CustomerController) DelCustomer() {
	m := c.Input()
	DelArr := m["CustomerIds[]"]
	user := c.GetSession(SESSION_USER).(*Account)
	ok := customerModel.DelCustomer(user, DelArr)

	c.JSonCode(ok)
}

func (c *CustomerController) AddCustomer() {
	user := c.GetSession(SESSION_USER).(*Account)
	customer := &Customer{
		Province: &Province{},
		City:     &City{},
	}
	m := c.Input()
	url_type, _ := c.GetInt("url_type")
	customer.RTXNum, _ = c.GetInt("buin")
	customer.EntName = c.GetString("entName")
	customer.Contacts = c.GetString("contacts")
	customer.Phone = c.GetString("phone")
	customer.Mobile = c.GetString("mobile")
	customer.QQ = c.GetString("qq")
	customer.Mail = c.GetString("mail")
	customer.Province.Id, _ = c.GetInt("province")
	customer.City.Id, _ = c.GetInt("city")
	customer.Note = c.GetString("remarks")
	customer.Assign_status, _ = c.GetInt("assign_status", 0)
	customer.EmpCount, _ = c.GetInt("emp_count")
	customer.Agent = &Agent{}
	tags := m["tags[]"]
	var ok bool
	switch url_type {
	case 0:
		customer.Agent.AgentId = NoAllAgentId
		ok = customerModel.AdminAddCustomer(customer, user, tags)
	case 1:
		customer.Mobile = c.GetString("phone")
		customer.Agent.AgentId = user.AgentId
		ok = customerModel.AgentAddCustomer(customer, user, tags)
	default:
		customer.Mobile = c.GetString("phone")
		customer.Agent.AgentId = user.AgentId
		ok = customerModel.EmpAddCustomer(customer, user, tags)
	}
	c.JSonCode(ok)
}

func (this *CustomerController) AllocationCustomer() {
	var rsp map[string]interface{}
	user := this.GetSession(SESSION_USER).(*Account)
	if user.RoleId == ROLE_EMPLOYEE {
		rsp = NewResponse(SYSTEMERROR)

	} else {
		m := this.Input()
		customers := m["CustomerIds[]"]
		employeeId := this.GetString("EmployeeId")
		ok := customerModel.AllocationCustomer(user, customers, employeeId)
		if !ok {
			rsp = NewResponse(SYSTEMERROR)
		}
		rsp = NewResponse(SYSTEMSUCCESS)
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) AllocationCustomerAgent() {
	var rsp map[string]interface{}
	user := this.GetSession(SESSION_USER).(*Account)
	if user.RoleId == ROLE_EMPLOYEE || user.RoleId == ROLE_AGENT {
		rsp = NewResponse(SYSTEMERROR)
	} else {
		m := this.Input()
		customers := m["CustomerIds[]"]
		agentId, _ := this.GetInt("AgentId")
		ok := customerModel.AllocationCustomerAgent(user, customers, agentId)
		if !ok {
			rsp = NewResponse(SYSTEMERROR)
		}
		rsp = NewResponse(SYSTEMSUCCESS)
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) ExcelCustomEmp() {
	f, _, err := this.GetFile("file")
	user := this.GetSession(SESSION_USER).(*Account)
	if err != nil {
		beego.Error(err)
		return
	}
	_byte, _ := ioutil.ReadAll(f)
	err = this._excelCustomEmp(_byte, user)
	defer func() {
		f.Close()
		code := SYSTEMSUCCESS
		if err != nil {
			code = SYSTEMERROR
		}
		rsp := NewResponse(code)
		if err != nil {
			rsp["msg"] = err.Error()
		}
		this.Data["json"] = rsp
		this.ServeJSON()
	}()
}

func (this *CustomerController) Prepare() {
	if this.GetSession(SESSION_USER) == nil {
		this.Redirect("login.html", 302)
		return
	}
	user := this.GetSession(SESSION_USER).(*Account)
	this.Data["user"] = user
}

// get方法权限判断
func (this *CustomerController) GetAuthority(aim string) {
	account := this.GetSession(SESSION_USER).(*Account)
	if account.RoleId == ROLE_EMPLOYEE { // 渠道的信息只能由超级管理员修改
		this.Redirect(NotAuthority, 302)
	} else {
		this.TplName = aim
	}
}

func (this *CustomerController) JSonCode(ok bool) {
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

func (this *CustomerController) GetAllProvince() {
	list := customerModel.GetAllProvince()
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetProvinceByKey() {
	key := this.GetString("key")
	list := customerModel.GetProvinceByKey(key)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetTagByKey() {
	key := this.GetString("key")
	list := customerModel.GetTagByKey(key)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetCity() {
	id, _ := this.GetInt("province_id")
	list := customerModel.GetCity(id)
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) GetAllTag() {
	list := customerModel.GetAllTags()
	rsp := NewResponse(SYSTEMSUCCESS)
	rsp["list"] = list
	this.Data["json"] = rsp
	this.ServeJSON()
}

func (this *CustomerController) ExcelModelDownLoad() {
	file, err := os.Open("./excel/demo.xlsx")

	if err != nil {
		beego.Error(err)
		this.Redirect(NotAuthority, 302)
		return
	}

	defer file.Close()
	fileName := path.Base("./excel/demo.xlsx")
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	this.Ctx.Output.Header("Content-Type", "application/octet-stream")
	this.Ctx.Output.Header("content-disposition", "attachment; filename=\""+fileName+"\"")
	_, err = io.Copy(this.Ctx.ResponseWriter, file)
	if err != nil {
		beego.Error(err)
		this.Redirect(NotAuthority, 302)
		return
	}
}

func (this *CustomerController) _excelCustomEmp(b []byte, acc *Account) (err error) {
	var xlFile *xlsx.File
	xlFile, err = xlsx.OpenBinary(b)
	if err != nil {
		beego.Error(err)
		return
	}
	customs := make([]*Customer, 0)
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			c := &Customer{
				Account: &Account{},
				City:    &City{},
				Tags:    make([]*Tag, 0),
			}

			for _i, cell := range row.Cells {
				t, err := cell.String()
				if err != nil {
					beego.Error(err)
					break
				}
				c = this._getInfoFromExcel(c, _i, t)
			}
			c.Account.AccountId = acc.AccountId
			if c.EntName == "" {
				continue
			}
			customs = append(customs, c)
		}
	}
	err = customerModel.ExcelCustomEmp(customs, acc)
	return err
}

func (this *CustomerController) _getInfoFromExcel(c *Customer, i int, text string) *Customer {
	switch i {
	case 0:
		c.EntName = text
	case 1:
		d, _ := strconv.Atoi(text)
		c.RTXNum = d
	case 2:
		c.Contacts = text
	case 3:
		c.Phone = text
	case 4:
		c.Mobile = text
	case 5:
		c.Mail = text
	case 6:
		c.QQ = text
	case 7:
		c.Note = text
	case 8:
		d, _ := strconv.Atoi(text)
		c.City.Id = d
	case 10:
		t := &Tag{}
		if text == "A" {
			t.TagId = 20

		} else if text == "B" {
			t.TagId = 21
		} else if text == "C" {
			t.TagId = 22
		} else {
			t.TagId = 28
		}
		c.Tags = append(c.Tags, t)
	}
	return c
}

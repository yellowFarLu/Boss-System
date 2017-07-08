package models

import (
	"database/sql"
	"errors"
	"time"

	"strconv"

	"fmt"

	"strings"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	. "xinda.im/boss/common"
)

type CustomerModel struct{}

const (
	CUSBAHAVEINSERT = "新增客户"
)

// 查看该客户名称是否已经存在
func (this *CustomerModel) CustomerIsExists(name string) (tag bool) {
	query := `select count(1) from t_customer where t_customer.name=?;`
	var a int
	err := Db.QueryRow(query, name).Scan(&a)

	if err != nil {
		beego.Error(err)
		return
	}

	if a != 0 {
		return
	}

	return true
}

//获取该客户的运营平台数据
func (this *CustomerModel) GetYesInfo(c *Customer) *Customer {
	c_info := customerCache.GetCustomInfo(c.RTXNum)
	c.Staff = int(c_info.Staff)
	c.Active = int(c_info.Active)

	c.EntYesterdayInfo = customerCache.GetCustomYesInfo(c.RTXNum)
	return c
}

// 详细页

// 增加备注
func (this *CustomerModel) AddComment(account *Account, comment *Comment) (tag bool) {
	comment.Committer = account.AccountId
	comment.Timex = GetCurrentTime()

	tx, err := Db.Begin()

	defer func() {
		if err == nil {
			tag = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err != nil {
		beego.Error(err)
		return
	}

	query := `insert into t_customer_comments(custom_id, committer, comments, timex) values(?, ?, ?, ?);`
	com := make(map[string]interface{}, 0)
	com["备注内容"] = comment.Comments
	_comment := NewComments(account, "", com)
	_, err = tx.Exec(query, comment.CustomerId, comment.Committer, _comment, comment.Timex)

	if err != nil {
		beego.Error(err)
		return
	}

	this.AlertLastFollowTime(comment.CustomerId, tx)
	return
}

// 修改流水
func (this *CustomerModel) AlertComments(commentId int, comments string) (ok bool) {
	query := `update t_customer_comments set comments = ? where t_customer_comments.comment_id = ?;`
	tx, err := Db.Begin()
	defer func() {
		if err == nil {
			ok = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return
	}

	_, err = tx.Exec(query, comments, commentId)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

// 获取流水
func (this *CustomerModel) GetComments(customerId, _type int, page int) (comments []*Comment) {
	offset := GetOffsetByPage(page)
	comments = make([]*Comment, 0)
	var query string
	if _type == 0 {
		query = `select t1.comment_id, t1.custom_id, t2.NAME, t1.comments, t1.timex ,t1.type
		from t_customer_comments t1 LEFT JOIN t_account t2 on t1.committer = t2.account_id where t1.custom_id = ?
		order by t1.timex desc, t1.comment_id desc limit ?, ?;`
	} else {
		query = `select t1.comment_id, t1.custom_id, t2.NAME, t1.comments, t1.timex , t1.type
		from t_customer_comments t1 LEFT JOIN t_account t2 on t1.committer = t2.account_id where t1.custom_id = ? and 
		t1.type <>1 order by t1.timex desc, t1.comment_id desc limit ?, ?;`
	}
	rows, err := Db.Query(query, customerId, offset, 5)
	if err != nil {
		beego.Error(err)
		return
	}

	for rows.Next() {
		comment := &Comment{}
		err = rows.Scan(&comment.CommentId, &comment.CustomerId, &comment.Committer, &comment.Comments, &comment.Timex, &comment.Type)
		comments = append(comments, comment)

		if err != nil {
			beego.Error(err)
			return
		}
	}

	return
}

func (this *CustomerModel) CheckCustomerInfo(oldCus *Customer) (customer *Customer) {
	customer = oldCus
	if customer.LastFollowTime == nil {
		lst := ""
		customer.LastFollowTime = &lst
	}
	return
}

//获取Tags
func (this *CustomerModel) GetTagsForCustmer(oldCus *Customer) (customer *Customer) {
	customer = oldCus
	tags := this.GetTagList(customer.CustomerId)
	customer.Tags = tags

	customer = this.CheckCustomerInfo(customer)

	return this.GetYesInfo(customer)
}

func (this *CustomerModel) SuperCustomerDetail(customerId string) (customer *Customer) {
	customer = &Customer{
		Agent:    &Agent{Manager: &Account{}},
		Account:  &Account{},
		Province: &Province{},
		City:     &City{},
		Tags:     make([]*Tag, 0),
		Comments: make([]*Comment, 0),
	}
	var account_id sql.NullString
	query := `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 and t1.custom_id = ?;`

	err := Db.QueryRow(query, customerId).Scan(&customer.CustomerId, &customer.EntName, &customer.RTXNum, &customer.Contacts,
		&customer.Phone, &customer.Mobile, &customer.QQ, &customer.Mail, &customer.Agent.AgentId, &customer.Timex, &customer.Assign_status, &customer.LastFollowTime,
		&customer.Note, &customer.Agent.Name, &account_id, &customer.City.Id, &customer.City.Name, &customer.Province.Id, &customer.Province.Name, &customer.EmpCount)
	if err != nil {
		beego.Error(err)
		return
	}
	customer.Account.AccountId = GetNullString(account_id)

	return this.GetTagsForCustmer(customer)
}

func (this *CustomerModel) EmpCustomerDetail(customerId string, account *Account) (customer *Customer) {
	customer = &Customer{
		Agent:    &Agent{},
		Account:  &Account{},
		Province: &Province{},
		City:     &City{},
		Tags:     make([]*Tag, 0),
		Comments: make([]*Comment, 0),
	}
	var last_follow_time, account_id sql.NullString
	query := `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 and t1.custom_id = ? and t3.account_id = ?;`

	err := Db.QueryRow(query, customerId, account.AccountId).Scan(&customer.CustomerId, &customer.EntName, &customer.RTXNum, &customer.Contacts,
		&customer.Phone, &customer.Mobile, &customer.QQ,
		&customer.Mail, &customer.Agent.AgentId, &customer.Timex, &customer.Assign_status, &last_follow_time, &customer.Note,
		&account_id, &customer.City.Id, &customer.City.Name, &customer.Province.Id, &customer.Province.Name, &customer.EmpCount)
	if err != nil {
		beego.Error(err)
		return
	}
	lft := GetNullString(last_follow_time)
	customer.LastFollowTime = &lft
	customer.Account.AccountId = GetNullString(account_id)
	customer.Tags = this.GetTagList(customer.CustomerId)
	customer = this.CheckCustomerInfo(customer)

	return this.GetYesInfo(customer)
}

func (this *CustomerModel) AgentCustomerDetail(customerId string, account *Account) (customer *Customer) {
	customer = &Customer{
		Agent:    &Agent{},
		Account:  &Account{},
		Province: &Province{}, City: &City{}, Tags: make([]*Tag, 0)}
	var last_follow_time, account_id sql.NullString
	query := `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 and t1.agent_id = ? and t1.custom_id = ?;`

	err := Db.QueryRow(query, account.AgentId, customerId).Scan(&customer.CustomerId, &customer.EntName, &customer.RTXNum, &customer.Contacts,
		&customer.Phone, &customer.Mobile, &customer.QQ, &customer.Mail, &customer.Agent.AgentId, &customer.Timex, &customer.Assign_status, &customer.LastFollowTime,
		&customer.Note, &customer.Agent.Name, &account_id, &customer.City.Id, &customer.City.Name, &customer.Province.Id, &customer.Province.Name, &customer.EmpCount)
	if err != nil {
		beego.Error(err)
		return
	}
	lft := GetNullString(last_follow_time)
	customer.LastFollowTime = &lft
	customer.Account.AccountId = GetNullString(account_id)
	customer.Tags = this.GetTagList(customer.CustomerId)
	customer = this.CheckCustomerInfo(customer)

	return this.GetYesInfo(customer)
}

// 判断用户id对应的用户是不是这个账号的
func (this *CustomerModel) CustomerBelongAccount(customerId string, account *Account) (ok bool) {
	query := `select count(1) from t_customer t1 where t1.custom_id 
		in(select t2.custom_id from t_account_customers t2 where t2.account_id=?) and t1.custom_id = ?;`
	var value int
	err := Db.QueryRow(query, account.AccountId, customerId).Scan(&value)
	if err != nil {
		beego.Error(err)
		return
	}

	if value != 0 {
		ok = true
	}

	return
}

// 判断用户id对应的渠道是不是当前渠道
func (this *CustomerModel) CustomerBelongAgent(customerId string, account *Account) (ok bool) {
	query := `select count(1) from t_customer t1 where t1.custom_id = ? and t1.agent_id=?;`
	var value int
	err := Db.QueryRow(query, customerId, account.AgentId).Scan(&value)
	if err != nil {
		beego.Error(err)
		return
	}

	if value != 0 {
		ok = true
	}

	return
}

//超级管理员和渠道的筛选查询封装
func (this *CustomerModel) SuperAndAgentQuery(sqlStr string) (slice []*Customer) {
	slice = make([]*Customer, 0)

	rows, err := Db.Query(sqlStr)
	if err != nil {
		beego.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		c := &Customer{
			Agent:            &Agent{Manager: &Account{}},
			Account:          &Account{},
			Province:         &Province{},
			City:             &City{},
			EntYesterdayInfo: &EntYesterdayInfo{},
			Tags:             make([]*Tag, 0),
		}
		var last_follow_time, account_id, account_name sql.NullString
		err := rows.Scan(&c.CustomerId, &c.EntName, &c.RTXNum, &c.Contacts, &c.Phone, &c.Mobile, &c.QQ,
			&c.Mail, &c.Agent.AgentId, &c.Timex, &c.Assign_status, &last_follow_time, &c.Note,
			&c.Agent.Name, &account_id, &c.City.Id, &c.City.Name, &c.Province.Id, &c.Province.Name, &account_name, &c.EmpCount)
		if err != nil {
			beego.Error(err)
			return
		}
		lft := GetNullString(last_follow_time)
		c.LastFollowTime = &lft
		c.Account.AccountId = GetNullString(account_id)
		c.Account.Name = GetNullString(account_name)
		c.Tags = this.GetTagList(c.CustomerId)

		c = this.GetYesInfo(c)

		slice = append(slice, c)
	}

	return slice
}

//企业号或企业名称
func (this *CustomerModel) ParseEntNameAndRTX(filter *FilterParam) *FilterParam {
	if len(filter.G["keyword"]) != 0 {
		keyWorld := filter.G["keyword"][0]
		delete(filter.G, "keyword")

		if keyWorld != "" {
			_tem := `%` + keyWorld + `%`
			entArr := make([]string, 0)
			entArr = append(entArr, _tem)
			nameArr := make([]string, 0)
			nameArr = append(entArr, _tem)

			filter.G["ent_name"] = entArr
			filter.G["rtx_number"] = nameArr
		}
	}

	if len(filter.G["account_name"]) != 0 {
		filter.G["account_name"][0] = (`%` + filter.G["account_name"][0] + `%`)
	}

	return filter
}

// 超级管理员筛选功能
func (this *CustomerModel) FilterAllCustomer(filter *FilterParam) (slice []*Customer) {
	slice = make([]*Customer, 0)

	filter = this.ParseEntNameAndRTX(filter)
	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return slice
	}

	query := `select distinct t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	`
	var sqlStr = QuerySqlForMapArr(query, keyArr, CustomerSqlMap, *filter, TABLENAME_CUSTOMER)
	slice = this.SuperAndAgentQuery(sqlStr)
	return
}

//  超级管理员筛选功能获取数量
func (this *CustomerModel) FilterAllCustomerCount(filter *FilterParam) (record, page_total int) {
	filter = this.ParseEntNameAndRTX(filter)

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return 0, 0
	}
	filter.SortKey = "t1." + filter.SortKey
	query := `select distinct count(1)
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	`
	var sqlStr = QuerySqlForMapArr(query, keyArr, CustomerSqlMap, *filter, TABLENAME_CUSTOMER)

	sqlStr = strings.Split(sqlStr, " limit")[0]
	rows, err := Db.Query(sqlStr)

	return PageCount(err, rows)
}

//渠道筛选功能
func (this *CustomerModel) FilterAgentCustomer(account *Account, filter *FilterParam) (slice []*Customer) {
	slice = make([]*Customer, 0)

	filter = this.ParseEntNameAndRTX(filter)
	var agentId = fmt.Sprintf("%d", account.AgentId)
	agentIdArr := make([]string, 0)
	agentIdArr = append(agentIdArr, agentId)
	filter.G["agent_id"] = agentIdArr

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return slice
	}
	filter.SortKey = "t1." + filter.SortKey
	query := `select distinct t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	`
	var sqlStr = QuerySqlForMapArr(query, keyArr, CustomerSqlMap, *filter, TABLENAME_CUSTOMER)
	slice = this.SuperAndAgentQuery(sqlStr)
	return
}

//渠道筛选功能获取数目
func (this *CustomerModel) FilterAgentCustomerCount(account *Account, filter *FilterParam) (record, page_total int) {
	filter = this.ParseEntNameAndRTX(filter)

	var agentId = fmt.Sprintf("%d", account.AgentId)
	agentIdArr := make([]string, 0)
	agentIdArr = append(agentIdArr, agentId)
	filter.G["agent_id"] = agentIdArr

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return 0, 0
	}

	query := `select distinct count(1)
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	`
	var sqlStr = QuerySqlForMapArr(query, keyArr, CustomerSqlMap, *filter, TABLENAME_CUSTOMER)
	sqlStr = strings.Split(sqlStr, " limit")[0]
	rows, err := Db.Query(sqlStr)
	return PageCount(err, rows)
}

// 销售筛选功能
func (this *CustomerModel) FilterEmpCustomer(account *Account, filter *FilterParam) (slice []*Customer) {
	slice = make([]*Customer, 0)
	filter = this.ParseEntNameAndRTX(filter)
	var accountId = fmt.Sprintf("%s", account.AccountId)
	accountIdArr := make([]string, 0)
	accountIdArr = append(accountIdArr, accountId)
	filter.G["account_id"] = accountIdArr

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return slice
	}
	// filter.SortKey = "t1." + filter.SortKey
	query := `select distinct t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	`
	var sqlStr = QuerySqlForMapArr(query, keyArr, CustomerSqlMap, *filter, TABLENAME_CUSTOMER)

	// beego.Info(sqlStr)

	rows, err := Db.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		beego.Error(err)
		return
	}

	for rows.Next() {
		c := &Customer{
			Agent:            &Agent{Manager: &Account{}},
			Account:          &Account{},
			Province:         &Province{},
			City:             &City{},
			EntYesterdayInfo: &EntYesterdayInfo{},
			Tags:             make([]*Tag, 0),
		}
		var last_follow_time, account_id, account_name sql.NullString
		err := rows.Scan(&c.CustomerId, &c.EntName, &c.RTXNum, &c.Contacts, &c.Phone, &c.Mobile, &c.QQ,
			&c.Mail, &c.Agent.AgentId, &c.Timex, &c.Assign_status, &last_follow_time, &c.Note,
			&account_id, &c.City.Id, &c.City.Name, &c.Province.Id, &c.Province.Name, &account_name, &c.EmpCount)
		if err != nil {
			beego.Error(err)
			return
		}
		lft := GetNullString(last_follow_time)
		c.LastFollowTime = &lft
		c.Account.AccountId = GetNullString(account_id)
		c.Account.Name = GetNullString(account_name)
		c.Tags = this.GetTagList(c.CustomerId)

		c = this.GetYesInfo(c)
		slice = append(slice, c)
	}

	return
}

// 销售筛选获取数量
func (this *CustomerModel) FilterEmpCustomerCount(account *Account, filter *FilterParam) (record, page_total int) {
	filter = this.ParseEntNameAndRTX(filter)
	var accountId = fmt.Sprintf("%s", account.AccountId)
	accountIdArr := make([]string, 0)
	accountIdArr = append(accountIdArr, accountId)
	filter.G["account_id"] = accountIdArr

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return 0, 0
	}

	filter.SortKey = "t1." + filter.SortKey
	query := `select distinct count(1)
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	`
	var sqlStr = QuerySqlForMapArr(query, keyArr, CustomerSqlMap, *filter, TABLENAME_CUSTOMER)

	sqlStr = strings.Split(sqlStr, " limit")[0]
	rows, err := Db.Query(sqlStr)
	if err != nil {
		beego.Error(err)
		return
	}

	return PageCount(err, rows)
}

//超级管理员和渠道的查询封装
func (this *CustomerModel) SuperAndAgentScane(rows *sql.Rows) (slice []*Customer) {
	defer rows.Close()
	slice = make([]*Customer, 0)
	for rows.Next() {
		c := &Customer{
			Agent:            &Agent{Manager: &Account{}},
			Account:          &Account{},
			Province:         &Province{},
			City:             &City{},
			EntYesterdayInfo: &EntYesterdayInfo{},
			Tags:             make([]*Tag, 0),
		}
		var last_follow_time, account_id, account_name sql.NullString
		err := rows.Scan(&c.CustomerId, &c.EntName, &c.RTXNum, &c.Contacts, &c.Phone, &c.Mobile, &c.QQ,
			&c.Mail, &c.Agent.AgentId, &c.Timex, &c.Assign_status, &last_follow_time, &c.Note,
			&c.Agent.Name, &account_id, &c.City.Id, &c.City.Name, &c.Province.Id, &c.Province.Name, &account_name, &c.EmpCount)
		if err != nil {
			beego.Error(err)
			return
		}
		lft := GetNullString(last_follow_time)
		c.LastFollowTime = &lft
		c.Account.AccountId = GetNullString(account_id)
		c.Account.Name = GetNullString(account_name)
		c.Tags = this.GetTagList(c.CustomerId)

		c = this.GetYesInfo(c)
		slice = append(slice, c)
	}

	return
}

//增加流水  behave行为
func (this *CustomerModel) InsertBizComments(tx *sql.Tx, g map[string][]string, customer *Customer, account *Account, tagArr []int, behave string) (rsp bool) {
	keyArr, err := ArrMapToMapArr(g)
	if err != nil {
		beego.Error(err)
		return
	}

	// 记录插入商机流水
	query := `insert into t_customer_comments (custom_id, committer, comments, timex) values (?, ?, ?, ?);`
	com := make(map[string]interface{}, 0)
	for _, v := range keyArr {
		for k, _v := range v {
			if behave == CUSBAHAVEINSERT && _v == "''" {
				continue
			}

			com[alertBizMap[k]] = _v
		}
	}
	comment := NewComments(account, behave, com)
	timex := GetCurrentTime()
	_, err = tx.Exec(query, customer.CustomerId, account.AccountId, comment, timex)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

// 超级管理员增加客户
func (this *CustomerModel) AdminAddCustomer(customer *Customer, account *Account, tags []string) (rsp bool) {
	rsp = false
	if !CheckCustomerInfo(customer) {
		return
	}
	tx, err := Db.Begin()
	defer func() {
		if err == nil {
			rsp = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return
	}
	// t_customer
	query := `insert into t_customer (name, rtx_number, contacts, phone, mobile, qq, mail, agent_id, city_id, timex, assign_status, note, emp_count) 
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	timex := time.Now().Format("2006-01-02 15:04:05")
	res, err := tx.Exec(query, customer.EntName, customer.RTXNum, customer.Contacts,
		customer.Phone, customer.Mobile, customer.QQ, customer.Mail, NoAllAgentId,
		customer.City.Id, timex, customer.Assign_status, customer.Note, customer.EmpCount)
	if err != nil {
		beego.Error(err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		beego.Error(err)
		return
	}

	//t_customer_tags
	tagArr := make([]int, 0)
	query = `insert into t_customer_tags (custom_id, tag_id, timex) values (?,?,?)`
	for _, v := range tags {
		var tag_id, _ = strconv.Atoi(v)
		_, err = tx.Exec(query, id, tag_id, timex)
		if err != nil {
			beego.Error(err)
			break
		}
		tagArr = append(tagArr, tag_id)
	}

	//t_customer_comments
	query = `insert into t_customer_comments (custom_id, committer, comments, timex, type) values (?, ?, ?, ?, 1);`
	m := this.ReMarkMap(customer)
	if len(tagArr) != 0 {
		m["标签"] = GetAllTagName(tagArr)
	}

	comment := NewComments(account, "新增客户", m)
	res, err = tx.Exec(query, id, account.AccountId, comment, timex)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

// 渠道增加客户
func (this *CustomerModel) AgentAddCustomer(customer *Customer, account *Account, tags []string) (rsp bool) {
	rsp = false
	if !CheckCustomerInfo(customer) {
		return
	}
	tx, err := Db.Begin()
	defer func() {
		if err == nil {
			rsp = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return
	}
	// t_customer
	query := `insert into t_customer (name, rtx_number, contacts, phone, mobile, qq, mail, agent_id, city_id, timex, assign_status, note, emp_count) 
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	timex := GetCurrentTime()
	res, err := tx.Exec(query, customer.EntName, customer.RTXNum, customer.Contacts, customer.Phone, customer.Mobile, customer.QQ,
		customer.Mail, account.AgentId, customer.City.Id, timex, customer.Assign_status, customer.Note, customer.EmpCount)
	if err != nil {
		beego.Error(err)
		return
	}
	id, err := res.LastInsertId()
	customer.CustomerId = int(id)
	if err != nil {
		beego.Error(err)
		return
	}

	//t_customer_tags
	tagArr := make([]int, 0)
	query = `insert into t_customer_tags (custom_id, tag_id, timex) values (?,?,?)`
	for _, v := range tags {
		var tag_id, _ = strconv.Atoi(v)
		_, err = tx.Exec(query, id, tag_id, timex)
		if err != nil {
			beego.Error(err)
			break
		}
		tagArr = append(tagArr, tag_id)
	}

	//t_customer_comments
	query = `insert into t_customer_comments (custom_id, committer, comments, timex, type) values (?, ?, ?, ?, 1);`
	com := this.ReMarkMap(customer)
	if len(tagArr) != 0 {
		com["标签"] = GetAllTagName(tagArr)
	}
	comment := NewComments(account, "新增客户", com)
	_, err = tx.Exec(query, id, account.AccountId, comment, timex)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

// 销售增加客户
func (this *CustomerModel) EmpAddCustomer(customer *Customer, account *Account, tags []string) (rsp bool) {
	tx, err := Db.Begin()
	if err != nil {
		beego.Error(err)
	}

	rsp = this.EmpAddCustomerWithTx(tx, customer, account, tags)

	defer func() {
		if rsp == true {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	return
}

func (this *CustomerModel) EmpAddCustomerWithTx(tx *sql.Tx, customer *Customer, account *Account, tags []string) (rsp bool) {
	var err error

	if customer.City.Id == 0 {
		customer.City.Id = 1
	}

	if !CheckCustomerInfo(customer) {
		err = fmt.Errorf("employee add customer error: customer info error")
		beego.Error(err)
		return
	}
	if err != nil {
		beego.Error(err)
		return
	}

	query := `insert into t_customer (name, rtx_number, contacts, phone, mobile, qq, mail, agent_id, city_id, timex, assign_status, note, emp_count) 
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	timex := GetCurrentTime()
	res, err := tx.Exec(query, customer.EntName, customer.RTXNum, customer.Contacts, customer.Phone, customer.Mobile, customer.QQ,
		customer.Mail, account.AgentId, customer.City.Id, timex, customer.Assign_status, customer.Note, customer.EmpCount)
	if err != nil {
		beego.Error(err, account.AgentId)
		return
	}
	id, err := res.LastInsertId()
	customer.CustomerId = int(id)
	if err != nil {
		beego.Error(err)
		return
	}

	//t_account_customers
	query = `insert into t_account_customers (custom_id, account_id, assign_time, status) values (?, ? ,?, 0)`
	_, err = tx.Exec(query, id, account.AccountId, timex)
	if err != nil {
		beego.Error(err)
		return
	}

	//t_customer_tags
	tagArr := make([]int, 0)
	query = `insert into t_customer_tags (custom_id, tag_id, timex) values (?,?,?)`
	for _, v := range tags {
		var tag_id, _ = strconv.Atoi(v)
		_, err = tx.Exec(query, id, tag_id, timex)
		if err != nil {
			beego.Error(err)
			return
		}
		tagArr = append(tagArr, tag_id)
	}

	query = `insert into t_customer_comments (custom_id, committer, comments, timex, type) values (?, ?, ?, ?, 1);`
	customer.City.Name = areaCache.GetCityByCityId(customer.City.Id)
	com := this.ReMarkMap(customer)
	if len(tagArr) != 0 {
		com["标签"] = GetAllTagName(tagArr)
	}
	comment := NewComments(account, "新增客户", com)
	_, err = tx.Exec(query, id, account.AccountId, comment, timex)
	if err != nil {
		beego.Error(err)
		return
	}

	return true
}

//分配客户给员工
func (this *CustomerModel) DistributionCustomer(account *Account, customer_id int, account_id int) (rsp bool) {
	rsp = false
	tx, err := Db.Begin()

	//出异常回滚
	defer func() {
		if err == nil {
			rsp = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return
	}

	query := `update t_account_customers set account_id = ? , set assign_time = ? where custom_id = ?;`
	timex := time.Now().Format("2006-01-02 15:04:05")
	_, err = tx.Exec(query, customer_id, timex, account_id)
	if err != nil {
		beego.Error(err)
		return
	}

	query = `insert into t_customer_assign_history (custom_id, assigner, agent_id, assignee, timex);`
	_, err = tx.Exec(query, customer_id, account.AccountId, account.AgentId, account_id, timex)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

func (this *CustomerModel) GetAllTagName(tagArr []int) (allTagName string) {
	_tagArr := tagsCache.GetTagList(tagArr)
	end := (len(_tagArr) - 1)
	for i, v := range _tagArr {
		if i == end {
			allTagName += v.Name
		} else {
			allTagName += (v.Name + ", ")
		}
	}

	return
}

// 修改客户数据
func (this *CustomerModel) AlertCustomer(g map[string][]string, customer *Customer, account *Account, tags []string, tagsAlert int) (rsp bool) {

	if !CheckCustomerInfo(customer) {
		beego.Error("customer Info error")
		return false
	}

	keyArr, err := ArrMapToMapArr(g)
	if err != nil {
		beego.Error(err)
		return
	}

	if strings.Compare(*customer.LastFollowTime, NULLSTR) == 0 {
		customer.LastFollowTime = nil
	}

	tx, err := Db.Begin()

	//出异常回滚
	defer func() {
		if err == nil {
			rsp = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return false
	}

	//t_customer
	query := `update t_customer  set `
	isHas := false
	for _, v := range keyArr {
		for k, _v := range v {
			if cusDbField[k] == "rtx_number" {
				_, err = strconv.Atoi(DelMao(_v))
				if err != nil {
					beego.Error(`rtx_number error `, DelMao(_v), " ", err)
					return
				}
			} else if cusDbField[k] == "city_id" {
				cyId, err := strconv.Atoi(DelMao(_v))
				if err != nil || (cyId == 0) {
					beego.Error(`city_id error `, DelMao(_v), " ", err)
					return
				}
			} else if cusDbField[k] == NULLSTR {
				continue
			}

			query += (cusDbField[k] + `=` + _v + `,`)
			isHas = true
		}
	}

	if isHas || tagsAlert == TAG_ALERT {
		if isHas {
			query = DelLastDou(query)
			query += ` where custom_id=?;`
			_, err = tx.Exec(query, customer.CustomerId)
			if err != nil {
				beego.Error(err)
				return false
			}
		}

		query = `insert into t_customer_comments (custom_id, committer, comments, timex, type) values (?, ?, ?, ?, 1);`
		m := make(map[string]interface{}, 0)
		for _, v := range keyArr {
			for k, _v := range v {
				if alertMap[k] == NULLSTR {
					continue
				}

				m[alertMap[k]] = _v
				if k == "city" {
					// 取出城市id对应的城市
					_v = strings.TrimPrefix(_v, "'")
					_v = strings.TrimSuffix(_v, "'")
					cid, _ := strconv.Atoi(_v)
					m[alertMap[k]] = areaCache.GetCityByCityId(cid)
				}
			}
		}

		timex := GetCurrentTime()
		if tagsAlert == TAG_ALERT {
			_query := `delete from t_customer_tags where custom_id =?`
			_, err = tx.Exec(_query, customer.CustomerId)
			if err != nil {
				beego.Error(err)
				return
			}

			tagArr := make([]int, 0)

			_query = `insert into t_customer_tags (custom_id, tag_id, timex) values(?,?,?)`
			for _, v := range tags {
				id, _ := strconv.Atoi(v)
				_, err = tx.Exec(_query, customer.CustomerId, id, timex)
				if err != nil {
					beego.Error(err)
					break
				}
				tagArr = append(tagArr, id)
			}

			tem := GetAllTagName(tagArr)
			if tem != NULLSTR {
				m["标签"] = tem
			}
		}

		comment := NewComments(account, "修改客户资料", m)

		_, err = tx.Exec(query, customer.CustomerId, account.AccountId, comment, timex)
		if err != nil {
			beego.Error(err)
			return
		}

		this.AlertLastFollowTime(customer.CustomerId, tx)
	}

	return
}

//修改了的信息进入备注
var alertMap = map[string]string{
	"rtx_num":   "RTX总机号",
	"contacts":  "联系人",
	"phone":     "手机号",
	"mobile":    "座机",
	"mail":      "邮箱",
	"city":      "城市",
	"qq":        "QQ",
	"remarks":   "备注",
	"entName":   "企业名字",
	"emp_count": "员工数目",
}

//取得客户数据库字段
var cusDbField = map[string]string{
	"rtx_num":   "rtx_number",
	"contacts":  "contacts",
	"phone":     "phone",
	"mobile":    "mobile",
	"mail":      "mail",
	"city":      "city_id",
	"qq":        "qq",
	"remarks":   "note",
	"entName":   "name",
	"emp_count": "emp_count",
}

// 备注map
func (this *CustomerModel) ReMarkMap(customer *Customer) (m map[string]interface{}) {
	m = make(map[string]interface{}, 0)
	m["客户id"] = customer.CustomerId
	if customer.RTXNum != 0 {
		m["企业号"] = customer.RTXNum
	}

	m["企业名称"] = customer.EntName

	if customer.Contacts != "" {
		m["联系人"] = customer.Contacts
	}

	if customer.Phone != "" {
		m["联系电话"] = customer.Phone
	}

	if customer.Mobile != "" {
		m["座机"] = customer.Mobile
	}

	if customer.Mail != "" {
		m["邮箱"] = customer.Mail
	}

	if customer.QQ != "" {
		m["QQ"] = customer.QQ
	}

	if customer.Note != "" {
		m["备注"] = customer.Note
	}

	if customer.LastFollowTime != nil {
		m["最后跟踪时间"] = customer.LastFollowTime
	}

	m["城市"] = customer.City

	return
}

// 删除客户
func (this *CustomerModel) DelCustomer(account *Account, DelArr []string) (rsp bool) {
	rsp = false
	tx, err := Db.Begin()

	defer func() {
		if err == nil {
			rsp = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return
	}
	timex := GetCurrentTime()
	query := `update t_customer set enabled=0 where custom_id = ?;`
	stmt, err := tx.Prepare(query)
	defer stmt.Close()
	if err != nil {
		beego.Error(err)
		return
	}

	query = `insert into t_customer_comments (custom_id, committer, comments, timex) values (?, ?, ?, ?);`
	stmt2, err2 := tx.Prepare(query)
	defer stmt2.Close()
	if err2 != nil {
		beego.Error(err2)
		return
	}
	for _, v := range DelArr {
		_v, _ := strconv.Atoi(v)
		_, err = stmt.Exec(_v)
		if err != nil {
			beego.Error(err)
			return
		}
		m := make(map[string]interface{}, 0)
		m["客户id"] = v
		comment := NewComments(account, "删除客户", m)
		_, err = stmt2.Exec(_v, account.AccountId, comment, timex)
		if err != nil {
			beego.Error(err)
			return
		}

		this.AlertLastFollowTime(_v, tx)
	}

	return
}

// 超级管理员查看客户列表（关键字或者全部)
func (this *CustomerModel) AdminShowKeyWorldCustomer(keyWord string, account *Account, page int) (slice []*Customer) {
	slice = make([]*Customer, 0)
	var offset int
	offset = GetOffsetByPage(page)
	keyWord = "%" + keyWord + "%"
	query := `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1
	and (t1.name like ? or t1.rtx_number like ?)
	order by t1.timex desc limit ?, ?;`
	rows, err := Db.Query(query, keyWord, keyWord, offset, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}

	slice = this.SuperAndAgentScane(rows)

	return
}

// 渠道客户列表（关键字或者全部)
func (this *CustomerModel) AgentShowKeyWorldCustomer(keyWord string, account *Account, page int) (slice []*Customer) {
	slice = make([]*Customer, 0)
	var offset int
	offset = GetOffsetByPage(page)
	keyWord = "%" + keyWord + "%"
	query := `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 and t1.agent_id = ?
	and (t1.name like ? or t1.rtx_number like ?)
	order by t1.timex desc limit ?, ?;`

	rows, err := Db.Query(query, account.AgentId, keyWord, keyWord, offset, PAGESIZE)
	defer rows.Close()
	if err != nil {
		beego.Error(err)
		return
	}
	for rows.Next() {
		c := &Customer{
			Agent:            &Agent{Manager: &Account{}},
			Account:          &Account{},
			Province:         &Province{},
			EntYesterdayInfo: &EntYesterdayInfo{},
			City:             &City{},
			Tags:             make([]*Tag, 0),
		}
		var last_follow_time, account_id, account_name sql.NullString
		err := rows.Scan(&c.CustomerId, &c.EntName, &c.RTXNum, &c.Contacts, &c.Phone, &c.Mobile, &c.QQ,
			&c.Mail, &c.Agent.AgentId, &c.Timex, &c.Assign_status, &last_follow_time, &c.Note,
			&account_id, &c.City.Id, &c.City.Name, &c.Province.Id, &c.Province.Name, &account_name, &c.EmpCount)
		if err != nil {
			beego.Error(err)
			return
		}
		lft := GetNullString(last_follow_time)
		c.LastFollowTime = &lft
		c.Account.AccountId = GetNullString(account_id)
		c.Account.Name = GetNullString(account_name)
		c.Tags = this.GetTagList(c.CustomerId)
		c = this.GetYesInfo(c)
		slice = append(slice, c)
	}
	return
}

// 销售客户列表（关键字或者全部)
func (this *CustomerModel) ShowEmployeeKeyWordCustomer(keyWord string, account *Account, page int) (slice []*Customer) {
	slice = make([]*Customer, 0)
	keyWord = "%" + keyWord + "%"

	offset := GetOffsetByPage(page)

	query := `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name, t1.emp_count
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 and t3.account_id = ?
	and (t1.name like ? or t1.rtx_number like ?)
	order by t1.timex desc limit ?, ?;`
	rows, err := Db.Query(query, account.AccountId, keyWord, keyWord, offset, PAGESIZE)
	defer rows.Close()
	if err != nil {
		beego.Error(err)
		return
	}

	for rows.Next() {
		c := &Customer{
			Agent:            &Agent{Manager: &Account{}},
			Account:          &Account{},
			Province:         &Province{},
			City:             &City{},
			EntYesterdayInfo: &EntYesterdayInfo{},
			Tags:             make([]*Tag, 0),
		}
		var last_follow_time, account_id, account_name sql.NullString
		err := rows.Scan(&c.CustomerId, &c.EntName, &c.RTXNum, &c.Contacts, &c.Phone, &c.Mobile, &c.QQ,
			&c.Mail, &c.Agent.AgentId, &c.Timex, &c.Assign_status, &last_follow_time, &c.Note,
			&account_id, &c.City.Id, &c.City.Name, &c.Province.Id, &c.Province.Name, &account_name, &c.EmpCount)
		if err != nil {
			beego.Error(err)
			return
		}
		lft := GetNullString(last_follow_time)
		c.LastFollowTime = &lft
		c.Account.AccountId = GetNullString(account_id)
		c.Account.Name = GetNullString(account_name)
		c.Tags = this.GetTagList(c.CustomerId)

		c = this.GetYesInfo(c)
		slice = append(slice, c)
	}
	return
}

// 获取数量信息
func (this *CustomerModel) GetCountInfo(AgentId string, KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error
	if KeyWorld == "" { // 展本渠道所有客户
		rows, err = Db.Query("select count(1) from t_customer  where agent_id = ? and enabled=1", AgentId)
	} else { //通过BusinessName商业名字搜索
		key := "%" + KeyWorld + "%"
		rows, err = Db.Query(`select count(1) from t_customer t1
		natural join t_agent t2 
		where t2.agent_id = ? and t1.enabled=1 and 
		(t2.name like ? or t1.rtx_number like ?)`, AgentId, key, key)
	}
	return PageCount(err, rows)
}

// 获取超级管理员的客户的数量信息
func (this *CustomerModel) GetAdminCountInfo(AgentId string, KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error
	if KeyWorld == "" { // 所有客户
		rows, err = Db.Query("select count(1) from t_customer  where enabled=1")
	} else { //通过BusinessName商业名字搜索
		key := "%" + KeyWorld + "%"
		rows, err = Db.Query(`select count(1) from t_customer
		where enabled=1 and 
		(name like ? or rtx_number like ?)`, key, key)
	}
	return PageCount(err, rows)
}

// 获取渠道的客户的数量信息
func (this *CustomerModel) GetAgentCountInfo(AgentId int, KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error
	if KeyWorld == "" { // 展本渠道所有客户
		rows, err = Db.Query("select count(1) from t_customer  where enabled=1 and agent_id = ?", AgentId)
	} else { //通过BusinessName商业名字搜索
		key := "%" + KeyWorld + "%"
		rows, err = Db.Query(`select count(1) from t_customer
		where enabled=1 and agent_id = ? and
		(name like ? or rtx_number like ?)`, AgentId, key, key)
	}
	return PageCount(err, rows)
}

// 获取销售的客户的数量信息
func (this *CustomerModel) GetEmpCountInfo(account_id string, KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error
	if KeyWorld == "" { // 展本销售的所有客户
		rows, err = Db.Query(`select count(1) from t_customer t1 natural join t_account_customers t2
		where t2.account_id = ? and t1.enabled=1;`, account_id)
	} else { //通过BusinessName商业名字搜索
		key := "%" + KeyWorld + "%"
		rows, err = Db.Query(`select count(1) from t_customer t1 natural join t_account_customers t2 
		where t2.account_id = ? and t1.enabled=1 and (t1.name like ? or t1.rtx_number like ?) ;`, account_id, key, key)
	}
	return PageCount(err, rows)
}

//分配客户给销售
func (this *CustomerModel) AllocationCustomer(account *Account, customers []string, account_id string) (rsp bool) {
	rsp = false
	tx, err := Db.Begin()

	defer func() {
		beego.Error(err)
		if err == nil {
			rsp = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return
	}
	query := `insert into t_account_customers (custom_id,account_id,assign_time,status) values(?,?,?,?) on duplicate key update account_id = ?`
	_query := `insert into t_customer_comments (custom_id, committer, comments, timex) values (?, ?, ?, ?);`
	timex := GetCurrentTime()
	qh := `insert into t_customer_assign_history (custom_id, assigner, agent_id, assignee, timex) values(?, ?, ?, ?, ?);`

	for _, v := range customers {
		cId, _ := strconv.Atoi(v)
		_, err = tx.Exec(query, cId, account_id, timex, 0, account_id)
		if err != nil {
			beego.Error(err)
			break
		}
		com := make(map[string]interface{}, 0)
		com["客户id"] = cId
		com["分配销售id"] = account_id
		comment := NewComments(account, "分配客户到销售", com)
		_, err = tx.Exec(_query, cId, account.AccountId, comment, timex)
		if err != nil {
			beego.Error(err)
			break
		}

		_, err = tx.Exec(qh, cId, account.AccountId, account.AgentId, account_id, timex)
		if err != nil {
			beego.Error(err)
			return
		}
	}

	return
}

// 分配客户给渠道
func (this *CustomerModel) AllocationCustomerAgent(account *Account, customers []string, agentId int) (rsp bool) {
	rsp = false
	tx, err := Db.Begin()

	defer func() {
		if err == nil {
			rsp = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error(err)
		return
	}
	query := `update t_customer set agent_id = ? where custom_id =?;`
	_query := `insert into t_customer_comments (custom_id, committer, comments, timex) values (?, ?, ?, ?);`
	timex := time.Now().Format("2006-01-02 15:04:05")
	for _, v := range customers {
		cId, _ := strconv.Atoi(v)
		_, err = tx.Exec(query, agentId, cId)
		if err != nil {
			beego.Error(err)
			return
		}
		com := make(map[string]interface{}, 0)
		com["客户id"] = cId
		com["分配渠道id"] = agentId
		comment := NewComments(account, "分配客户到渠道", com)
		_, err = tx.Exec(_query, cId, account.AccountId, comment, timex)
		if err != nil {
			beego.Error(err)
			return
		}

		this.AlertLastFollowTime(cId, tx)
	}

	return
}

// 修改最后跟踪时间
func (this *CustomerModel) AlertLastFollowTime(customerId int, tx *sql.Tx) {
	lst := GetCurrentTime()
	query := `update t_customer set last_follow_time=? where custom_id = ?;`
	_, err := tx.Exec(query, lst, customerId)
	if err != nil {
		beego.Error(err)

		return
	}
}

func (this *CustomerModel) GetAllProvince() (list []*Province) {
	return areaCache.GetAllProvince()
}

func (this *CustomerModel) GetProvinceByKey(key string) (list []*Province) {
	return areaCache.GetProvinceByKey(key)
}

func (this *CustomerModel) GetTagByKey(key string) (list []*Tag) {
	return tagsCache.GetTagByKey(key)
}

func (this *CustomerModel) GetCity(id int) (list []*City) {
	return areaCache.GetCity(id)
}
func (this *CustomerModel) GetAllTags() (list []*Tag) {
	return tagsCache.GetAllTags()
}

func (this *CustomerModel) GetTagList(customer_Id int) []*Tag {
	ids := make([]int, 0)
	query := `select tag_id from t_customer_tags where custom_id = ?`
	rows, _ := Db.Query(query, customer_Id)
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return tagsCache.GetTagList(ids)
}

func (this *CustomerModel) ExcelCustomEmp(customs []*Customer, acc *Account) (err error) {
	tx, err := Db.Begin()
	if err != nil {
		beego.Error(err)
	}

	rc, err := redis.Dial(REDIS_PRODOCOL, REDIS_ADDRESS)
	if err != nil {
		beego.Error(err)
		return errors.New("add customer error")
	}
	defer rc.Close()

	for _, v := range customs {
		entNameRs, err := rc.Do("EXISTS", v.EntName)
		if err != nil {
			beego.Error(err)
			return errors.New("add customer error")
		}
		buinRs, err := rc.Do("EXISTS", v.RTXNum)
		if err != nil {
			beego.Error(err)
			return errors.New("add customer error")
		}
		// _, nameOk := customerCache.cusname_m[v.EntName]
		// _, buinOk := customerCache.cusbuin_m[v.RTXNum]
		var tema int64 = 1
		beego.Info(entNameRs.(int64))
		beego.Info(buinRs.(int64))
		if entNameRs.(int64) == tema {
			beego.Error(fmt.Sprintf("企业名重复：%s", v.EntName))
			return errors.New(fmt.Sprintf("企业名重复：%s", v.EntName))
		}

		if buinRs.(int64) == tema {
			beego.Error(fmt.Sprintf("RTX企业号重复： %d", v.RTXNum))
			return errors.New(fmt.Sprintf("RTX企业号重复： %d", v.RTXNum))
		}

		var tags = make([]string, 0)
		for _, _v := range v.Tags {
			tags = append(tags, fmt.Sprintf("%d", _v.TagId))
		}
		if ok := this.EmpAddCustomerWithTx(tx, v, acc, tags); !ok {
			beego.Error("add customer error: ", v)
			tx.Rollback()
			return errors.New("add customer error")
		}
	}

	tx.Commit()

	return nil
}

func (this *CustomerModel) GetAccountComments(id string, page int, Type string, account *Account) (list []*Comment, err error) {

	list = make([]*Comment, 0)

	if !accountModal.checkAccount(id, account) {
		return
	}

	var rows *sql.Rows
	offset := (page - 1) * COMMENTPAGESIZE
	var query string
	if Type == "custom" {
		query = `select t1.comments, t1.custom_id, t2.name, t1.timex from t_customer_comments t1 LEFT JOIN t_customer t2 on t1.custom_id = t2.custom_id 
	where t1.custom_id in (select custom_id from t_account_customers where account_id = ?) ORDER BY t1.timex desc limit ?, ?`
	} else {
		query = `select t1.comments, t1.biz_id, t2.title, t1.timex from t_biz_comments t1 LEFT JOIN t_biz t2 on t1.biz_id = t2.biz_id 
	where t2.custom_id in (select custom_id from t_account_customers where account_id = ?) ORDER BY t1.timex desc limit ?, ?;`
	}
	rows, err = Db.Query(query, id, offset, COMMENTPAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	for rows.Next() {
		c := &Comment{}
		rows.Scan(&c.Comments, &c.CustomerId, &c.Name, &c.Timex)
		list = append(list, c)
	}
	return
}

/**
* All Customer Sort By [Key & Index]
 */
func (this *CustomerModel) GetAllCustomerBySort(key string, index, page int) (list []*Customer, page_info map[string]interface{}) {
	list = make([]*Customer, 0)
	page_info = make(map[string]interface{}, 0)
	asc := fmt.Sprintf("%d", index)
	if v, ok := SortKeyMap[key]; ok {
		key = v
		if _v, _ok := SortIndexMap[asc]; _ok {
			asc = _v
		} else {
			asc = SortIndexMap["default"]
		}
	} else {
		key = SortKeyMap["default"]
		asc = SortIndexMap["default"]
	}

	_q := fmt.Sprintf("order by t1.%s %s limit ?, ?;", key, asc)
	query := `select count(1) from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1;`

	var record, page_total, to int
	err := Db.QueryRow(query).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 ` + _q

	rows, err := Db.Query(query, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetCustomListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
* Agent Customer Sort By [Key & Index]
 */
func (this *CustomerModel) GetAgentCustomerBySort(key string, index, page int, user *Account) (list []*Customer, page_info map[string]interface{}) {
	list = make([]*Customer, 0)
	page_info = make(map[string]interface{}, 0)
	asc := fmt.Sprintf("%d", index)
	if v, ok := SortKeyMap[key]; ok {
		key = v
		if _v, _ok := SortIndexMap[asc]; _ok {
			asc = _v
		} else {
			asc = SortIndexMap["default"]
		}
	} else {
		key = SortKeyMap["default"]
		asc = SortIndexMap["default"]
	}

	_q := fmt.Sprintf("order by t1.%s %s limit ?, ?;", key, asc)
	query := `select count(1) from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 and t2.agent_id = ?;`

	var record, page_total, to int
	err := Db.QueryRow(query, user.AgentId).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	WHERE t1.enabled = 1 and t2.agent_id = ? ` + _q

	rows, err := Db.Query(query, user.AgentId, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetCustomListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
* Sale Customer Sort By [Key & Index]
 */
func (this *CustomerModel) GetSaleCustomerBySort(key string, index, page int, user *Account) (list []*Customer, page_info map[string]interface{}) {
	list = make([]*Customer, 0)
	page_info = make(map[string]interface{}, 0)
	asc := fmt.Sprintf("%d", index)
	if v, ok := SortKeyMap[key]; ok {
		key = v
		if _v, _ok := SortIndexMap[asc]; _ok {
			asc = _v
		} else {
			asc = SortIndexMap["default"]
		}
	} else {
		key = SortKeyMap["default"]
		asc = SortIndexMap["default"]
	}

	_q := fmt.Sprintf("order by t1.%s %s limit ?, ?;", key, asc)
	query := `select count(1) from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	LEFT JOIN t_account_customers  t7 on t1.custom_id = t7.custom_id
	WHERE t1.enabled = 1 and t7.account_id = ?;`

	var record, page_total, to int
	err := Db.QueryRow(query, user.AccountId).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	LEFT JOIN t_account_customers  t7 on t1.custom_id = t7.custom_id
	WHERE t1.enabled = 1 and t7.account_id = ? ` + _q

	rows, err := Db.Query(query, user.AccountId, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetCustomListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
 * GetSaleCustomerTypeList
 */
func (this *CustomerModel) GetSaleCustomerTypeList(_type string, page int, user *Account) (list []*Customer, page_info map[string]interface{}) {
	var id int
	switch _type {
	case "A":
		id = 20
	case "B":
		id = 21
	case "C":
		id = 22
	default:
		id = 28
	}
	list = make([]*Customer, 0)
	page_info = make(map[string]interface{}, 0)

	query := `select count(1) from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	LEFT JOIN t_account_customers  t7 on t1.custom_id = t7.custom_id
	LEFT JOIN t_customer_tags t8 on t1.custom_id  = t8.custom_id
	WHERE t1.enabled = 1 and t8.tag_id = ? and t7.account_id = ?`

	var record, page_total, to int
	err := Db.QueryRow(query, id, user.AccountId).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.custom_id,t1.name,t1.rtx_number,t1.contacts,t1.phone,t1.mobile,t1.qq,t1.mail,t1.agent_id,t1.timex,t1.assign_status,t1.last_follow_time,t1.note,
	t2.name,t4.account_id,t5.city_id,t5.city,t6.province_id,t6.province, t4.name
	from t_customer t1 
	LEFT JOIN t_agent t2 on t1.agent_id = t2.agent_id
	LEFT JOIN t_account_customers t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t3.account_id = t4.account_id
	LEFT JOIN t_city t5 on t1.city_id = t5.city_id
	LEFT JOIN t_province t6 on t5.province_id = t6.province_id
	LEFT JOIN t_account_customers  t7 on t1.custom_id = t7.custom_id
	LEFT JOIN t_customer_tags t8 on t1.custom_id  = t8.custom_id
	WHERE t1.enabled = 1 and t8.tag_id = ? and t7.account_id = ? limit?, ?`
	rows, err := Db.Query(query, id, user.AccountId, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetCustomListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
* Get Customer List By Query [Rows]
 */
func (this *CustomerModel) GetCustomListByRows(list []*Customer, rows *sql.Rows) []*Customer {
	for rows.Next() {
		c := &Customer{
			Agent:            &Agent{Manager: &Account{}},
			Account:          &Account{},
			Province:         &Province{},
			City:             &City{},
			EntYesterdayInfo: &EntYesterdayInfo{},
			Tags:             make([]*Tag, 0),
		}
		var last_follow_time, account_id, account_name sql.NullString
		err := rows.Scan(&c.CustomerId, &c.EntName, &c.RTXNum, &c.Contacts, &c.Phone, &c.Mobile, &c.QQ,
			&c.Mail, &c.Agent.AgentId, &c.Timex, &c.Assign_status, &last_follow_time, &c.Note,
			&c.Agent.Name, &account_id, &c.City.Id, &c.City.Name, &c.Province.Id, &c.Province.Name, &account_name)
		if err != nil {
			beego.Error(err)
			return list
		}
		lft := GetNullString(last_follow_time)
		c.LastFollowTime = &lft
		c.Account.AccountId = GetNullString(account_id)
		c.Account.Name = GetNullString(account_name)
		c.Tags = this.GetTagList(c.CustomerId)

		c_info := customerCache.GetCustomInfo(c.RTXNum)
		c.Staff = int(c_info.Staff)
		c.Active = int(c_info.Active)

		c.EntYesterdayInfo = customerCache.GetCustomYesInfo(c.RTXNum)
		list = append(list, c)
	}
	rows.Close()
	return list
}

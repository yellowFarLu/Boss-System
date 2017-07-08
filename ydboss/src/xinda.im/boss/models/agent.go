package models

//渠道的数据库实现
import (
	"database/sql"
	"strings"

	"fmt"

	"github.com/astaxie/beego"
	. "xinda.im/boss/common"
)

type AgentModel struct{}

//筛选排序渠道
func (this *AgentModel) FilterAgentList(filter FilterParam) (slice []*Agent) {
	slice = make([]*Agent, 0)

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return slice
	}

	query := `select t1.agent_id, t1.name, 
		t1.contacts, t1.mobile, t1.Mail, t1.note, t1.timex from t_agent t1  `
	var sql = QuerySqlForMapArr(query, keyArr, AgentSqlMap, filter, TABLENAME_AGENT)

	beego.Info(sql)
	rows, err := Db.Query(sql)
	if err != nil {
		beego.Error("[agent db] db Query err: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		agent := &Agent{}
		err := rows.Scan(&agent.AgentId, &agent.Name, &agent.Contacts, &agent.Mobile, &agent.Mail, &agent.Note, &agent.Timex)
		if err != nil {
			beego.Error(err)
			return
		}
		slice = append(slice, agent)
	}
	return
}

//修改渠道信息
func (this *AgentModel) AlertAgentInfo(agent *Agent) (rsp bool) {
	rsp = false
	// 渠道信息有效性检查
	if !CheckAgentInfo(agent) {
		return
	}

	//开启事务
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

	// 更新渠道

	_, err = Db.Exec("update t_agent t1 set t1.name=?, t1.contacts=?, t1.mobile=?, t1.mail=?, t1.note=? where t1.agent_id = ?;",
		agent.Name, agent.Contacts, agent.Mobile, agent.Mail, agent.Note, agent.AgentId)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

// 查看渠道
func (this *AgentModel) GetAgentList(EmployeeId string, page int) (slice []*Agent) {
	var rows *sql.Rows
	var err error

	query := `select t1.agent_id, t1.name, t1.contacts, t1.mobile, t1.Mail, t1.note, t1.timex, t1.account_id from t_agent t1 
		where t1.enabled =1 and t1.agent_id <>1 order by t1.agent_id desc`
	if page == -1 {
		rows, err = Db.Query(query)
	} else {
		var offset int = 0
		offset = GetOffsetByPage(page)
		rows, err = Db.Query(query+` limit ?, ?;`, offset, PAGESIZE)
	}

	if err != nil {
		beego.Error(err)
	}
	defer rows.Close()

	slice = make([]*Agent, 0)
	for rows.Next() {
		agent := &Agent{Manager: &Account{}}
		err := rows.Scan(&agent.AgentId, &agent.Name, &agent.Contacts, &agent.Mobile, &agent.Mail, &agent.Note, &agent.Timex, &agent.Manager.AccountId)
		if err != nil {
			beego.Error(err)
		}

		slice = append(slice, agent)
	}

	return
}

// 增加一条渠道（分配这个渠道负责账号）
func (this *AgentModel) AddAgent(agent *Agent, account *Account) (tag bool) {
	var err error

	// 渠道信息有效性检查
	if !CheckAgentInfo(agent) {
		err = fmt.Errorf("agent info error")
		beego.Error(err)
		return
	}

	if !CheckAccountInfo(account) {
		err = fmt.Errorf("account info error")
		beego.Error(err)
		return
	}

	tx, err := Db.Begin()
	if err != nil {
		beego.Error("[agent] add faild: create tx error: ", err)
		return
	}
	defer func() {
		if err == nil {
			tag = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	//插入渠道表
	agent.Timex = GetCurrentTime()
	query := `INSERT INTO t_agent(name, contacts, mobile, mail, note, timex) VALUES (?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(query, agent.Name, agent.Contacts, agent.Mobile, agent.Mail, agent.Note, agent.Timex)
	if err != nil {
		beego.Error("[agent] add faild: insert agent to db error: ", err)
		return
	}
	agentId, _ := res.LastInsertId()

	account.Timex = GetCurrentTime()
	account.AgentId = int(agentId)
	account.RoleId = ROLE_EMPLOYEE
	account.Pwd = Sh1(account.Pwd)

	query = `insert into t_account(account_id, agent_id, pwd, name, gender, mobile, timex) 
	values (?, ?, ?, ?, ?, ?, ?);`
	_, err = tx.Exec(query, account.AccountId, account.AgentId, account.Pwd, account.Name, account.Gender, account.Mobile, account.Timex)
	if err != nil {
		beego.Error("insert into Account error:", err)
		return
	}

	//插入角色表
	_, err = tx.Exec(`insert into t_account_roles (role_id, account_id) values (?,?)`, ROLE_AGENT, account.AccountId)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

//删除渠道
func (this *AgentModel) DelAgent(DelArr []string) (ok bool) {
	//开启事务
	tx, err := Db.Begin()
	//出异常回滚
	defer func() {
		if err == nil {
			ok = true
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error("[agent db] db begin error:", err)
		return
	}

	for _, v := range DelArr {
		//信达渠道无法删除
		if strings.Compare(v, "1") == 0 || strings.Compare(v, "2") == 0 {
			err = fmt.Errorf("Error xinda or default agent can not delete")
			beego.Error(err)
			return
		}

		//处于渠道账号的agent_id要置为0
		_, err = tx.Exec(`update t_agent set enabled=0 where agent_id = ?;`, v)
		if err != nil {
			beego.Error("[agent db] db exec error:", err)
			return
		}

		_, err = tx.Exec(`update t_account as t1 set enabled = 0 where t1.agent_id = ?;`, v)
		if err != nil {
			beego.Error("[agent db] db update err", err)
			return
		}
	}

	return
}

// 搜索符合关键词的所有渠道
func (this *AgentModel) ShowKeyWorldAgent(KeyWorld string, page int) (slice []*Agent) {
	var offset int
	offset = GetOffsetByPage(page)
	KeyWorld = "%" + KeyWorld + "%"
	rows, err := Db.Query(`select t1.agent_id, t1.name, t1.contacts, t1.mobile, t1.Mail, t1.note, t1.timex from t_agent t1 
	where t1.enabled =1 and t1.name like ? and t1.agent_id <>1 order by t1.agent_id desc limit ?, ?;`, KeyWorld, offset, PAGESIZE)
	if err != nil {
		beego.Info("[agent db] db Query err: ", err)
	}
	defer rows.Close()

	slice = make([]*Agent, 0)
	for rows.Next() {
		agent := &Agent{}
		err := rows.Scan(&agent.AgentId, &agent.Name, &agent.Contacts, &agent.Mobile, &agent.Mail, &agent.Note, &agent.Timex)
		if err != nil {
			beego.Error(err)
		}

		slice = append(slice, agent)
	}

	return
}

// 获取渠道数量信息
func (this *AgentModel) GetCountInfo(KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error
	if KeyWorld == "" { // 获取所有渠道
		rows, err = Db.Query(`select count(1) from t_agent where enabled=1 and agent_id <>1;`)
	} else {
		KeyWorld = "%" + KeyWorld + "%"
		rows, err = Db.Query(`select count(1) from t_agent where Name like ? and enabled=1 and agent_id <>1;`, KeyWorld)
	}

	return PageCount(err, rows)
}

func (this *AgentModel) FilterAgentCount(filter FilterParam) (record, page_total int) {
	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return
	}
	var sql = QuerySqlForMapArr(`select count(1) from t_agent t1 `,
		keyArr, AgentSqlMap, filter, TABLENAME_AGENT)
	sql = strings.Split(sql, " limit")[0]
	beego.Info("计算个数的 ", sql)
	rows, err := Db.Query(sql)
	if err != nil {
		beego.Error(err)
		return
	}
	return PageCount(err, rows)
}

/**
* All agent Sort By [Key & Index]
 */
func (this *AgentModel) GetAgentBySort(key string, index, page int) (list []*Agent, page_info map[string]interface{}) {
	list = make([]*Agent, 0)
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
	query := `select count(1) from t_agent t1 
		where t1.enabled =1 and t1.agent_id <>1;`

	var record, page_total, to int
	err := Db.QueryRow(query).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.agent_id, t1.name, t1.contacts, t1.mobile, t1.Mail, t1.note, t1.timex, t1.account_id from t_agent t1 
		where t1.enabled =1 and t1.agent_id <>1  ` + _q

	rows, err := Db.Query(query, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetAgentListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
* Get Agent List By Query [Rows]
 */
func (this *AgentModel) GetAgentListByRows(list []*Agent, rows *sql.Rows) []*Agent {
	for rows.Next() {
		agent := &Agent{Manager: &Account{}}
		err := rows.Scan(&agent.AgentId, &agent.Name, &agent.Contacts, &agent.Mobile, &agent.Mail, &agent.Note, &agent.Timex, &agent.Manager.AccountId)
		if err != nil {
			beego.Error(err)
			return list
		}
		list = append(list, agent)
	}
	rows.Close()
	return list
}

package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"

	"strings"

	"github.com/astaxie/beego"
	. "xinda.im/boss/common"
)

//对字符串进行SHA1哈希
func Sh1(data string) string {
	t := sha1.New()
	io.WriteString(t, data+"xindaboss")
	return fmt.Sprintf("%x", t.Sum(nil))
}

type AccountModel struct{}

func (this *AccountModel) checkAccount(subAccName string, account *Account) (ok bool) {
	subAccount := accountModal.GetUserInfo(subAccName)
	if (subAccount.RoleId == ROLE_AGENT) || (subAccount.AgentId != account.AgentId) {
		return
	}

	return true
}

//销售筛选功能
func (this *AccountModel) FilterAccountList(filter FilterParam, account *Account) (slice []*Account) {
	slice = make([]*Account, 0)

	keyArr, err := this.addTerm(account, &filter)
	if err != nil {
		beego.Error(err)
		return slice
	}

	query := `select t1.account_id, t1.gender, t1.mobile, t1.name, t1.timex, t1.agent_id from t_account t1 `
	var sql = QuerySqlForMapArr(query, keyArr, AccountSqlMap, filter, TABLENAME_ACCOUNT)
	// beego.Info(sql)
	rows, err := Db.Query(sql)
	if err != nil {
		beego.Info("[account db] db Query err: ", err)
	}
	defer rows.Close()

	slice = make([]*Account, 0)
	for rows.Next() {
		e := &Account{}
		err := rows.Scan(&e.AccountId, &e.Gender,
			&e.Mobile, &e.Name, &e.Timex, &e.AgentId)
		if err != nil {
			beego.Error(err)
		}
		slice = append(slice, e)
	}

	return
}

//销售筛选功能获取个数
func (this *AccountModel) FilterAccountCount(filter FilterParam, account *Account) (record, page_total int) {
	keyArr, err := this.addTerm(account, &filter)
	if err != nil {
		beego.Error(err)
		return 0, 0
	}

	query := `select count(1) from t_account t1 `
	var sql = QuerySqlForMapArr(query,
		keyArr, AccountSqlMap, filter, TABLENAME_ACCOUNT)
	sql = strings.Split(sql, " limit")[0]
	// beego.Info(sql)
	rows, err := Db.Query(sql)
	if err != nil {
		beego.Info("[account db] db Query err: ", err)
	}

	return PageCount(err, rows)
}

func (this *AccountModel) addTerm(account *Account, filter *FilterParam) (keyArr []map[string]string, err error) {
	var agentId = fmt.Sprintf("%d", account.AgentId)
	agentIdArr := make([]string, 0)
	agentIdArr = append(agentIdArr, agentId)
	filter.G["agent_id"] = agentIdArr

	keyArr, err = ArrMapToMapArr(filter.G)

	return
}

//获取登录者的个人信息
func (this *AccountModel) GetUserInfo(Username string) (account *Account) {
	account = &Account{}
	// 查询登录的基本信息
	query := `select t1.account_id,t1.agent_id, t1.gender, t1.mobile, t1.name, 
	t1.pwd, t1.timex, t2.role_id from t_account t1 LEFT JOIN t_account_roles t2 ON t1.account_id=t2.account_id
	where t1.account_id = ?;`
	err := Db.QueryRow(query, Username).
		Scan(&account.AccountId, &account.AgentId, &account.Gender, &account.Mobile, &account.Name, &account.Pwd, &account.Timex, &account.RoleId)
	if err != nil {
		beego.Error(err)
	}
	return
}

//修改员工信息
func (this *AccountModel) AlertAccountInfo(account *Account) (tag bool) {
	if !CheckAccountInfo(account) {
		return false
	}

	query := `update t_account set name=?, gender=?, mobile=? where account_id=?;`
	_, err := Db.Exec(query, account.Name, account.Gender, account.Mobile, account.AccountId)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

// 删除员工
func (this *AccountModel) DelAccount(DelArr []string) (tag bool) {
	//开启事务
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

	for _, v := range DelArr {

		if strings.Compare(v, "admin") == 0 {
			err = fmt.Errorf("admin can not delete")
			beego.Error(err)
			return
		}

		//使用tx
		query := `update t_account set enabled=0 where account_id = ?;`
		_, err = tx.Exec(query, v)
		if err != nil {
			beego.Error(err)
			return
		}
	}
	return
}

// 搜索符合关键词的员工
func (this *AccountModel) ShowKeyWorldAccount(KeyWorld string, account *Account, page int) (slice []*Account) {
	slice = make([]*Account, 0)

	var offset int = 0
	offset = GetOffsetByPage(page)

	KeyWorld = "%" + KeyWorld + "%"
	query := `select t1.account_id, t1.gender, t1.mobile, t1.name, t1.timex, t1.agent_id from t_account t1
			where t1.agent_id = ? and t1.enabled = 1 and (t1.Name like ? or t1.account_id like ?)
			order by t1.timex desc limit ?, ?;`
	rows, err := Db.Query(query, account.AgentId, KeyWorld, KeyWorld, offset, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		e := &Account{}
		err := rows.Scan(&e.AccountId, &e.Gender,
			&e.Mobile, &e.Name, &e.Timex, &e.AgentId)
		if err != nil {
			beego.Error(err)
		}
		slice = append(slice, e)
	}
	return
}

//插入销售员工信息
func (this *AccountModel) AddAccount(account *Account) (tag bool) {

	if !CheckAccountInfo(account) {
		return
	}
	account.Timex = GetCurrentTime()
	account.RoleId = ROLE_EMPLOYEE

	//首先对密码加密
	account.Pwd = Sh1(account.Pwd)

	//开启事务
	tx, err := Db.Begin()

	//出异常回滚
	defer func() {
		if err == nil {
			tag = true
			tx.Commit()
		} else {
			beego.Error(err)
			tx.Rollback()
		}
	}()

	if err != nil {
		beego.Error("db begin error:", err)
		return
	}

	query := `insert into t_account(account_id, agent_id, pwd, name, gender, mobile, timex) values (?, ?,  ?, ?, ?, ?, ?);`
	_, err = tx.Exec(query, account.AccountId, account.AgentId,
		account.Pwd, account.Name, account.Gender, account.Mobile, account.Timex)
	if err != nil {
		beego.Error("insert into Account error:", err)
		return
	}

	//插入角色表
	query = `insert into t_account_roles (role_id, account_id) values (?,?)`
	_, err = tx.Exec(query, ROLE_EMPLOYEE, account.AccountId)
	if err != nil {
		beego.Error(err)
		return
	}

	return
}

// 获取数量信息
func (this *AccountModel) GetCountInfo(AgentId string, KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error
	if KeyWorld == "" {
		rows, err = Db.Query(`select count(1) from t_account 
	where agent_id = ? and enabled=1;`, AgentId)
	} else {
		KeyWorld = "%" + KeyWorld + "%"
		rows, err = Db.Query(`select count(1) from t_account 
	where agent_id = ? and enabled=1 and (Name like ? or account_id like ?);`, AgentId, KeyWorld, KeyWorld)
		if err != nil {
			beego.Error(err)
		}
	}

	return PageCount(err, rows)
}

// 获取该账号的客户列表
func (this *AccountModel) GetEmpCustomerList(account *Account, page, role int) (slice []*Customer, total, page_total int) {
	slice = make([]*Customer, 0)
	total = 0
	page_total = 0
	query := `select count(1) from t_customer t1 
	left join t_account_customers t2 on t1.custom_id = t2.custom_id 
	where t2.account_id = ? and t1.enabled=1;`
	err := Db.QueryRow(query, account.AccountId).Scan(&total)
	if err != nil {
		beego.Error(err)
		return
	}

	var to int
	to, page_total = PageInfoCount(page, total)
	query = `select t1.custom_id, t1.rtx_number, t1.name, t1.contacts, t1.phone, t1.mail, t1.timex, t1.note, t1.assign_status,
	t6.note, t7.note, t3.city, t4.province, t1.agent_id  from t_customer t1 
	left join t_account_customers t2 on t1.custom_id=t2.custom_id 
	left join t_city t3 on t1.city_id=t3.city_id
	left join t_province t4 on t3.province_id = t4.province_id
	left join t_customer_tags t5 on t1.custom_id = t5.custom_id
	left join t_tag t6 on t5.tag_id = t6.tag_id and t6.type like 'cusd'
	left join t_tag t7 on t5.tag_id = t7.tag_id and t7.type like 'cusc'
	where t2.account_id = ? and t1.enabled=1
	order by t1.custom_id limit ?, ?;`
	rows, err := Db.Query(query, account.AccountId, to, PAGESIZE)

	if err != nil {
		beego.Error(err)
		return
	}

	for rows.Next() {
		customer := &Customer{
			Account: &Account{},
			Agent:   &Agent{},
		}
		err := rows.Scan(&customer.CustomerId, &customer.RTXNum, &customer.EntName, &customer.Contacts, &customer.Phone, &customer.Mail,
			&customer.Timex, &customer.Note, &customer.Assign_status, &customer.Province.Name, &customer.City.Name, &customer.Agent.AgentId)
		if err != nil {
			beego.Error(err)
			return
		}
		slice = append(slice, customer)
	}

	rows.Close()

	return
}

// 查看是否已经存在账户信息
func (this *AccountModel) AcountIsExisted(account_id string) (tag bool) {
	rows, err := Db.Query("select * from t_account where account_id like ?", account_id)

	defer rows.Close()
	if err != nil {
		beego.Error(err)
	}

	if rows.Next() {
		tag = true
	}

	return
}

func (this *AccountModel) UpdateAccountPsw(id, new_psw string) (err error) {
	psw := Sh1(new_psw)
	query := `update t_account set pwd = ? where account_id = ?`
	_, err = Db.Exec(query, psw, id)
	if err != nil {
		beego.Error(err)
		return
	}
	return
}

func (this *AccountModel) CheckRightByAccount(user *Account, id string) (ok bool) {
	ok = false
	if user.RoleId > ROLE_AGENT {
		return
	}
	var agent_id int
	err := Db.QueryRow(`select agent_id from t_account where account_id = ?`, id).Scan(&agent_id)
	if err != nil {
		beego.Error(err)
		return
	}
	if agent_id != user.AgentId {
		return
	}
	return true
}

/**
* Agent Employee Sort By [Key & Index]
 */
func (this *AccountModel) GetAccountBySort(key string, index, page int, user *Account) (list []*Account, page_info map[string]interface{}) {
	list = make([]*Account, 0)
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
	query := `select count(1) from t_account t1 where t1.agent_id = ? and t1.enabled = 1;`

	var record, page_total, to int
	err := Db.QueryRow(query, user.AgentId).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.account_id, t1.gender, t1.mobile, t1.name, t1.timex, t1.agent_id from t_account t1
			where t1.agent_id = ? and t1.enabled = 1  ` + _q

	rows, err := Db.Query(query, user.AgentId, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetAccountListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
* Get Employee List By Query [Rows]
 */
func (this *AccountModel) GetAccountListByRows(list []*Account, rows *sql.Rows) []*Account {
	for rows.Next() {
		e := &Account{}
		err := rows.Scan(&e.AccountId, &e.Gender,
			&e.Mobile, &e.Name, &e.Timex, &e.AgentId)
		if err != nil {
			beego.Error(err)
			return list
		}
		list = append(list, e)
	}
	rows.Close()
	return list
}

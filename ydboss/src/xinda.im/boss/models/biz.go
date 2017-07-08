package models

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"

	"strings"

	. "xinda.im/boss/common"
)

type BizModel struct{}

const (
	BIZBAHAVEINSERT = "新增商机"
)

// 增加流水(封装)
func (this *BizModel) InsertCommentsFromC(content string, bizId int, account *Account) (rsp bool) {
	tx, err := Db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err != nil {
		beego.Error(err)
		return
	}

	query := `insert into t_biz_comments(biz_id, committer, comments, timex) values(?, ?, ?, ?);`
	com := make(map[string]interface{}, 0)
	com["备注内容"] = content
	comment := NewComments(account, "添加备注", com)
	timex := GetCurrentTime()
	_, err = tx.Exec(query, bizId, account.AccountId, comment, timex)
	if err != nil {
		beego.Error(err)
		return
	}
	return true
}

//增加流水  behave行为
func (this *BizModel) InsertBizComments(tx *sql.Tx, g map[string][]string, biz *Biz, account *Account, tagArr []int, behave string, tagsAlert int) (rsp bool) {
	keyArr, err := ArrMapToMapArr(g)
	if err != nil {
		beego.Error(err)
		return
	}

	// 记录插入商机流水
	query := `insert into t_biz_comments(biz_id, committer, comments, timex, type) values(?, ?, ?, ?, 1);`

	com := make(map[string]interface{}, 0)
	if tagsAlert == TAG_ALERT {
		com["标签"] = GetAllTagName(tagArr)
	}

	for _, v := range keyArr {
		for k, _v := range v {
			if behave == BIZBAHAVEINSERT && _v == "''" {
				continue
			}
			if k == "tagsAlert" {
				continue
			}
			com[alertBizMap[k]] = _v
			if k == "customerId" {
				_v = DelFu(_v)
				cid, _ := strconv.Atoi(_v)
				com[alertBizMap[k]] = this.GetCustomerNameById(cid)
			} else if k == "status" {
				_v = DelFu(_v)
				com[alertBizMap[k]] = statusMap[_v]
			}
		}
	}
	comment := NewComments(account, behave, com)

	timex := GetCurrentTime()
	_, err = tx.Exec(query, biz.BizId, account.AccountId, comment, timex)
	if err != nil {
		beego.Error(err)
	}
	return
}

func (this *BizModel) GetCustomerNameById(customerId int) (customerName string) {
	query := `select name from t_customer where custom_id = ?;`
	err := Db.QueryRow(query, customerId).Scan(&customerName)
	if err != nil {
		beego.Error(err)
	}

	return
}

// 商机状态
var statusMap = map[string]string{
	// 商机状态： 0初步交流 1需求沟通 2商务沟通 3签约交款 4商务失败
	"0": "初步交流",
	"1": "需求沟通",
	"2": "商务沟通",
	"3": "签约交款",
	"4": "商务失败",
}

//修改了的商机信息进入备注
var alertBizMap = map[string]string{
	"title":        "商机标题",
	"content":      "商机内容",
	"customerId":   "企业名称", // 根据企业id去拉企业名称
	"status":       "商机状态", //
	"quota":        "销售额",
	"real_time":    "真实签约时间",
	"estimatetime": "预计签约时间",
	"tagsAlert":    "tagsAlert",
}

//取得商机数据库字段
var bizDbField = map[string]string{
	"title":        "title",
	"content":      "content",
	"customerId":   "custom_id",
	"status":       "status", //
	"quota":        "amount",
	"real_time":    "real_time",
	"estimatetime": "estimate_time",
	"tagsAlert":    "tagsAlert",
}

// 销售获取商机详细页
func (this *BizModel) EmpBizDetail(bizId string, account *Account) (biz *Biz) {
	biz = &Biz{
		Tags:     make([]*Tag, 0),
		Customer: &Customer{},
	}
	query := `select t1.biz_id, t1.account_id, t1.custom_id, t1.title, t1.content, t1.amount, 
	t1.status, t1.timex, t1.estimate_time, t1.real_time, t3.name from t_biz t1
    left join t_customer t3 on t1.custom_id = t3.custom_id
    where t1.biz_id = ?;`
	err := Db.QueryRow(query, bizId).Scan(&biz.BizId, &biz.AccountId, &biz.Customer.CustomerId, &biz.Title, &biz.Content,
		&biz.Amount, &biz.Status, &biz.Timex, &biz.EstimateTime, &biz.RealTime, &biz.Customer.EntName)
	if err != nil {
		beego.Error(err)
	}

	biz, err = this.GetBizTags(biz)
	if err != nil {
		beego.Error("[opport] get tag erroer:", err)
	}

	if biz.RealTime == nil {
		tem := " "
		biz.RealTime = &tem
	}

	if biz.EstimateTime == nil {
		tem := " "
		biz.EstimateTime = &tem
	}

	return
}

// 客户标题或者客户名称
func (this *BizModel) ParseEntNameAndTitle(filter *FilterParam) *FilterParam {
	if len(filter.G["keyword"]) != 0 {
		keyWorld := filter.G["keyword"][0]
		delete(filter.G, "keyword")

		if keyWorld != "" {
			_tem := `%` + keyWorld + `%`

			entArr := make([]string, 0)
			entArr = append(entArr, _tem)

			titleArr := make([]string, 0)
			titleArr = append(entArr, _tem)

			filter.G["ent_name"] = entArr
			filter.G["title"] = titleArr
		}
	}

	return filter
}

//销售商机筛选
func (this *BizModel) FilterBizList(filter *FilterParam, account *Account) (slice []*Biz) {
	slice = make([]*Biz, 0)

	filter = this.ParseEntNameAndTitle(filter)

	// 是我的账号的才显示出来
	// G         map[string][]string
	accountArr := make([]string, 0)
	accountArr = append(accountArr, account.AccountId)
	filter.G["accountArr"] = accountArr

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return
	}
	filter.SortKey = "t1." + filter.SortKey
	query := `select t1.biz_id, t1.account_id, t1.custom_id, t1.title, t1.content, t1.amount, 
	t1.status, t1.timex, t1.estimate_time, t1.real_time, t3.name, t1.enabled from t_biz t1
    left join t_customer t3 on t1.custom_id = t3.custom_id
	`
	var sql = BizQuerySqlForMapArr(query, keyArr, BizSqlMap, *filter, TABLENAME_BIZ)
	beego.Info(sql)
	rows, err := Db.Query(sql)
	if err != nil {
		beego.Info("[biz db] db Query err: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		biz := &Biz{Tags: make([]*Tag, 0)}
		biz.Customer = &Customer{}
		err := rows.Scan(&biz.BizId, &biz.AccountId, &biz.Customer.CustomerId, &biz.Title, &biz.Content,
			&biz.Amount, &biz.Status, &biz.Timex, &biz.EstimateTime, &biz.RealTime, &biz.Customer.EntName, &biz.Enabled)
		if err != nil {
			beego.Error("[opport] get list erroer:", err)
			return
		}

		biz, err = this.GetBizTags(biz)
		if err != nil {
			beego.Error("[opport] get tag erroer:", err)
		}

		slice = append(slice, biz)
	}

	return
}

//渠道商机筛选
func (this *BizModel) FilterAgentBizList(filter *FilterParam, account *Account) (slice []*Biz) {
	slice = make([]*Biz, 0)

	filter = this.ParseEntNameAndTitle(filter)

	// 是这个渠道的才显示出来
	agentIdArr := make([]string, 0)
	agentIdArr = append(agentIdArr, strconv.Itoa(account.AgentId))
	filter.G["agent_id"] = agentIdArr

	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return
	}
	filter.SortKey = "t1." + filter.SortKey
	query := `select t1.biz_id, t1.account_id, t1.custom_id, t1.title, t1.content, t1.amount, 
	t1.status, t1.timex, t1.estimate_time, t1.real_time, t3.name, t4.account_id, t4.name, t1.enabled from t_biz t1
    left join t_customer t3 on t1.custom_id = t3.custom_id
	left join t_account t4 on t1.account_id = t4.account_id
	`
	var sql = BizQuerySqlForMapArr(query, keyArr, BizSqlMap, *filter, TABLENAME_BIZ)
	beego.Info(sql)
	rows, err := Db.Query(sql)
	if err != nil {
		beego.Info("[biz db] db Query err: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		biz := &Biz{Tags: make([]*Tag, 0)}
		biz.Customer = &Customer{}
		biz.Account = &Account{}
		err := rows.Scan(&biz.BizId, &biz.AccountId, &biz.Customer.CustomerId, &biz.Title, &biz.Content,
			&biz.Amount, &biz.Status, &biz.Timex, &biz.EstimateTime, &biz.RealTime, &biz.Customer.EntName,
			&biz.Account.AccountId, &biz.Account.Name, &biz.Enabled)
		if err != nil {
			beego.Error("[opport] get list erroer:", err)
			return
		}

		biz, err = this.GetBizTags(biz)
		if err != nil {
			beego.Error("[opport] get tag erroer:", err)
		}

		slice = append(slice, biz)
	}

	return
}

// 封装获取商数据
func (this *BizModel) GetBizs(rows *sql.Rows) (slice []*Biz) {
	slice = make([]*Biz, 0)

	for rows.Next() {
		biz := &Biz{Tags: make([]*Tag, 0)}

		biz.Customer = &Customer{}
		biz.Account = &Account{}
		err := rows.Scan(&biz.BizId, &biz.AccountId, &biz.Customer.CustomerId, &biz.Title, &biz.Content,
			&biz.Amount, &biz.Status, &biz.Timex, &biz.EstimateTime, &biz.RealTime, &biz.Customer.EntName,
			&biz.Account.AccountId, &biz.Account.Name, &biz.Enabled)
		if err != nil {
			beego.Error("[opport] get list erroer:", err)
			return
		}
		biz, err = this.GetBizTags(biz)
		if err != nil {
			beego.Error("[opport] get tag erroer:", err)
		}
		slice = append(slice, biz)
	}
	rows.Close()

	return
}

// 展示某个销售商机（全部或者搜索）
func (this *BizModel) ShowBizs(account *Account, key string, page int) (slice []*Biz) {
	slice = make([]*Biz, 0)

	var keyWord = "%" + key + "%"
	offset := GetOffsetByPage(page)

	query := `select t1.biz_id, t1.account_id, t1.custom_id, t1.title, t1.content, t1.amount, 
	t1.status, t1.timex, t1.estimate_time, t1.real_time, t3.name, t4.account_id, t4.name, t1.enabled from t_biz t1
    left join t_customer t3 on t1.custom_id = t3.custom_id
	left join t_account t4 on t1.account_id = t4.account_id
	where t1.account_id = ? and (t1.title like ? or t1.content like ? or t3.name like ?)
	order by t1.biz_id desc limit ?, ?;`
	rows, err := Db.Query(query, account.AccountId, keyWord, keyWord, keyWord, offset, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}

	return this.GetBizs(rows)
}

// 渠道展示商机（全部或者搜索）
func (this *BizModel) AgentShowBizs(account *Account, key string, page int) (slice []*Biz) {

	slice = make([]*Biz, 0)

	var keyWord = "%" + key + "%"
	offset := GetOffsetByPage(page)

	query := `select t1.biz_id, t1.account_id, t1.custom_id, t1.title, t1.content, t1.amount,
	t1.status, t1.timex, t1.estimate_time, t1.real_time, t3.name, t4.account_id, t4.name, t1.enabled from t_biz t1
	left join t_customer t3 on t1.custom_id = t3.custom_id
	left join t_account t4 on t1.account_id = t4.account_id
	where t4.agent_id = ? and (t1.title like ? or t1.content like ? or t3.name like ?) order by t1.biz_id desc limit ?, ?;`
	rows, err := Db.Query(query, account.AgentId, keyWord, keyWord, keyWord, offset, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}

	return this.GetBizs(rows)
}

//拉去商机对应的Tag
func (this *BizModel) GetBizTags(oldBbiz *Biz) (biz *Biz, err error) {
	biz = oldBbiz

	// 拉取对应的tag
	query := `select t_biz_tags.tag_id, name, type, note from t_tag, t_biz_tags where t_tag.tag_id=t_biz_tags.tag_id and t_biz_tags.biz_id = ?;`
	rows2, err := Db.Query(query, biz.BizId)
	if err != nil {
		beego.Error(err)
		return
	}
	for rows2.Next() {
		tag := &Tag{}
		err := rows2.Scan(&tag.TagId, &tag.Name, &tag.Type, &tag.Note)
		if err != nil {
			return nil, err
		}
		biz.Tags = append(biz.Tags, tag)
	}
	rows2.Close()

	return biz, nil
}

// 修改商机对应的流水
func (this *BizModel) AlertCommentsFromC(commentId int, comments string) (ok bool) {
	query := `update t_biz_comments set comments = ? where comment_id = ?;`
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

// 将商机阶段变成关闭
func (this *BizModel) DeleteBizs(account *Account, DelArr []string, status int) (ok bool) {
	//开启事务
	tx, err := Db.Begin()
	if err != nil {
		beego.Error(err)
		return
	}

	query := `insert into t_biz_comments(
		biz_id, committer, comments, timex) values(?, ?, ?, ?);`
	stmtL, err := tx.Prepare(query)
	defer stmtL.Close()
	if err != nil {
		beego.Error(err)
	}
	timex := GetCurrentTime()
	for _, v := range DelArr {
		id, _ := strconv.Atoi(v)
		_, err = Db.Exec(`update t_biz set enabled=0 where biz_id = ?`, id)
		if err != nil {
			beego.Error(err)
			break
		}

		_, err = Db.Exec(`update t_biz set status=? where biz_id = ?`, status, id)
		if err != nil {
			beego.Error(err)
			break
		}

		com := make(map[string]interface{}, 0)
		com["客户id"] = id
		comment := NewComments(account, "关闭商机", com)
		// 记录流水
		_, err := stmtL.Exec(id, account.AccountId, comment, timex)
		if err != nil {
			beego.Error(err)
			break
		}
	}

	if err == nil {
		ok = true
		tx.Commit()
	} else {
		tx.Rollback()
	}

	return
}

// 插入商机
func (this *BizModel) InsertBiz(g map[string][]string, account *Account, biz *Biz, tags []string) (tag bool) {

	if biz.RealTime != nil && strings.Compare(*biz.RealTime, "") == 0 {
		biz.RealTime = nil
	}

	if biz.EstimateTime != nil && strings.Compare(*biz.EstimateTime, "") == 0 {
		biz.EstimateTime = nil
	}

	if !CheckOpportunities(biz) {
		return false
	}

	biz.Timex = GetCurrentTime()

	tx, err := Db.Begin()
	if err != nil {
		beego.Error("[Biz] add faild: create tx error: ", err)
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

	//使用tx
	// 插入商机基本信息
	query := `insert into t_biz(account_id, custom_id, title, content, amount, status, 
	timex,
	estimate_time, real_time, agent_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	res, err := tx.Exec(query, biz.AccountId, biz.Customer.CustomerId, biz.Title, biz.Content,
		biz.Amount, biz.Status,
		biz.Timex, biz.EstimateTime, biz.RealTime, account.AgentId)
	if err != nil {
		beego.Error("[biz] add faild: insert agent to db error: ", err)
		return
	}

	id, err := res.LastInsertId()
	biz.BizId = int(id)

	// 插入tag
	err, tagArr := this.AlertTags(tx, biz, tags)
	if err != nil {
		beego.Error(err)
		return
	}
	query = `insert into t_customer_comments (custom_id, committer, comments, timex, type) values (?, ?, ?, ?,1);`
	m := make(map[string]interface{}, 0)
	m["商机Id"] = id
	m["商机标题"] = "<a href='/get_emp_biz_detail.html?id=" + strconv.Itoa(int(id)) + "'>" + biz.Title + "</a>"
	_comment := NewComments(account, "转化客户为商机", m)
	res, err = tx.Exec(query, biz.Customer.CustomerId, account.AccountId, _comment, biz.Timex)
	if err != nil {
		beego.Error(err)
		return
	}

	// 记录插入商机流水
	var tem int
	if len(tagArr) == 0 {
		tem = TAG_NOT_ALERT
	} else {
		tem = TAG_ALERT
	}
	tag = this.InsertBizComments(tx, g, biz, account, tagArr, BIZBAHAVEINSERT, tem)

	return
}

// 更新Tags
func (this *BizModel) AlertTags(tx *sql.Tx, biz *Biz, tags []string) (err error, tagArr []int) {
	query := `delete from t_biz_tags where biz_id =?`
	tagArr = make([]int, 0)
	_, err = tx.Exec(query, biz.BizId)
	if err != nil {
		return err, tagArr
	}

	timex := GetCurrentTime()
	query = `insert into t_biz_tags(biz_id, tag_id, timex) values(?, ?, ?);`
	for _, v := range tags {
		id, _ := strconv.Atoi(v)
		_, err = tx.Exec(query, biz.BizId, id, timex)
		if err != nil {
			return err, tagArr
		}

		tagArr = append(tagArr, id)
	}

	return nil, tagArr
}

// 更新商机（更新基本信息, 记流水账, 更新tag）
func (this *BizModel) AlertBiz(g map[string][]string, account *Account, biz *Biz, tags []string, tagsAlert int) (tag bool) {
	if !CheckOpportunities(biz) {
		return
	}
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
		beego.Error("[Biz] add faild: create tx error: ", err)
		return
	}

	keyArr, err := ArrMapToMapArr(g)
	if err != nil {
		beego.Error(err)
		return
	}

	// title=?, content=?, amount=?, status=?, real_time=? where t_biz.biz_id=?;
	isHas := false
	query := `update t_biz set `
	for _, v := range keyArr {
		for k, _v := range v {
			if bizDbField[k] == "tagsAlert" {
				continue
			}
			if bizDbField[k] == "status" {
				_, err := strconv.Atoi(DelFu(_v))
				if err != nil {
					return
				}
			} else if bizDbField[k] == "real_time" && _v == "''" {
				continue
			}
			query += (bizDbField[k] + `=` + _v + `,`)
			isHas = true
		}
	}

	if isHas == true || tagsAlert == TAG_ALERT {
		//把最后一个','去掉
		query = DelLastDou(query)
		query += ` where t_biz.biz_id=?;`
		beego.Info(query, biz.BizId)
		if isHas {
			_, err = tx.Exec(query, biz.BizId)
			if err != nil {
				beego.Error(err)
				return
			}
		}

		var tagArr []int
		if tagsAlert == TAG_ALERT {
			//insert tags
			err, tagArr = this.AlertTags(tx, biz, tags)
			if err != nil {
				beego.Error(err)
				return
			}
		}

		// 记录更新商机流水
		tag = this.InsertBizComments(tx, g, biz, account, tagArr, "修改商机", tagsAlert)
	}

	return
}

func (this *BizModel) GetBizComments(id, _type int, page int) ([]*Comment, error) {
	list := make([]*Comment, 0)
	offset := GetOffsetByPage(page)
	var query string
	if _type == 0 {
		query = `select t1.comment_id, t1.biz_id, t2.NAME, t1.comments, t1.timex, t1.type 
		from t_biz_comments t1 LEFT JOIN t_account t2 on t1.committer = t2.account_id where t1.biz_id = ? 
		order by t1.timex desc limit ?, ?;`
	} else {
		query = `select t1.comment_id, t1.biz_id, t2.NAME, t1.comments, t1.timex, t1.type
		from t_biz_comments t1 LEFT JOIN t_account t2 on t1.committer = t2.account_id where 
		t1.biz_id = ? and t1.type <>1 order by t1.timex desc limit ?, ?;`
	}
	rows, err := Db.Query(query, id, offset, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		c := &Comment{}
		err := rows.Scan(&c.CommentId, &c.CustomerId, &c.Committer, &c.Comments, &c.Timex, &c.Type)
		if err != nil {
			beego.Error(err)
			return list, err
		}
		list = append(list, c)
	}
	return list, nil
}

// 销售获取商机数量信息
func (this *BizModel) GetCountInfo(AccountId string, KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error

	KeyWorld = "%" + KeyWorld + "%"
	rows, err = Db.Query(`select count(1) from t_biz where account_id = ? and (title like ? or content like ?) and t_biz.status!=4;`, AccountId, KeyWorld, KeyWorld)

	return PageCount(err, rows)
}

// 渠道获取商机数量信息
func (this *BizModel) GetAgentBizCountInfo(AgentId int, KeyWorld string) (int, int) {
	var rows *sql.Rows
	var err error
	KeyWorld = "%" + KeyWorld + "%"

	query := `select count(1) from t_biz t1
	left join t_customer t3 on t1.custom_id = t3.custom_id
	left join t_account t4 on t1.account_id = t4.account_id
	where t4.agent_id = ? and (t1.title like ? or t1.content like ? or t3.name like ?) and t1.status!=4;`

	rows, err = Db.Query(query, AgentId, KeyWorld, KeyWorld, KeyWorld)

	return PageCount(err, rows)
}

//商机筛选数量
func (this *BizModel) FilterBizCount(filter FilterParam) (record, page_total int) {
	keyArr, err := ArrMapToMapArr(filter.G)
	if err != nil {
		beego.Error(err)
		return 0, 0
	}
	query := `select count(1) from t_biz t1
    left join t_customer t3 on t1.custom_id = t3.custom_id
	`
	var sql = BizQuerySqlForMapArr(query, keyArr, BizSqlMap, filter, TABLENAME_BIZ)
	sql = strings.Split(sql, " limit")[0]
	beego.Info(sql)
	rows, err := Db.Query(sql)
	if err != nil {
		beego.Info("[biz db] db Query err: ", err)
		return 0, 0
	}

	return PageCount(err, rows)
}

/**
*  Agent Opport Sort By [Key & Index]
 */
func (this *BizModel) GetAgentOpportBySort(key string, index, page int, user *Account) (list []*Biz, page_info map[string]interface{}) {
	list = make([]*Biz, 0)
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
	query := `select count(1) from t_biz t1
	LEFT JOIN t_customer t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t1.account_id = t4.account_id
	where t4.agent_id = ?  and t1.status!=4;`

	var record, page_total, to int
	err := Db.QueryRow(query, user.AgentId).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.biz_id, t1.account_id, t1.custom_id, t1.title, t1.content, t1.amount, 
	t1.status, t1.timex, t1.estimate_time, t1.real_time, t3.name, t4.NAME from t_biz t1
    	LEFT JOIN t_customer t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t1.account_id = t4.account_id
	where t4.agent_id = ?  and t1.status!=4  ` + _q

	rows, err := Db.Query(query, user.AgentId, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetOpportListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
*  Opport Sort By [Key & Index]
 */
func (this *BizModel) GetOpportBySort(key string, index, page int, user *Account) (list []*Biz, page_info map[string]interface{}) {
	list = make([]*Biz, 0)
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
	query := `select count(1) from t_biz t1
    	LEFT JOIN t_customer t3 on t1.custom_id = t3.custom_id
	where t1.account_id = ?  and t1.status!=4;`

	var record, page_total, to int
	err := Db.QueryRow(query, user.AccountId).Scan(&record)
	if err != nil {
		beego.Error(err)
		return
	}
	to, page_total = PageInfoCount(page, record)

	query = `select t1.biz_id, t1.account_id, t1.custom_id, t1.title, t1.content, t1.amount, 
	t1.status, t1.timex, t1.estimate_time, t1.real_time, t3.name, t4.NAME from t_biz t1
    	LEFT JOIN t_customer t3 on t1.custom_id = t3.custom_id
	LEFT JOIN t_account t4 on t1.account_id = t4.account_id
	where t1.account_id = ?  and t1.status!=4  ` + _q

	rows, err := Db.Query(query, user.AccountId, to, PAGESIZE)
	if err != nil {
		beego.Error(err)
		return
	}
	list = this.GetOpportListByRows(list, rows)
	page_info = NewPageInfo(page, page_total, record)
	return
}

/**
* Get Opport List By Query [Rows]
 */
func (this *BizModel) GetOpportListByRows(list []*Biz, rows *sql.Rows) []*Biz {
	for rows.Next() {
		biz := &Biz{
			Tags:    make([]*Tag, 0),
			Account: &Account{},
		}
		biz.Customer = &Customer{}
		err := rows.Scan(&biz.BizId, &biz.AccountId, &biz.Customer.CustomerId, &biz.Title, &biz.Content,
			&biz.Amount, &biz.Status, &biz.Timex, &biz.EstimateTime, &biz.RealTime, &biz.Customer.EntName, &biz.Account.Name)
		if err != nil {
			beego.Error(err)
			return list
		}
		biz, err = this.GetBizTags(biz)
		if err != nil {
			beego.Error(err)
			return list
		}
		list = append(list, biz)
	}
	rows.Close()
	return list
}

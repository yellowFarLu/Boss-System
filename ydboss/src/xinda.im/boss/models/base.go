package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"regexp"
	//"strconv"
	//	"strconv"
	"time"

	"github.com/astaxie/beego"

	"strings"

	. "xinda.im/boss/common"
)

// 表名
const (
	TABLENAME_ACCOUNT  = "account"
	TABLENAME_AGENT    = "agent"
	TABLENAME_BIZ      = "biz"
	TABLENAME_CUSTOMER = "customer"
)

// redis
const (
	REDIS_PRODOCOL = "tcp"
	REDIS_ADDRESS  = "localhost:6379"
)

// 标签是否修改了
const TAG_ALERT = 1
const TAG_NOT_ALERT = -1
const NULLSTR = ""

var (
	areaCache     *AreaCache
	tagsCache     *TagsCache
	customerCache *CustomerCache
	config        *Config
	Db            *sql.DB
	StatiscDb     *sql.DB //备份数据库
)

var (
	accountModal *AccountModel
)

type Config struct {
	MainDB   string `json:"mainDB"`
	LogName  string `json:"logName"`
	OpsdPath string `json:"opsdPath"`
}

var SortKeyMap map[string]string = map[string]string{
	"default":       "timex",
	"timex":         "timex",
	"follow_time":   "last_follow_time",
	"estimate_time": "estimate_time",
	"real_time":     "real_time",
}

var SortIndexMap map[string]string = map[string]string{
	"default": "desc",
	"1":       "desc",
	"0":       "asc",
}

// 渠道筛选Map
var AgentSqlMap map[string]string = map[string]string{
	"name":        "name like ",
	"contacts":    "contacts like ",
	"timex_begin": "timex >= ",
	"timex_end":   "timex <= ",
}

// 销售的筛选排序Map
var AccountSqlMap map[string]string = map[string]string{
	"name":        "name like ",
	"account_id":  "account_id like ",
	"timex_begin": "timex >= ",
	"timex_end":   "timex <= ",
	"agent_id":    "t1.agent_id = ",
}

//客户的筛选排序Map
var CustomerSqlMap map[string]string = map[string]string{
	"last_follow_time_begin": "last_follow_time >= ",
	"last_follow_time_end":   "last_follow_time <= ",
	"assign_status":          "assign_status = ",
	"rtx_number":             "rtx_number like ",
	"ent_name":               "t1.name like ",
	"timex_begin":            "t1.timex >= ",
	"timex_end":              "t1.timex <= ",
	"agent_id":               "t1.agent_id = ",
	"account_id":             "t3.account_id = ",

	"emp_count_from":   "t1.emp_count >= ",
	"emp_count_to":     "t1.emp_count <= ",
	"timex_from":       "t1.timex >= ",
	"timex_to":         "t1.timex <= ",
	"follow_time_from": "t1.last_follow_time >= ",
	"follow_time_to":   "t1.last_follow_time <= ",
	"account_name":     "t4.name like ",
}

// 商机的筛选排序Map
var BizSqlMap map[string]string = map[string]string{
	"title":               "title like ",
	"status":              "status = ",
	"estimate_time_begin": "estimate_time >= ",
	"estimate_time_end":   "estimate_time <= ",
	"real_time_begin":     "real_time >= ",
	"real_time_end":       "real_time <= ",
	"amount_max":          "amount <= ",
	"amount_min":          "amount >= ",
	"timex_begin":         "t1.timex >= ",
	"timex_end":           "t1.timex <= ",
	"accountArr":          "account_id like ",
	"agent_id":            "agent_id = ",

	//新的
	"sales_from":         "amount >= ",
	"sales_to":           "amount <= ",
	"timex_from":         "t1.timex >= ",
	"timex_to":           "t1.timex <= ",
	"estimate_time_from": "estimate_time >= ",
	"estimate_time_to":   "estimate_time <= ",
	"real_time_from":     "real_time >= ",
	"real_time_to":       "real_time <= ",
	"ent_name":           "t3.name like",
}

// 把JSON转化成map
func JsonToMap(js string) (rs map[string]string, err error) {
	if err := json.Unmarshal([]byte(js), &rs); err != nil {
		beego.Error(err)
		return nil, err
	} else {
		return rs, nil
	}
}

// 把JSON数组转化为Map数组
func JsonArrToMapArr(filterStrArr []string) (keyArr []map[string]string, err error) {
	keyArr = make([]map[string]string, 0)
	for _, v := range filterStrArr {
		temM, err := JsonToMap(v)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		beego.Info("转化后的Map：", temM)
		keyArr = append(keyArr, temM)
	}

	return keyArr, nil
}

// 把数组Map转化成Map数组
func ArrMapToMapArr(g map[string][]string) (keyArr []map[string]string, err error) {
	keyArr = make([]map[string]string, 0)
	for k, v := range g {
		if k == "sortKey" || k == "sortIndex" || k == "tags[]" ||
			k == "provinces[]" || k == "customerId" || k == "province" ||
			k == "opportunities_id" {
			continue
		}

		temM := make(map[string]string)
		temM[k] = "'" + v[0] + "'"
		keyArr = append(keyArr, temM)
	}

	return keyArr, nil
}

// 客户SQL语句拼接
func QuerySqlForMapArr(sql string, keyArr []map[string]string, keyMap map[string]string, filter FilterParam, tagTable string) (rs string) {

	var query bytes.Buffer
	_, err := query.WriteString(sql)
	if err != nil {
		beego.Error(err)
	}

	var tem string

	// 筛选tag或者省份
	if len(filter.Tags) != 0 {
		if tagTable == TABLENAME_CUSTOMER {
			tem = `left join t_customer_tags t7 on t1.custom_id = t7.custom_id
    		   left join t_tag t8 on t7.tag_id = t8.tag_id WHERE t1.enabled = 1 and t8.tag_id in(`
		}

	} else {

		if tagTable == TABLENAME_AGENT {
			tem = `where t1.enabled =1 and t1.agent_id!=1 and `
		} else {
			tem = `WHERE t1.enabled = 1 and `
		}
	}

	query.WriteString(tem)

	mo := len(filter.Tags) - 1
	for i, tag := range filter.Tags {
		tem = ""
		if i == mo {
			tem += (tag + ") and ")
		} else {
			tem += (tag + ",")
		}
		query.WriteString(tem)
	}

	mo = len(filter.Provinces) - 1
	if (mo + 1) != 0 {
		tem = `t6.province_id in(`
		query.WriteString(tem)

		for i, pro := range filter.Provinces {
			tem = ""
			if i == mo {
				tem += (pro + ") and ")
			} else {
				tem += (pro + ",")
			}
			query.WriteString(tem)
		}
	}

	for _, m := range keyArr {
		tem = ""

		for k, v := range m {
			if v == "''" || k == "agent_id" || k == "account_id" ||
				k == "rtx_number" || k == "ent_name" {
				continue
			}

			tem += (keyMap[k] + v + " and ")

		}
		query.WriteString(tem)
	}

	// 把最后一个and去掉
	tem = query.String()
	end := strings.LastIndex(tem, ` and `)
	tem = tem[0:end]

	query.Reset()
	query.WriteString(tem)

	// 查看是否要渠道或者销售id判定
	if len(filter.G["account_id"]) != 0 {
		temOne := (` and ` + CustomerSqlMap["account_id"] + `'` + filter.G["account_id"][0] + `'`)
		query.WriteString(temOne)
	}

	if len(filter.G["agent_id"]) != 0 {
		temOne := (` and ` + CustomerSqlMap["agent_id"] + `'` + filter.G["agent_id"][0] + `'`)
		query.WriteString(temOne)
	}

	if tagTable == TABLENAME_CUSTOMER && len(filter.G["ent_name"]) != 0 || len(filter.G["rtx_number"]) != 0 {
		temOne := (` and (` + keyMap["ent_name"] + `'` + filter.G["ent_name"][0] + `'`)
		temOne += (` or ` + keyMap["rtx_number"] + `'` + filter.G["ent_name"][0] + `' )`)
		query.WriteString(temOne)
	}

	desc := " desc"
	if filter.SortIndex != 0 {
		desc = " asc"
	}

	tem = ` order by ` + filter.SortKey + desc
	query.WriteString(tem)

	rs = query.String()

	return
}

// 商机SQL语句拼接
func BizQuerySqlForMapArr(sql string, keyArr []map[string]string, keyMap map[string]string, filter FilterParam, tagTable string) (rs string) {

	var query bytes.Buffer
	_, err := query.WriteString(sql)
	if err != nil {
		beego.Error(err)
	}

	var tem string

	// 筛选tag或者省份
	if len(filter.Tags) != 0 {
		tem = `left join t_biz_tags t7 on t1.biz_id = t7.biz_id
			   left join t_tag t8 on t7.tag_id = t8.tag_id WHERE t1.enabled = 1 and t8.tag_id in(`
	} else {

		if tagTable == TABLENAME_BIZ {
			tem = `where t1.status!=4 and `
		} else if tagTable == TABLENAME_AGENT {
			tem = `where t1.enabled =1 and t1.agent_id!=1 and `
		} else {
			tem = `WHERE t1.enabled = 1 and `
		}
	}

	query.WriteString(tem)

	mo := len(filter.Tags) - 1
	for i, tag := range filter.Tags {
		tem = ""
		if i == mo {
			tem += (tag + ") and ")
		} else {
			tem += (tag + ",")
		}
		query.WriteString(tem)
	}

	mo = len(filter.Provinces) - 1
	if (mo + 1) != 0 {
		tem = `t6.province_id in(`
		query.WriteString(tem)

		for i, pro := range filter.Provinces {
			tem = ""
			if i == mo {
				tem += (pro + ") and ")
			} else {
				tem += (pro + ",")
			}
			query.WriteString(tem)
		}
	}

	for _, m := range keyArr {
		tem = ""

		for k, v := range m {
			if v == "''" || k == "agent_id" || k == "account_id" ||
				k == "ent_name" || k == "title" {
				continue
			}

			tem += (keyMap[k] + v + " and ")

		}
		query.WriteString(tem)
	}

	// 把最后一个and去掉
	tem = query.String()
	end := strings.LastIndex(tem, ` and `)
	tem = tem[0:end]

	query.Reset()
	query.WriteString(tem)

	// 查看是否要渠道或者销售id判定
	if len(filter.G["account_id"]) != 0 {
		temOne := (` and ` + CustomerSqlMap["account_id"] + `'` + filter.G["account_id"][0] + `'`)
		query.WriteString(temOne)
	}

	if len(filter.G["agent_id"]) != 0 {
		temOne := (` and ` + CustomerSqlMap["agent_id"] + `'` + filter.G["agent_id"][0] + `'`)
		query.WriteString(temOne)
	}

	if tagTable == TABLENAME_BIZ && len(filter.G["ent_name"]) != 0 || len(filter.G["title"]) != 0 {
		temOne := (` and (` + keyMap["ent_name"] + `'` + filter.G["ent_name"][0] + `'`)
		temOne += (` or ` + keyMap["title"] + `'` + filter.G["title"][0] + `' )`)
		query.WriteString(temOne)
	}

	desc := " desc"
	if filter.SortIndex != 0 {
		desc = " asc"
	}

	tem = ` order by ` + filter.SortKey + desc
	query.WriteString(tem)

	rs = query.String()

	return
}

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//如果最后一个符号是',' 把最后一个','去掉
func DelLastDou(str string) (rs string) {
	end := strings.LastIndex(str, ",")
	if end == -1 {
		return str
	} else {
		return string(str[0:end])
	}

}

// 得到TAG的所有名字
func GetAllTagName(tagArr []int) (allTagName string) {
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

//去掉冒号
func DelFu(_v string) (rs string) {
	_v = strings.TrimPrefix(_v, "'")
	_v = strings.TrimSuffix(_v, "'")

	return _v
}

// 去掉符号 ''
func DelMao(str string) (rs string) {
	start := strings.Index(str, "'")
	end := strings.LastIndex(str, "'")

	return SubStr(str, start+1, end)
}

// 截取字符串
func SubStr(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

/*
有效性检查
*/
type RegexCheck struct{}

var Regex *RegexCheck

func CacheInit() {
	areaCache = NewAreaServer()
	tagsCache = NewTagsServer()
	customerCache = NewCustomerServer()
}

// 手机号码
func (Regex *RegexCheck) PhoneCheck(Phone string) bool {
	if Phone == "" {
		return true
	}

	reg := regexp.MustCompile(`^[1-9]\d*$`)
	return reg.MatchString(Phone)
}

// 邮箱
func (Regex *RegexCheck) EmailCheck(mail string) bool {
	if mail == "" {
		return true
	}
	reg := regexp.MustCompile(`^$|^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`)
	return reg.MatchString(mail)
}

// RTX
func (Regex *RegexCheck) RTXNumberCheck(RTX string) bool {
	if RTX == "" || RTX == "0" {
		return true
	}
	reg := regexp.MustCompile(`^$|^[1-9][0-9]{7}$`)
	return reg.MatchString(RTX)
}

// 商机类别判定
func (Regex *RegexCheck) OptCategoryCheck(CategoryDescription string) bool {
	// reg := regexp.MustCompile(`^[0-4]$`)
	// return reg.MatchString(CategoryDescription)
	return true
}

// 商机阶段判定
func (Regex *RegexCheck) OptStageCheck(OpportunitiesStage string) bool {
	// reg := regexp.MustCompile(`^[0-4]$`)
	// return reg.MatchString(OpportunitiesStage)
	return true
}

// 商机状态
func (Regex *RegexCheck) OptStatusCheck(StatusDescription string) bool {
	// reg := regexp.MustCompile(`^[0-1]$`)
	// return reg.MatchString(StatusDescription)
	return true
}

///////////////////////////整体信息判断////////////////////////

// 测试渠道信息
func CheckAgentInfo(agent *Agent) bool {
	if !Regex.PhoneCheck(agent.Mobile) {
		return false
	}

	if !Regex.EmailCheck(agent.Mail) {
		return false
	}

	return true
}

// 测试客户信息
func CheckCustomerInfo(customer *Customer) bool {
	//	if !Regex.PhoneCheck(customer.Phone) {
	//		return false
	//	}

	//	if !Regex.EmailCheck(customer.Mail) {
	//		return false
	//	}

	//	if !Regex.RTXNumberCheck(strconv.Itoa(customer.RTXNum)) {
	//		return false
	//	}

	return true
}

// 测试销售信息
func CheckAccountInfo(account *Account) bool {
	if !Regex.PhoneCheck(account.Mobile) {
		return false
	}

	return true
}

// 测试商机信息
func CheckOpportunities(biz *Biz) bool {

	return true
}

// 获取偏移
func GetOffsetByPage(page int) (offset int) {
	return ((page - 1) * 5)
}

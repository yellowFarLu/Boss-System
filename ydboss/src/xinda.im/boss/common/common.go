package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
//	"strings"

	"github.com/astaxie/beego"
)

//商机阶段map
var ops = &map[int]interface{}{
	1: "初步交流（10%）",
	2: "需求沟通（30%）",
	3: "商务沟通（50%）",
	4: "签约交款（100%）",
	0: "商务失败（0%）",
}

//商机类别map
var opc = &map[int]interface{}{
	1: "有度销售",
	2: "RTX销售",
	3: "RTX短信",
	4: "企业微信项目",
	0: "有度项目",
}

//商机
const (
	OPPORTUNITIES_PRE_COMMUNICATION = 0
)

// 没分配的客户使用的渠道ID
const (
	NoAllAgentId = 1 //没分配的agentid
)

const (
	SYSTEMSUCCESS = 0
	SYSTEMERROR   = -1
)

//角色
const (
	ROLE_ADMIM     = -1
	ROLE_XIN_ADMIN = 0
	ROLE_AGENT     = 1
	ROLE_EMPLOYEE  = 2
)

// 删除结果
const (
	DELOK   = 0
	DELFAIL = 1
)

const (
	PAGESIZE        = 50 //一页多少条记录
	COMMENTPAGESIZE = 20
	PAGE            = 1
)

// 单点登录
const (
	KStatusOK = iota
	KThridAuthFaild
	KThridAuthError
	KUserNotExist
	KUserBan
)

// 账号及登录
const (
	LOGINSTATUSOK = iota
	LOGINSTATUSNAMENOTEXIST
	LOGINSTATUSPWDERROR
	ACCOUNTEXIST
)

//用户
const (
	CUSTOMER_EXISTS     = 1
	CUSTOMER_NOT_EXISTS = 0
)

// 筛选参数
type FilterParam struct {
	G         map[string][]string
	SortKey   string
	SortIndex int
	Page      int
	Tags      []string
	Provinces []string
}

type Account struct {
	AccountId string `json:"employee_id"`
	Pwd       string `json:"pwd"`
	Name      string `json:"name"`
	Gender    int    `json:"sex"`
	Mobile    string `json:"phone"`
	Timex     string `json:"timex"`
	AgentId   int    `json:"agent_id"`
	RoleId    int    `json:"role"`
}

// 备注
type Comment struct {
	Name       string `json:"name"`
	CommentId  int    `json:"comment_id"`
	CustomerId int    `json:"customer_id"`
	Committer  string `json:"committer"`
	Comments   string `json:"comments"`
	Timex      string `json:"timex"`
	Type       int    `json:"type"`
}

// 客户
type Customer struct {
	CustomerId     int     `json:"customer_id"`
	EntName        string  `json:"entName"`
	RTXNum         int     `json:"rxt_num"`
	Contacts       string  `json:"contacts"`
	Phone          string  `json:"phone"`
	Mobile         string  `json:"mobile"`
	QQ             string  `json:"qq"`
	Mail           string  `json:"email"`
	Timex          string  `json:"timex"`
	Assign_status  int     `json:"dis_state"` // 0未分配    1分配给渠道    2分配给销售
	LastFollowTime *string `json:"last_follow_time"`
	Note           string  `json:"remarks"`
	Agent          *Agent
	Account        *Account
	Province       *Province
	City           *City
	Tags           []*Tag
	Comments       []*Comment

	Staff            int `json:"staff"` // 组织架构数
	Active           int `json:"active"`
	EmpCount         int `json:"emp_count"` // 员工人数
	EntYesterdayInfo *EntYesterdayInfo
}

type CustomerInfo struct {
	Buin   int32  `json:"buin"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`

	SerVersion    string `json:"ser_version"`
	Osver         string `json:"osver"`
	Ostype        string `json:"ostype"`
	Addr          string `json:"addr"`
	AgVersion     string `json:"ag_version"`
	AgOsver       string `json:"ag_osver"`
	AgOstype      string `json:"ag_ostype"`
	AgAddr        string `json:"ag_addr"`
	RtxsvrVersion string `json:"rtxsvr_version"`

	ShortName     string `json:"name"`
	Timex         string `json:"registTime"`
	MachineKey    string `json:"machineKey"`
	AuthNum       int32  `json:"auth"`
	ExpireTime    string `json:"expireTime"`
	PushKey       string `json:"pushKey"`
	NeedPush      int32  `json:"needPush"`
	Active        int32  `json:"act_num"`
	Staff         int32  `json:"all_num"`
	OnlineNum     int32  `json:"online_num"`
	RtxV          string `json:"rtx_v"`
	SrvV          string `json:"srv_v"`
	LicenseDate   string `json:"licenseDate"`
	Version       int    `json:"version"`
	CustomVersion int    `json:"customVersion"`
}

type EntStatiscInfo struct {
	Buin      int32
	Staff     int32
	Online    int32
	Limits    int32
	Active    int32
	Block     int32
	Expire    int64
	Clientmsg int32
	Rtxmsg    int32
	Timex     int64
	AgLive    string
	AgOnline  bool
	AgVersion string
}

type EntYesterdayInfo struct {
	Buin    int32 `json:"buin"`
	Total   int   `json:"total"`
	IOS     int   `json:"iOS"`
	Android int   `json:"android"`
	Pc      int   `json:"pc"`
	Web     int   `json:"web"`
}

// 渠道
type Agent struct {
	AgentId  int    `json:"agent_id"`
	Name     string `json:"agent_name"`
	Contacts string `json:"contacts"`
	Mobile   string `json:"contacts_phone"`
	Mail     string `json:"contacts_mail"`
	Note     string `json:"note"`
	Timex    string `josn:"timex"`

	Manager *Account `json:"account"`
}

//商机
type Biz struct {
	BizId        int     `json:"opportunities_id"`
	AccountId    string  `json:"account_id"`
	Title        string  `json:"title"`
	Content      string  `json:"content"`
	Amount       float64 `json:"sales"`  // 销售量
	Status       int     `json:"status"` // 商机状态： 0初步交流 1需求沟通 2商务沟通 3签约交款 4商务失败
	Timex        string  `json:"beginTime"`
	EstimateTime *string `json:"estimate_time"` //预计签约时间
	RealTime     *string `json:"real_time"`     // 真实签约时
	Enabled      int     `json:"enabled"`       // 商机是否关闭

	Customer *Customer `json:"customer"`
	Tags     []*Tag
	Account  *Account `json:"account"`
}

// tag
type Tag struct {
	TagId int    `json:"tag_id"`
	Name  string `josn:"name"`
	Type  string `json:"type"`
	Note  string `json:"note"`
}

// 省份
type Province struct {
	Id   int    `json:"province_id"`
	Name string `json:"province_name"`
}

// 城市
type City struct {
	Id   int    `json:"city_id"`
	Name string `json:"city_name"`
}

// 生成回复对象
func NewResponse(status int) map[string]interface{} {
	rsp := make(map[string]interface{}, 0)
	rsp["code"] = status
	return rsp
}

// 生成分页对象
func NewPageInfo(page int, page_total int, record int) map[string]interface{} {
	pif := make(map[string]interface{}, 0)
	if page == 0 {
		page = 1
	}
	pif["page"] = page             // 当前多少页
	pif["page_size"] = PAGESIZE    //一页多少条
	pif["page_total"] = page_total //总共多少页
	pif["record"] = record         // 总共的记录
	return pif
}

// 获取分页数据
func PageCount(err error, rows *sql.Rows) (int, int) {
	if err != nil {
		beego.Info(err)
		os.Exit(1)
	}
	// 帖子总数
	var Count int
	for rows.Next() {
		rows.Scan(&Count)
	}
	var page_total = (Count / PAGESIZE)
	if Count%PAGESIZE > 0 {
		page_total++
	}

	rows.Close()

	return Count, page_total
}

// 获取分页数据
func PageInfoCount(page, total int) (to, page_total int) {
	to = (page - 1) * PAGESIZE
	if to > total {
		to = total
	}
	page_total = total / PAGESIZE
	if total%PAGESIZE > 0 {
		page_total++
	}
	return
}

func NewComments(account *Account, action string, m map[string]interface{}) (comment string) {
	comment = ""
	m["action"] = action
	_byte, _ := json.Marshal(m)
	comment = string(_byte)
	return
}

func getString(v interface{}) (rsp string) {
	rsp = ""
	switch v.(type) {
	case string:
		rsp = v.(string)
	case int:
		rsp = fmt.Sprintf("%d", v.(int))
	case int64:
		rsp = fmt.Sprintf("%d", int(v.(int64)))
	case *City:
		rsp = v.(*City).Name
	}
	return
}

func GetNullString(v sql.NullString) string {
	if !v.Valid {
		return ""
	} else {
		return v.String
	}
}

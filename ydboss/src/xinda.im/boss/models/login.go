package models

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"cindasoft.com/library/entity"
	"cindasoft.com/library/proto/protohttp"
	"cindasoft.com/library/slog"

	"github.com/astaxie/beego"
	. "xinda.im/boss/common"
)

const (
	kjgGetIndentifyTokenUrl = "/v3/api/jginfo/identify.gen"
	YDIdentifyAddr          = "https://youdu.im"
)

type LoginModel struct{}

// 判断是不是第一次登陆
func (this *LoginModel) FirstLogin(accountId string) (firstLogin bool) {
	var fLogin = 0
	query := `select f_login from t_account t1 
	where t1.account_id = ?;`
	err := Db.QueryRow(query, accountId).Scan(&fLogin)
	if err != nil {
		beego.Error(err)
	}

	if fLogin == 0 {
		return true
	}

	return false
}

func (this *LoginModel) Login(account *Account) (code int) {
	code = LOGINSTATUSOK
	account.Pwd = Sh1(account.Pwd)
	query := `select pwd from t_account t1  where account_id=? and enabled=1;`
	var pwd string
	err := Db.QueryRow(query, account.AccountId).Scan(&pwd)

	if err != nil && err != sql.ErrNoRows {
		code = SYSTEMERROR
		beego.Error(err)
		return
	}
	if err != nil && err == sql.ErrNoRows {
		code = LOGINSTATUSNAMENOTEXIST
		return
	}

	if strings.Compare(pwd, account.Pwd) == 0 {
		code = LOGINSTATUSOK
	} else {
		code = LOGINSTATUSPWDERROR
	}

	return
}

func (this *LoginModel) YDLogin(token string) (string, Response) {
	if token == "" {
		return "", NewYDResponse(KThridAuthFaild)
	}
	rsp, err := http.Get(YDIdentifyAddr + "/v3/api/jginfo/identify?token=" + token)

	if err != nil {
		slog.Warn("YD identify error:", err)
		return "", NewYDResponse(KThridAuthError)
	}

	if rsp.StatusCode != http.StatusOK {
		slog.Warn("YD identify error:", rsp.Status)
		return "", NewYDResponse(KThridAuthError)
	}

	var resp struct {
		Status   *protohttp.Status
		UserInfo *entity.UserInfo
	}

	bs, _ := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	if err = json.Unmarshal(bs, &resp); err != nil {
		slog.Warn("YD identify error:", string(bs), err)
		return "", NewYDResponse(KThridAuthError)
	}
	if resp.Status == nil || !resp.Status.IsStatusOK() || resp.UserInfo == nil {
		slog.Warn("YD identify error:", string(bs), err)
		return "", NewYDResponse(KThridAuthFaild)
	}

	var ban int
	err = Db.QueryRow("select enabled from t_account where account_id=?", resp.UserInfo.Email).Scan(&ban)
	if err == sql.ErrNoRows {
		return "", NewYDResponse(KUserNotExist)
	}
	if err != nil {
		slog.Exit(err)
	}
	if ban == 0 {
		return "", NewYDResponse(KUserBan)
	}

	return resp.UserInfo.Email, NewYDResponse(KStatusOK)
}

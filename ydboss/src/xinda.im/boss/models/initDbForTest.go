package models

import (
	"database/sql"

	"github.com/astaxie/beego"
)

func InitTest() {
	var err error
	if Db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/emoa_boss?charset=utf8"); err != nil {
		beego.Error("[open mysql err]: ", err)
		return
	}

	//判断数据库是否还处于有效连接状态
	err = Db.Ping()
	if err != nil {
		beego.Error("[mysql ping failed]: ", err)
		return
	}
	//连接池的最大空闲连接跟最大连接数
	Db.SetMaxIdleConns(1)
	Db.SetMaxOpenConns(10)
}

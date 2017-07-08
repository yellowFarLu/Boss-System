package models

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	//change work dir
	os.Chdir(filepath.Dir(os.Args[0]))
	beego.Info(os.Getwd())
	file, err := os.Open("config.json")
	if err != nil {
		beego.Error("open config.json err:", err)
		return
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		beego.Error("open config.json err:", err)
	}

	//读取json Unmarshal存到map里
	err = json.Unmarshal(bs, &config)
	if err != nil {
		beego.Error("[Json unmarshal error]: ", err)
		return
	}

	InitDB()
	InitLog()
	CacheInit()

	InitBackup()
}

func InitBackup() {
	backupTool := NewBackupTool()
	backupTool.load()
}

func InitDB() {
	var err error
	if Db, err = sql.Open("mysql", config.MainDB); err != nil {
		beego.Error("[open mysql emoa_boss err]: ", err)
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

func InitLog() {
	err := os.MkdirAll("./log/", os.ModePerm)
	if err != nil {
		beego.Error("makDir log error: ", err)
		return
	}
	f, err := os.OpenFile(config.LogName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		beego.Error("create log error: ", err)
		return
	}
	defer f.Close()
}

package models

import (
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/astaxie/beego"

	"cindasoft.com/library/utils"
)

type BackupTool struct{}

// 初始化工具类
func NewBackupTool() (backupTool *BackupTool) {
	backupTool = &BackupTool{}
	return
}

/*
shell> mysqladmin -h 'other_hostname' create db_name
（mysqladmin -h other_hostname -u username -p create db_name）
shell> mysqldump --opt db_name | mysql -h 'other_hostname' db_name
(mysqldump -h local_host -u username -p local_db_name | mysql -h other_hostname -u username -p db_name)
*/
// 备份方法
func (this *BackupTool) _load() {
	one := `mysqldump`
	cmdOne := exec.Command(one, `-uroot`, `-proot`, `emoa_boss`)

	rs, err := cmdOne.Output()
	if err != nil {
		beego.Error(err)
	}
	err = ioutil.WriteFile(`boss_backup.sql`, rs, 0666)

	if err != nil {
		beego.Error(`cmd write error   `, err)
	}

	// 备份到另外一台机器
	// beego.Info("开始备份")
	// one := `mysqladmin`
	// cmdOne := exec.Command(one, `-h`, `10.0.0.83`, `-uroot`, `-proot`, `create`, `boss_backup`)
	// err := cmdOne.Run()
	// if err != nil {
	// 	beego.Error(err)
	// }

	// two := `mysqldump`
	// cmdTwo := exec.Command(two, `-h`, `182.254.180.76`, `-uroot`, `-proot`, `boss_backup`)
	// err = cmdTwo.Run()
	// if err != nil {
	// 	beego.Error(err)
	// }

	// two := `mysql`
	// cmdTwo := exec.Command(two, `-h`, `10.0.0.83`, `-uroot`, `-proot`, `boss_backup`)
	// err := cmdTwo.Run()
	// if err != nil {
	// 	beego.Error(err)
	// }
}

// 定时备份
func (this *BackupTool) load() {
	this._load()

	utils.Go(func() {
		for {
			select {
			case <-time.After(24 * time.Hour):
				this._load()
			}
		}
	})
}

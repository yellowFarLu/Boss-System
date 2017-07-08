package models

import (
	"encoding/json"
	"sync"
	"time"

	ph "cindasoft.com/library/proto/protohttp"
	"cindasoft.com/library/utils"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"

	. "xinda.im/boss/common"
)

const (
	API_CUSTOMER_INFO_LIST = "/v3/rpc/jgopsd/customer.info.list"
)

type CustomerCache struct {
	sync.RWMutex
	api       *ph.Client
	ent_m     map[int32]*CustomerInfo
	statisc_m map[int32]*EntStatiscInfo
	yes_m     map[int32]*EntYesterdayInfo
	cusname_m map[string]*Customer // key是name
	cusbuin_m map[int]int          // key是总机号
}

func NewCustomerServer() *CustomerCache {
	a := &CustomerCache{
		ent_m:     make(map[int32]*CustomerInfo, 0),
		statisc_m: make(map[int32]*EntStatiscInfo, 0),
		yes_m:     make(map[int32]*EntYesterdayInfo, 0),
		cusname_m: make(map[string]*Customer, 0),
		cusbuin_m: make(map[int]int, 0),
		api:       &ph.Client{},
	}
	a.load()
	return a
}

func (this *CustomerCache) load() {
	this._load()
	utils.Go(func() {
		for {
			select {
			case <-time.After(5 * time.Minute):
				this._load()
			}
		}
	})
}

func (this *CustomerCache) _load() {
	req := ph.NewJsonRequest()
	token := "7FF1809CF33B4A1F9D21F9F396CB956B"
	req.AddField("token", token)
	cr := this.api.Post2(config.OpsdPath+API_CUSTOMER_INFO_LIST, req)
	st := cr.ReadStatus()
	if !st.IsStatusOK() {
		beego.Error("[read api error]:", API_CUSTOMER_INFO_LIST)
		return
	}
	ent_m := make(map[int32]*CustomerInfo, 0)
	statisc_m := make(map[int32]*EntStatiscInfo, 0)
	yes_m := make(map[int32]*EntYesterdayInfo, 0)
	cr.Read("ent_data", &ent_m)
	cr.Read("statisc", &statisc_m)
	cr.Read("yes_m", &yes_m)
	this.Lock()
	defer this.Unlock()

	this.ent_m = ent_m
	this.statisc_m = statisc_m
	this.yes_m = yes_m

	// 5分钟更新一次Redis
	rc, err := redis.Dial(REDIS_PRODOCOL, REDIS_ADDRESS)
	if err != nil {
		beego.Error(err)
		return
	}
	defer rc.Close()

	query := `select t1.name,t1.rtx_number from t_customer t1`
	rows, err := Db.Query(query)
	if err != nil {
		beego.Error(err)
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

		rows.Scan(&c.EntName, &c.RTXNum)
		// this.cusname_m[c.EntName] = c
		// this.cusbuin_m[c.RTXNum] = c.RTXNum

		//把用户数据缓存到Redis
		_, err := rc.Do("APPEND", c.RTXNum, c.RTXNum)
		if err != nil {
			beego.Error(err)
		}
		cb, err := json.Marshal(c)
		cjson := string(cb)
		if err != nil {
			beego.Error(err)
		}
		_, err = rc.Do("APPEND", c.EntName, cjson)
		if err != nil {
			beego.Error(err)
		}
	}

}

func (this *CustomerCache) GetCustomInfo(buin int) (c *CustomerInfo) {
	c = &CustomerInfo{}
	this.RLock()
	defer this.RUnlock()

	b := int32(buin)
	if _, exist := this.ent_m[b]; !exist {
		return
	}
	c = this.ent_m[b]
	return
}

func (this *CustomerCache) GetCustomStatiscInfo(buin int) (c *EntStatiscInfo) {
	c = &EntStatiscInfo{}
	this.RLock()
	defer this.RUnlock()

	b := int32(buin)
	if _, exist := this.statisc_m[b]; !exist {
		return
	}
	c = this.statisc_m[b]
	return
}

func (this *CustomerCache) GetCustomYesInfo(buin int) (c *EntYesterdayInfo) {
	c = &EntYesterdayInfo{}
	this.RLock()
	defer this.RUnlock()

	b := int32(buin)
	if _, exist := this.yes_m[b]; !exist {
		return
	}
	c = this.yes_m[b]
	return
}

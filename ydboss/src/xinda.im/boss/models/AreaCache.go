package models

import (
	"strings"
	"sync"

	"cindasoft.com/library/utils"
	"github.com/astaxie/beego"
	. "xinda.im/boss/common"
)

type AreaCache struct {
	sync.RWMutex
	province []*Province
	m        map[int][]*City
}

func NewAreaServer() *AreaCache {
	a := &AreaCache{
		province: make([]*Province, 0),
		m:        make(map[int][]*City, 0),
	}
	a.load()
	return a
}

func (this *AreaCache) load() {
	this._load()
	//	utils.Go(func() {
	//		//		for t := range time.Tick(5 * time.Minute) {
	//		//			this._load()
	//		//		}
	//		for {
	//			select {
	//			case <-time.After(5 * time.Minute):
	//				this._load()
	//			}
	//		}
	//	})
}

func (this *AreaCache) _load() {
	provinces := this.GetProvinces()
	this.IntoCache(provinces)
}

func (this *AreaCache) GetProvinces() (p []*Province) {
	p = make([]*Province, 0)
	rows, err := Db.Query(`select province_id,province from t_province order by display_order desc;`)
	defer rows.Close()
	if err != nil {
		beego.Error(err)
		return
	}

	this.Lock()
	defer this.Unlock()
	for rows.Next() {
		province := &Province{}
		err := rows.Scan(&province.Id, &province.Name)
		if err != nil {
			beego.Error(err)
			break
		}
		p = append(p, province)
		this.province = append(this.province, province)
	}
	return
}

func (this *AreaCache) IntoCache(provinces []*Province) {
	query := `select city_id,city from t_city where province_id = ? order by display_order desc;`
	for _, v := range provinces {
		value := v
		// utils.Go(func() {}异步读取
		utils.Go(func() {
			rows, err := Db.Query(query, value.Id)
			defer rows.Close()
			if err != nil {
				beego.Error(err)
				return
			}
			this.Lock()
			defer this.Unlock()
			citys := make([]*City, 0)
			for rows.Next() {
				city := &City{}
				err := rows.Scan(&city.Id, &city.Name)
				if err != nil {
					beego.Error(err)
					break
				}
				citys = append(citys, city)
			}
			this.m[value.Id] = citys
		})
	}
}

func (this *AreaCache) GetAllProvince() (list []*Province) {
	this.RLock()
	list = this.province
	this.RUnlock()
	return
}

func (this *AreaCache) GetProvinceByKey(key string) (list []*Province) {
	this.RLock()
	list = make([]*Province, 0)
	for _, v := range this.province {
		if ok := strings.Contains(v.Name, key); ok {
			list = append(list, v)
		}
	}
	this.RUnlock()
	if len(list) > 5 {
		list = list[0:5]
	}
	return
}

func (this *AreaCache) GetCity(id int) (list []*City) {
	this.RLock()
	list = this.m[id]
	this.RUnlock()
	return
}

func (this *AreaCache) GetCityByCityId(cityId int) (cname string) {
	for _, v := range this.m {
		for _, _v := range v {
			if _v.Id == cityId {
				cname = _v.Name
				return
			}
		}
	}

	return "未知城市"
}

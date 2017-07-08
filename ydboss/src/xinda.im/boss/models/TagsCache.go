package models

import (
	"strings"
	"sync"
	"time"

	"cindasoft.com/library/utils"
	"github.com/astaxie/beego"

	. "xinda.im/boss/common"
)

type TagsCache struct {
	sync.RWMutex
	tags  []*Tag
	m     map[int]*Tag
	biz_m map[int]*Tag // key是商机Id
}

func NewTagsServer() *TagsCache {
	a := &TagsCache{
		tags:  make([]*Tag, 0),
		m:     make(map[int]*Tag, 0),
		biz_m: make(map[int]*Tag, 0),
	}
	a.load()
	return a
}

func (this *TagsCache) load() {
	this._load()
	utils.Go(func() {
		for {
			select {
			case <-time.After(1 * time.Hour):
				this._load()
			}
		}
	})
}

func (this *TagsCache) _load() {
	query := `select tag_id, name, type, note from t_tag order by tag_id desc;`
	rows, err := Db.Query(query)
	defer rows.Close()
	if err != nil {
		beego.Error(err)
		return
	}

	this.Lock()
	tags := make([]*Tag, 0)
	for rows.Next() {
		tag := &Tag{}
		err = rows.Scan(&tag.TagId, &tag.Name, &tag.Type, &tag.Note)
		if err != nil {
			beego.Error(err)
			return
		}
		tags = append(tags, tag)
		this.m[tag.TagId] = tag
	}
	this.tags = tags
	this.Unlock()
}

func (this *TagsCache) GetAllTags() (list []*Tag) {
	this.RLock()
	list = this.tags
	this.RUnlock()
	return
}

func (this TagsCache) GetTagList(ids []int) (list []*Tag) {
	this.RLock()
	list = make([]*Tag, 0)
	for _, v := range ids {
		list = append(list, this.m[v])
	}
	this.RUnlock()
	return
}

func (this *TagsCache) GetTagByKey(key string) (list []*Tag) {
	this.RLock()
	list = make([]*Tag, 0)
	for _, v := range this.tags {
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

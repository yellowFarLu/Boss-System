package models

import (
	"testing"

	"fmt"

	. "xinda.im/boss/common"
)

// func Test_FilterAllCustomer(t *testing.T) {
// 	InitTest()
// 	tagsCache = NewTagsServer()
// 	cusModel := &CustomerModel{}
// 	keyStrArr := make([]string, 0)
// 	keyStrArr = append(keyStrArr, `{"rtx_number": "'32423421'"}`)
// 	// keyStrArr = append(keyStrArr, `{"last_follow_time_begin": "'%w%'"}`)
// 	keyStrArr = append(keyStrArr, `{"timex_begin": "'2017-04-06 15:42:09'"}`)
// 	keyStrArr = append(keyStrArr, `{"timex_end": "'2017-04-06 15:42:09'"}`)
// 	sortKey := "timex"
// 	slice := cusModel.FilterAllCustomer(keyStrArr, sortKey, 1)

// 	for _, v := range slice {
// 		fmt.Println("获取的数据：", v.EntName)
// 	}
// }

// func Test_FilterAgentCustomer(t *testing.T) {
// 	InitTest()
// 	tagsCache = NewTagsServer()
// 	cusModel := &CustomerModel{}
// 	keyStrArr := make([]string, 0)
// 	keyStrArr = append(keyStrArr, `{"rtx_number": "'32423421'"}`)
// 	keyStrArr = append(keyStrArr, `{"timex_begin": "'2017-04-06 15:42:09'"}`)
// 	keyStrArr = append(keyStrArr, `{"timex_end": "'2017-04-06 15:42:09'"}`)
// 	sortKey := "timex"
// 	account := &Account{AgentId: 2}
// 	slice := cusModel.FilterAgentCustomer(account, keyStrArr, sortKey, 1)

// 	for _, v := range slice {
// 		fmt.Println("获取的数据：", v.EntName)
// 	}
// }

func Test_FilterEmpCustomer(t *testing.T) {
	InitTest()
	tagsCache = NewTagsServer()
	cusModel := &CustomerModel{}

	account := &Account{AccountId: "peyton.li@xinda.im"}

	arr1 := make([]string, 0)
	arr1 = append(arr1, "2017-05-06 15:42:09")
	g := map[string][]string{
		"follow_time_to": arr1,
	}
	sortKey := "timex"
	page := 1
	sortIndex := 0
	tags := g["tags[]"]
	provinces := g["provinces[]"]

	filter := &FilterParam{
		G:         g,
		SortKey:   sortKey,
		SortIndex: sortIndex,
		Page:      page,
		Tags:      tags,
		Provinces: provinces,
	}

	slice := cusModel.FilterEmpCustomer(account, filter)

	for _, v := range slice {
		fmt.Println("获取的数据：", v.EntName)
	}
}

// func Test_AdminAddCustomer(t *testing.T) {
// 	InitTest()
// 	cusModel := &CustomerModel{}
// 	var agent = &Agent{AgentId: NoAllAgentId}
// 	var city = &City{Id: 1}
// 	customer := &Customer{
// 		EntName:       "同仁",
// 		RTXNum:        87878787,
// 		Contacts:      "老北京",
// 		Mobile:        "29876781",
// 		Phone:         "13527206719",
// 		QQ:            "296947440",
// 		Mail:          "296947440@qq.com",
// 		Agent:         agent,
// 		City:          city,
// 		Timex:         "2016-13-04 09:09:09",
// 		Assign_status: 0,
// 		Note:          "note",
// 	}

// 	account := &Account{}
// 	account.AccountId = "admin"
// 	cusModel.AdminAddCustomer(customer, account, nil)
// }

// func Test_AgentAddCustomer(t *testing.T) {
// 	InitTest()
// 	cusModel := &CustomerModel{}
// 	var agent = &Agent{AgentId: NoAllAgentId}
// 	var city = &City{Id: 1}
// 	customer := &Customer{
// 		EntName:       "发的",
// 		RTXNum:        87878787,
// 		Contacts:      "官方发",
// 		Mobile:        "29876781",
// 		Phone:         "13527206719",
// 		QQ:            "296947440",
// 		Mail:          "296947440@qq.com",
// 		Agent:         agent,
// 		City:          city,
// 		Timex:         "2016-13-04 09:09:09",
// 		Assign_status: 0,
// 		Note:          "note",
// 	}

// 	account := &Account{}
// 	account.AccountId = "admin"
// 	ok := cusModel.AgentAddCustomer(customer, account, nil)
// 	fmt.Println(ok)
// }

// func Test_EmpAddCustomer(t *testing.T) {
// 	InitTest()
// 	cusModel := &CustomerModel{}
// 	var agent = &Agent{AgentId: NoAllAgentId}
// 	var city = &City{Id: 1}
// 	customer := &Customer{
// 		EntName:       "销售增加客户",
// 		RTXNum:        87878787,
// 		Contacts:      "销售增加客户",
// 		Mobile:        "29876781",
// 		Phone:         "13527206719",
// 		QQ:            "296947440",
// 		Mail:          "296947440@qq.com",
// 		Agent:         agent,
// 		City:          city,
// 		Timex:         "2016-13-04 09:09:09",
// 		Assign_status: 0,
// 		Note:          "note",
// 	}

// 	account := &Account{}
// 	account.AccountId = "admin"
// 	ok := cusModel.AgentAddCustomer(customer, account, nil)
// 	fmt.Println(ok)
// }

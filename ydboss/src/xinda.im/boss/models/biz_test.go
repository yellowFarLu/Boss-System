package models

// import (
// 	"testing"

// 	. "xinda.im/boss/common"
// )

// func Test_InsertBiz(t *testing.T) {
// 	InitTest()
// 	cus := &Customer{CustomerId: 2}
// 	biz := &Biz{AccountId: "admin", Customer: cus, Title: "dada",
// 		Content: "ct", Amount: 12.30, Status: 0, Timex: "2017-03-31 09:09:09",
// 		EstimateTime: "2017-03-31 09:09:09", RealTime: "2017-03-31 09:09:09"}
// 	bizModel := &BizModel{}
// 	bizModel.InsertBiz(biz)
// }

// func Test_AlertBiz(t *testing.T) {
// 	InitTest()
// 	cus := &Customer{CustomerId: 2}
// 	biz := &Biz{AccountId: "admin", Customer: cus, Title: "dada",
// 		Content: "ct阿达的", Amount: 12.30, Status: 0, Timex: "2017-03-31 09:09:09",
// 		EstimateTime: "2017-03-31 09:09:09", RealTime: "2017-03-31 09:09:09"}
// 	bizModel := &BizModel{}
// 	bizModel.AlertBiz(biz)
// }

// func Test_FilterBizList(t *testing.T) {
// 	InitTest()
// 	bizModel := &BizModel{}
// 	keyStrArr := make([]string, 0)
// 	keyStrArr = append(keyStrArr, `{"title": "'%t%'"}`)
// 	keyStrArr = append(keyStrArr, `{"status": "1"}`)
// 	keyStrArr = append(keyStrArr, `{"timex_begin": "'2017-04-06 15:42:34'"}`)
// 	keyStrArr = append(keyStrArr, `{"timex_end": "'2017-04-06 15:42:34'"}`)
// 	sortKey := "timex"
// 	bizModel.FilterBizList(keyStrArr, sortKey, 1)
// }

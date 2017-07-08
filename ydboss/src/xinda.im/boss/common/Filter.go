package common

import (
	"github.com/astaxie/beego/context"
)

var FilterCoding = func(ctx *context.Context) {
	// coding := ctx.Request.Header.Get("charset")
	// coding = strings.ToLower(coding)
	// content := ctx.Input.RequestBody
	// if coding != "utf8" {

	// 	ctx.Input.RequestBody = data
	// }
}

var FilterUser = func(ctx *context.Context) {
	// if ctx.Request.RequestURI == "/login" || ctx.Request.RequestURI == "/login.html" {
	// 	return
	// }

	// account := ctx.Input.Session("userInfo")
	// if account == nil {
	// 	ctx.Redirect(302, "/login.html")
	// 	return
	// }
}

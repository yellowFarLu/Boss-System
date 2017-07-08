package models

import "xinda.im/boss/common"

type Response struct {
	Code int
	Rsp  int
}

func NewYDResponse(code int) (rsp Response) {
	rsp = Response{Code: code}
	return
}

func (this *Response) IsOK() (ok bool) {
	if this.Code == common.KStatusOK {
		ok = true
		return
	}

	return
}

package response

import (
	"fmt"
	// http client driver
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"pangoo/pkg/handler"
)

type BaseResp struct {
	Code int8   `thrift:"code,1" form:"code" json:"code" query:"code"`
	Msg  string `thrift:"msg,2" form:"msg" json:"msg" query:"msg"`
	Data interface{} `thrift:"data,3" form:"data" json:"data" query:"data"`
}

func NewBaseResp() *BaseResp {
	return &BaseResp{}
}

func (p *BaseResp) GetCode() (v int8) {
	return p.Code
}

func (p *BaseResp) GetMsg() (v string) {
	return p.Msg
}

func (p *BaseResp) GetData() (v interface{}) {
	return p.Data
}

func (p *BaseResp) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaseResp(%+v)", *p)
}

type Response interface {
	Succ(body interface{}) Response
	Fail(err error) Response
}

type DefaultResponse struct {
	Msg handler.ExceptionHandler `json:"msg"`
	Data   interface{}             `json:"data"`
}

func (dr *DefaultResponse) Succ(body interface{}) Response {
	success, _ := handler.NewSuccess("成功").(*handler.ExceptionHandler)
	dr.Msg = *success
	dr.Data = body
	return dr
}

func (dr *DefaultResponse) Fail(err error) Response {
	errCustom, ok := err.(*handler.ExceptionHandler)
	if !ok {
		hlog.Error("error type invalid, please use custom error")
		unkownError, _ := handler.NewUnkownError(err.Error()).(*handler.ExceptionHandler)
		dr.Msg = *unkownError
	} else {
		dr.Msg = *errCustom
	}

	return dr
}

package res

import (
	// http client driver
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"primus/pkg/handler"
)

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

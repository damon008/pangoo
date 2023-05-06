package handler


import (
	"fmt"
	"primus/pkg/constants"
)


const (
	badRequestType = "[bad request]"
	innerErrorType = "[inner error]"
	successType    = "[success]"
)

type ExceptionHandler struct {
	Code      int    `json:"code"`
	ErrorType string `json:"-"`
	Msg       string `json:"msg"`
}

func (ec ExceptionHandler) Error() string {
	return ec.Msg
}

func NewBadReqError(code int, format string, paras ...interface{}) error {
	return newError(code, badRequestType, format, paras...)
}

func NewInnerError(code int, format string, paras ...interface{}) error {
	return newError(code, innerErrorType, format, paras...)
}

func NewSuccess(format string, paras ...interface{}) error {
	return newSuccess(constants.OPERATE_SUCCESS, successType, format, paras...)
}

func NewUnkownError(format string, paras ...interface{}) error {
	return newError(constants.UNKNOWN_ERROR, innerErrorType, format, paras...)
}

func IsInnerError(err error) bool {
	if errCustom, ok := err.(ExceptionHandler); ok {
		return errCustom.ErrorType == innerErrorType
	}
	return false
}

func IsBadRequestError(err error) bool {
	if errCustom, ok := err.(ExceptionHandler); ok {
		return errCustom.ErrorType == badRequestType
	}
	return false
}

func IsSuccess(err error) bool {
	if err == nil {
		return true
	}

	if errCustom, ok := err.(*ExceptionHandler); ok {
		return errCustom.ErrorType == successType
	}
	return false
}

func NewError(code int, format string, paras ...interface{}) *ExceptionHandler {
	return &ExceptionHandler{Code: code, Msg: fmt.Sprintf(format, paras...)}
}

func newSuccess(code int, errorType string, format string, paras ...interface{}) error {
	return &ExceptionHandler{Code: code, ErrorType: errorType, Msg: fmt.Sprintf(format, paras...)}
}

func newError(code int, errorType string, format string, paras ...interface{}) error {
	return &ExceptionHandler{Code: code, ErrorType: errorType, Msg: fmt.Sprintf(format, paras...)}
}
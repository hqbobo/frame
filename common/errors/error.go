package errors

import "github.com/micro/go-micro/errors"

// SvcError 统一错误封装
type SvcError errors.Error

var gNAME string

// SetID 设置名字
func SetID(name string) { gNAME = name }

// ParseAppError 反向解析错误
func ParseAppError(err error) *SvcError {
	if err == nil {
		return nil
	}
	oe := errors.Parse(err.Error())
	e := new(SvcError)
	e.Code = oe.Code
	e.Detail = oe.Detail
	e.Id = oe.Id
	return e
}

// NewAppError 新的错误
func NewAppError(code int32) error {
	detail := defs.getmsg(code)
	return errors.New(gNAME, detail, code)
}

// NewCustomError 新的错误
func NewCustomError(code int32, msg string) error {
	return errors.New(gNAME, msg, code)
}

// Error 错误描述
func (e *SvcError) Error() string {
	if e == nil {
		return ""
	}
	return e.Detail
}

// GetCode 错误Code
func (e *SvcError) GetCode() int32 {
	if e == nil {
		return 0
	}
	return e.Code
}

package restful

import (
	"strings"
)

const (
	CodeSuccess     = 200
	CodeClientError = 400
	CodeServerError = 500
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (r *Response) WithMessage(v string) *Response {
	r.Msg = v
	return r
}

func (r *Response) WithData(v any) *Response {
	r.Data = v
	return r
}

func Success(data any) *Response {
	return &Response{
		Code: CodeSuccess,
		Data: data,
	}
}

func Error(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}

func ParamError(msg ...string) *Response {
	r := Error(CodeClientError, "参数错误")
	if len(msg) > 0 {
		r.WithMessage(strings.Join(msg, "; "))
	}
	return r
}

type TypedResponse[T any] struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data T      `json:"data,omitempty"`
}

package pkgres

import (
	serviceError "cms-server/domain/service/error"
	serviceRes "cms-server/domain/service/response"
	"net/http"
)

type ErrorApp struct {
	Message string `json:",omitempty"`
	Data    any    `json:",omitempty"`
	code    int
}

func NewErr(msg string) serviceError.ErrorApp {
	return &ErrorApp{
		Message: msg,
		code:    http.StatusInternalServerError,
	}
}

func Err(err error) serviceError.ErrorApp {
	return &ErrorApp{
		Message: err.Error(),
		code:    http.StatusInternalServerError,
	}
}

func (r *ErrorApp) Error() string {
	return r.Message
}

func (r *ErrorApp) SetMessage(message string) serviceError.ErrorApp {
	r.Message = message
	return r
}

func (r *ErrorApp) SetData(data any) serviceError.ErrorApp {
	r.Data = data
	return r
}

func (r *ErrorApp) Code(code int) serviceError.ErrorApp {
	r.code = code
	return r
}

func (r *ErrorApp) GetCode() int {
	return r.code
}

func (r *ErrorApp) BadReq() serviceError.ErrorApp {
	return r.Code(http.StatusBadRequest)
}

func (r *ErrorApp) UnprocessableEntity() serviceError.ErrorApp {
	return r.Code(http.StatusUnprocessableEntity)
}

func (r *ErrorApp) InternalServerError() serviceError.ErrorApp {
	return r.Code(http.StatusInternalServerError)
}

func (r *ErrorApp) NotFound() serviceError.ErrorApp {
	return r.Code(http.StatusNotFound)
}

func (r *ErrorApp) Unauthorized() serviceError.ErrorApp {
	return r.Code(http.StatusUnauthorized)
}

func (r *ErrorApp) Forbidden() serviceError.ErrorApp {
	return r.Code(http.StatusForbidden)
}

func (r *ErrorApp) Conflict() serviceError.ErrorApp {
	return r.Code(http.StatusConflict)
}

type res struct {
	Message string `json:",omitempty"`
	Data    any    `json:",omitempty"`
	code    int
}

func NewRes(msg string) serviceRes.Response {
	return &res{
		Message: msg,
	}
}

func ResData(data any) serviceRes.Response {
	return &res{
		Data: data,
	}
}

func (r *res) New(msg string) serviceRes.Response {
	return &res{
		Message: msg,
		code:    http.StatusOK,
	}
}

func (r *res) SetMessage(message string) serviceRes.Response {
	r.Message = message
	return r
}

func (r *res) SetData(data any) serviceRes.Response {
	r.Data = data
	return r
}

func (r *res) Code(code int) serviceRes.Response {
	r.code = code
	return r
}

func (r *res) GetCode() int {
	return r.code
}

package puberr

import (
	"context"
	"database/sql"
	"errors"

	"net/http"
)

/*
цель PubErr - упростить обработку ошибок в приложении.
плюсы:
	- отделить текст ошибки который видит пользователь от текста ошибки который видит разработчик
	- можно задать HTTP-код ответа
	- можно задать код ошибки
*/

var (
	ErrGeneric          = NewPubErr("error occurred").SetCode(1)
	ErrNotImplemented   = NewPubErr("not implemented").SetCode(10)
	ErrNotFound         = NewPubErr("not found").SetCode(11).SetHTTPCode(http.StatusNotFound)
	ErrExists           = NewPubErr("already exists").SetCode(12)
	ErrNoElements       = NewPubErr("no more elements").SetCode(13)
	ErrNotAuthorized    = NewPubErr("not authorized").SetCode(14).SetHTTPCode(http.StatusUnauthorized)
	ErrForbidden        = NewPubErr("forbidden").SetCode(15).SetHTTPCode(http.StatusForbidden)
	ErrResourceNotFound = NewPubErr("resource not found").SetCode(16).SetHTTPCode(http.StatusNotFound)
	ErrInvalidParams    = NewPubErr("invalid params").SetCode(17)
	ErrInternal         = NewPubErr("internal error").SetCode(18).SetHTTPCode(http.StatusInternalServerError)
	ErrNotOwnedResource = NewPubErr("this resource is not owned by this author").SetCode(19)
	ErrInvalidRequest   = NewPubErr("request format is not valid").SetCode(20)
	ErrInvalidToken     = NewPubErr("invalid token").SetCode(21)
)

var CodeToErr = map[int]PubErr{
	1:  ErrGeneric,
	10: ErrNotImplemented,
	11: ErrNotFound,
	12: ErrExists,
	13: ErrNoElements,
	14: ErrNotAuthorized,
	15: ErrForbidden,
	16: ErrResourceNotFound,
	17: ErrInvalidParams,
	18: ErrInternal,
	19: ErrNotOwnedResource,
	20: ErrInvalidRequest,
	21: ErrInvalidToken,
}

type PubErr struct {
	Cause     error           `json:"-"`
	PublicMsg string          `json:"error,omitempty"`
	ErrCode   int             `json:"errCode,omitempty"`
	HTTPCode  int             `json:"-"`
	Ctx       context.Context `json:"-"`
}

func NewPubErr(msg string) PubErr {
	return PubErr{
		PublicMsg: msg,
		HTTPCode:  http.StatusBadRequest,
	}
}

func (it PubErr) SetHTTPCode(code int) PubErr {
	it.HTTPCode = code
	return it
}

func (it PubErr) SetCode(code int) PubErr {
	it.ErrCode = code
	return it
}

func (it PubErr) SetCause(cause error) PubErr {
	it.Cause = cause
	return it
}

func (it PubErr) SetContext(ctx context.Context) PubErr {
	it.Ctx = ctx
	return it
}

func (it PubErr) Error() string {
	if it.Cause != nil {
		return it.Cause.Error()
	}

	return it.PublicMsg
}

func (it PubErr) Unwrap() error {
	return it.Cause
}

func ErrToPubErr(err error) (PubErr, error) {
	var appErr PubErr
	ok := errors.As(err, &appErr)
	if ok {
		return appErr, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound, nil
	}

	return NewPubErr("generic error").
			SetHTTPCode(http.StatusInternalServerError).
			SetCause(err),
		err
}

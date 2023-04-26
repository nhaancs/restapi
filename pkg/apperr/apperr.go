package apperr

import (
	"errors"
	"net/http"
	"restapi/pkg/appconst"
)

type appError struct {
	code   string
	msg    string
	status int
	err    error
}

func (e appError) Status() int {
	return e.status
}

func (e appError) Code() string {
	return e.code
}

func (e appError) Message() string {
	return e.msg
}

func (e appError) rootError() error {
	if err, ok := e.err.(appError); ok {
		return err.rootError()
	}

	return e.err
}

func (e appError) Error() string {
	if rootErr := e.rootError(); rootErr != nil {
		return e.rootError().Error()
	}
	return e.msg
}

func Wrap(err error, code, msg string, status int) error {
	return appError{
		code:   code,
		msg:    msg,
		status: status,
		err:    err,
	}
}

func Convert(err error) appError {
	if e := new(appError); errors.As(err, e) {
		return *e
	}

	return appError{
		code:   appconst.CodeInternal,
		msg:    "internal error",
		status: http.StatusInternalServerError,
		err:    err,
	}
}

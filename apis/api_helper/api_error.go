/*
 *
 * api_error.go
 * api_helper
 *
 * Created by lintao on 2019-01-26 16:58
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package api_helper

import "errors"

type Error struct {
	Err  error
	Code ResponetsCode
}

func (this *Error) Error() string {
	return this.Err.Error()
}

func NewError(err error, code ResponetsCode) *Error {
	return &Error{
		Err:  err,
		Code: code,
	}
}

func NewParamError(err error) *Error {
	return &Error{
		Err:  err,
		Code: StatusParamErr,
	}
}

func NewDBError(err error) *Error {
	return &Error{
		Err:  err,
		Code: StatusDBErr,
	}
}

func NewMsgError(str string) *Error {
	return &Error{
		Err:  errors.New(str),
		Code: StatusAlerMsg,
	}
}

func NewReloginError(err error) *Error {
	return &Error{
		Err:  err,
		Code: StatusRelogin,
	}
}

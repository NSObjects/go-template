/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package tools

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

/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

import "errors"

type Error struct {
	Err  error
	Code StatusCode
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewError(err error, code StatusCode) *Error {
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
		Code: StatusAlterMsg,
	}
}

func NewAuthError(err error) *Error {
	return &Error{
		Err:  err,
		Code: StatusAuth,
	}
}

/*
 *
 * middleware.go
 * api_helper
 *
 * Created by lintao on 2019-01-29 09:18
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package api_helper

import (
	"gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	Validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

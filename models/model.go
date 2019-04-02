/*
 *
 * model.go
 * models
 *
 * Created by lintao on 2019-01-10 17:07
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package models

import (
	"fmt"
	"time"
)

type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		panic(err)
	}

	now, err := time.ParseInLocation(timeFormart, string(data), local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).Unix())), nil
}

type ApiQuery struct {
	Page  int `json:"page" form:"page" query:"page"`
	Count int `json:"count" form:"count" query:"count"`
}

func (q *ApiQuery) Limit() (count, offset int) {
	if q.Page > 0 {
		q.Page -= 1
	}

	if q.Count == 0 {
		q.Count = 20
	}
	return q.Count, q.Page * q.Count
}

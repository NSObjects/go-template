/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package model

import (
	"fmt"
	"time"
)

type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		panic(err)
	}

	now, err := time.ParseInLocation(timeFormat, string(data), local)
	*t = Time(now)
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(*t).Unix())), nil
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

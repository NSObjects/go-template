/*
 * Created by lintao on 2023/7/26 下午5:08
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

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

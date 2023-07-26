/*
 * Created by lintao on 2023/7/26 下午2:40
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package param

type APIQuery struct {
	Page  int `json:"page" form:"page" query:"page"`
	Count int `json:"count" form:"count" query:"count"`
}

func (q *APIQuery) Limit() int {
	if q.Count == 0 {
		q.Count = 10
	}

	return q.Count
}

func (q *APIQuery) Offset() int {
	if q.Page > 0 {
		q.Page -= 1
	}

	return q.Page * q.Count
}

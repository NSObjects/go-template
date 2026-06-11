/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

import (
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/samber/do/v2"
)

// Register 注册数据层依赖。
func Register(i do.Injector) {
	do.Provide[biz.UserRepository](i, func(i do.Injector) (biz.UserRepository, error) {
		dm, err := do.Invoke[*db.DataManager](i)
		if err != nil {
			return nil, err
		}
		return NewUserRepository(dm), nil
	})
}

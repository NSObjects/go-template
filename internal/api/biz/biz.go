/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package biz

import (
	"github.com/samber/do/v2"
)

// Register 注册业务层依赖。
func Register(i do.Injector) {
	do.Provide[UserUseCase](i, func(i do.Injector) (UserUseCase, error) {
		repo, err := do.Invoke[UserRepository](i)
		if err != nil {
			return nil, err
		}
		return NewUserHandler(repo), nil
	})
}

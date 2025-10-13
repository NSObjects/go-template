/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package biz

import (
	"go.uber.org/fx"
)

var Model = fx.Options(fx.Provide(NewUserHandler))

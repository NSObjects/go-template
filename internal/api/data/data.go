/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package data

import (
	"github.com/NSObjects/go-template/internal/api/data/db"
	"go.uber.org/fx"
)

var Model = fx.Options(
	fx.Provide(NewUserDataSource, db.NewDataSource),
)

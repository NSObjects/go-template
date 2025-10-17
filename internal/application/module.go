package application

import (
	appuser "github.com/NSObjects/go-template/internal/application/user"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(appuser.NewService),
)

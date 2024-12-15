package id

import "go.uber.org/fx"

var Module = fx.Module(
	"id",
	fx.Provide(NewGenerator),
)

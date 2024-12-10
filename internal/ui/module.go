package ui

import "go.uber.org/fx"

var Module = fx.Module(
	"ui",
	fx.Provide(NewViews),
	fx.Invoke(Register),
)

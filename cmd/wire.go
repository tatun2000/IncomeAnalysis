//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"telegrammbot.core/internal/config"
	"telegrammbot.core/internal/domains/oauth"
	"telegrammbot.core/internal/domains/sheet"
	"telegrammbot.core/internal/domains/telegram"
)

func InjectAppGod(ctx context.Context, cfgPath string) (*app, func(), error) {
	wire.Build(
		config.NewConfig,
		wire.FieldsOf(new(config.Config),
			"GeneralOpts",
			"BotOpts",
			"SheetOpts",
			"OauthOpts",
		),

		newTracer,

		OauthSet,
		SheetSet,
		TelegramSet,

		wire.Struct(new(app), "*"),
	)
	return &app{}, nil, nil
}

var (
	OauthSet = wire.NewSet(
		wire.Bind(new(sheet.IOauthService), new(*oauth.Service)),
		oauth.NewService,
	)
	SheetSet = wire.NewSet(
		wire.Bind(new(telegram.ISheetService), new(*sheet.Service)),
		sheet.NewService,
	)
	TelegramSet = wire.NewSet(
		telegram.NewService,
	)
)

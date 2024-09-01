package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var Validator = validator.New()

type Config struct {
	GeneralOpts GeneralOpts `mapstructure:"generalOpts" validate:"required"`
	BotOpts     BotOpts     `mapstructure:"botOpts" validate:"required"`
	SheetOpts   SheetOpts   `mapstructure:"sheetOpts" validate:"required"`
	OauthOpts   OauthOpts   `mapstructure:"oauthOpts" validate:"required"`
}

type GeneralOpts struct {
	WorkDir string `mapstructure:"workDir" validate:"required,dir"`
}

type BotOpts struct {
	Token string `mapstructure:"endpoint" validate:"required"`
}

type SheetOpts struct {
	SpreadsheetId string `mapstructure:"spreadsheetId" validate:"required"`
}

type OauthOpts struct {
	RefreshToken string `mapstructure:"refreshToken" validate:"required"`
}

func NewConfig(configPath string) (cfg Config, err error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(configPath)
	err = v.ReadInConfig()
	if err != nil {
		return cfg, fmt.Errorf("NewConfig: %w", err)
	}
	if err = v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("NewConfig: %w", err)
	}
	if err = Validator.Struct(cfg); err != nil {
		return cfg, fmt.Errorf("NewConfig: %w", err)
	}
	return cfg, nil
}

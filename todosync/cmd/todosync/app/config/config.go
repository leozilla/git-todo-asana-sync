package config

import (
	"context"
	"fmt"
	validate "github.com/go-playground/validator/v10"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

type Config struct {
	GitRepoUrl string `env:"GIT_REPO_URL" validate:"hostname|hostname_rfc1123|ip"`
}

func MustLoadFromEnv(ctx context.Context, logger *zap.Logger) Config {
	var cfg Config

	l := envconfig.PrefixLookuper("APP_", envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, &cfg, l); err != nil {
		logger.Error("Service config could not be loaded from env variables")
		panic(err)
	}

	validator := validate.New()
	if err := validator.Struct(cfg); err != nil {
		logger.Error("Service config is invalid",
			zap.String("config", fmt.Sprintf("%+v", cfg)))
		panic(err)
	}

	logger.Info("Service config loaded",
		zap.String("config", fmt.Sprintf("%+v", cfg)))

	return cfg
}

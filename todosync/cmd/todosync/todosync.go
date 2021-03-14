package main

import (
	"context"
	"git-todo-asana-sync/todosync/cmd/todosync/app/config"
	"go.uber.org/zap"
)

var GitCommit string
var GitSummary string
var BuildDate string

type VersionInfo struct {
	GitCommit  string
	GitSummary string
	BuildDate  string
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	info := VersionInfo{
		GitCommit:  GitCommit,
		GitSummary: GitSummary,
		BuildDate:  BuildDate,
	}
	logBuildInfo(logger, info)

	ctx := context.Background()
	config.MustLoadFromEnv(ctx, logger)
}

func logBuildInfo(logger *zap.Logger, info VersionInfo) {
	logger.Info("Build info",
		zap.String("commit", info.GitCommit),
		zap.String("date", info.BuildDate),
		zap.String("version", info.GitSummary))
}

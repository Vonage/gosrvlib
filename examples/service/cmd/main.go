package main

import (
	"github.com/nexmoinc/gosrvlib-sample-service/internal/cli"
	"github.com/nexmoinc/gosrvlib/pkg/logging"
	"go.uber.org/zap"
)

var (
	Version = "0.0.0"
	Release = "0"
)

func main() {
	_, _ = logging.NewDefaultLogger(cli.AppName, Version, Release, "json", "debug")

	rootCmd, err := cli.New(Version, Release)
	if err != nil {
		zap.L().Fatal("UNABLE TO START THE PROGRAM", zap.Error(err))
		return
	}
	// execute the root command and log errors (if any)
	if err = rootCmd.Execute(); err != nil {
		zap.L().Fatal("UNABLE TO RUN THE COMMAND", zap.Error(err))
	}
}

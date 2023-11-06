package cmd

import (
	"context"
	"log/slog"
	"os"
	"soulmonk/dumper/pkg/config"
	"soulmonk/dumper/pkg/rest"
)

func RunServer() error {
	ctx := context.Background()
	cfg := config.Load()

	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.Level(cfg.LogLevel))
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return rest.RunServer(ctx, cfg.HTTPPort)
}

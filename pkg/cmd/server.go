package cmd

import (
	"context"
	"soulmonk/dumper/pkg/config"
	"soulmonk/dumper/pkg/rest"
)

func RunServer() error {
	ctx := context.Background()
	cfg := config.Load()
	return rest.RunServer(ctx, cfg.HTTPPort)
}

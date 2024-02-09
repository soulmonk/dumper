package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"soulmonk/dumper/pkg/db"
	"time"
)

func RunServer(ctx context.Context, httpPort string, dao *db.Dao) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := setupRouter(dao.IdeasQuerier)

	addr := httpPort
	slog.Debug("listen on", "addr", addr)
	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: r,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			slog.Warn("shutting down HTTP/REST gateway...")
			_, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			_ = srv.Shutdown(ctx)
			<-ctx.Done()
		}

	}()

	slog.Info("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}

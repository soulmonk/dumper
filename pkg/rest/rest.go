package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	// Add the sloggin middleware to all routes.
	// The middleware will log all requests attributes.
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	r.GET("/ping", ping)
	//r.GET("/ideas", getIdeas)
	//r.POST("/ideas", createIdea)
	//r.GET("/ideas/:id", getIdea)
	return r
}

func RunServer(ctx context.Context, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := setupRouter()

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

type pingResponse struct {
	Status string `bson:"status" json:"status"`
}

func ping(c *gin.Context) {
	var data = pingResponse{"ok"}
	c.IndentedJSON(http.StatusOK, data)
}

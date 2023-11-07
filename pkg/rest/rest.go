package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"soulmonk/dumper/pkg/db"
	"soulmonk/dumper/pkg/db/ideas"
	"time"
)

func setupRouter(dao *db.Dao) *gin.Engine {
	r := gin.New()
	// Add the sloggin middleware to all routes.
	// The middleware will log all requests attributes.
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	r.GET("/ping", ping)
	r.GET("/ideas", getGetIdeasHandler(dao))
	r.POST("/ideas", getCreateIdeaHandler(dao))
	return r
}

func getGetIdeasHandler(dao *db.Dao) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := dao.IdeasQuerier.ListIdeas(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(200, result)
	}
}
func getCreateIdeaHandler(dao *db.Dao) gin.HandlerFunc {
	return func(c *gin.Context) {
		var idea ideas.CreateIdeaParams
		if err := c.Bind(&idea); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		result, err := dao.IdeasQuerier.CreateIdea(c, idea)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(200, result)
	}
}

func RunServer(ctx context.Context, httpPort string, dao *db.Dao) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := setupRouter(dao)

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

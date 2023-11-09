package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"soulmonk/dumper/pkg/db"
	"soulmonk/dumper/pkg/db/ideas"
	"time"
)

func setupRouter(querier ideas.Querier) *gin.Engine {
	r := gin.New()
	// Add the sloggin middleware to all routes.
	// The middleware will log all requests attributes.
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	r.GET("/ping", ping)

	r.GET("/ideas", getGetIdeasHandler(querier))
	r.GET("/ideas/random", getRandomIdeasHandler(querier))
	r.POST("/ideas", getCreateIdeaHandler(querier))
	r.POST("/ideas/:id/done", getDoneIdeaHandler(querier))

	return r
}

func getRandomIdeasHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		ids, err := querier.GetIdsOfActiveIdeas(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result, err := querier.GetIdea(c, ids[rand.Intn(len(ids))])
		c.JSON(200, result)
	}
}

type IdeaId struct {
	ID int64 `uri:"id" binding:"required"`
}

func getDoneIdeaHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		var idea IdeaId
		if err := c.ShouldBindUri(&idea); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := querier.DoneIdea(c, idea.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
	}
}

func getGetIdeasHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := querier.ListIdeas(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
	}
}
func getCreateIdeaHandler(querier ideas.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		var idea ideas.CreateIdeaParams
		if err := c.ShouldBindJSON(&idea); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := querier.CreateIdea(c, idea)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
	}
}

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

type pingResponse struct {
	Status string `bson:"status" json:"status"`
}

func ping(c *gin.Context) {
	var data = pingResponse{"ok"}
	c.IndentedJSON(http.StatusOK, data)
}

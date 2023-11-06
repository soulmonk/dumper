package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"soulmonk/dumper/pkg/rest/middleware"
	"time"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", ping)
	return r
}

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := setupRouter()

	addr := httpPort
	log.Debug().Str("listen on", addr).Send()
	//if err := http.ListenAndServe(addr, r); err != nil {
	//	log.Fatal(err)
	//}

	srv := &http.Server{
		Addr: ":" + httpPort,
		Handler: middleware.AddRequestID(
			middleware.AddLogger(r)),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Warn().Msg("shutting down HTTP/REST gateway...")
			_, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			_ = srv.Shutdown(ctx)
			<-ctx.Done()
		}

	}()

	log.Info().Msg("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}

type pingResponse struct {
	Status string `bson:"status" json:"status"`
}

// getAlbums responds with the list of all albums as JSON.
func ping(c *gin.Context) {
	var data = pingResponse{"ok"}
	c.IndentedJSON(http.StatusOK, data)
}

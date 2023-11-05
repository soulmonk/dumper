package rest

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	middleware2 "soulmonk/dumper/pkg/rest/middleware"
	"time"
)

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	r := mux.NewRouter()

	r.HandleFunc("/ping", status).Methods("GET")

	addr := httpPort
	log.Debug().Str("listen on", addr).Send()
	//if err := http.ListenAndServe(addr, r); err != nil {
	//	log.Fatal(err)
	//}

	srv := &http.Server{
		Addr: ":" + httpPort,
		Handler: middleware2.AddRequestID(
			middleware2.AddLogger(r)),
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

type statusResponse struct {
	Status string `bson:"status" json:"status"`
}

func status(w http.ResponseWriter, r *http.Request) {
	var data = statusResponse{"ok"}
	RespondWithJson(w, http.StatusOK, data)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJson(w, code, map[string]string{"error": msg})
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		log.Fatal().Err(err)
	}
}

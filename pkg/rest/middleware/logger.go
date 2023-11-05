package middleware

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

// AddLogger logs request/response pair
func AddLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// We do not want to be spammed by Kubernetes health check.
		// Do not log Kubernetes health check.
		// You can change this behavior as you wish.
		if r.Header.Get("X-Liveness-Probe") == "Healthz" {
			h.ServeHTTP(w, r)
			return
		}

		id := GetReqID(ctx)

		// Prepare fields to log
		var scheme string
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
		proto := r.Proto
		method := r.Method
		remoteAddr := r.RemoteAddr
		userAgent := r.UserAgent()
		uri := strings.Join([]string{scheme, "://", r.Host, r.RequestURI}, "")

		// Log HTTP request
		log.Debug().
			Str("request-id", id).
			Str("http-scheme", scheme).
			Str("http-proto", proto).
			Str("http-method", method).
			Str("remote-addr", remoteAddr).
			Str("user-agent", userAgent).
			Str("uri", uri).
			Msg("request started")

		t1 := time.Now()

		h.ServeHTTP(w, r)

		// Log HTTP response
		log.Debug().
			Str("request-id", id).
			Str("http-scheme", scheme).
			Str("http-proto", proto).
			Str("http-method", method).
			Str("remote-addr", remoteAddr).
			Str("user-agent", userAgent).
			Str("uri", uri).
			Float64("elapsed-ms", float64(time.Since(t1).Nanoseconds())/1000000.0).
			Msg("request completed")

	})
}

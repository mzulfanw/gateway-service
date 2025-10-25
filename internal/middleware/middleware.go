package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ctxKey string

const RequestIDKey ctxKey = "requestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", rid)
		ctx := context.WithValue(r.Context(), RequestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logger(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			rid, _ := r.Context().Value(RequestIDKey).(string)
			log.WithFields(logrus.Fields{
				"request_id": rid,
				"method":     r.Method,
				"path":       r.URL.Path,
				"duration":   time.Since(start).String(),
			}).Info("request completed")
		})
	}
}

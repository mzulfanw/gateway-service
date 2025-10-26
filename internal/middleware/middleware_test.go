package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRequestID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := w.Header().Get("X-Request-ID")
		if rid == "" {
			t.Error("X-Request-ID header not set")
		}
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header missing in response")
	}
}

func TestLogger(t *testing.T) {
	log := logrus.New()
	log.SetOutput(io.Discard)

	handler := Logger(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mzulfanw/gateway-service/internal/proxy"
)

func TestCreateOrder(t *testing.T) {
	tp := proxy.NewProxy("", "http://localhost:3001")
	h := &OrderHandler{proxy: tp}
	req := httptest.NewRequest("POST", "/orders", nil)
	rr := httptest.NewRecorder()
	h.CreateOrder(rr, req)
	if rr.Code != http.StatusBadRequest && rr.Code != http.StatusOK {
		t.Errorf("expected status 200 or 503, got %d", rr.Code)
	}
}

func TestGetByProductID_Success(t *testing.T) {
	tp := proxy.NewProxy("", "http://localhost:3001")
	h := &OrderHandler{proxy: tp}
	req := httptest.NewRequest("GET", "/orders/product/04b18607-1a30-402b-8a7b-26dd5d5d6235", nil)
	vars := map[string]string{"productId": "04b18607-1a30-402b-8a7b-26dd5d5d6235"}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	h.GetByProductID(rr, req)
	if rr.Code != http.StatusServiceUnavailable && rr.Code != http.StatusOK {
		t.Errorf("expected status 200 or 503, got %d", rr.Code)
	}
}

func TestGetByProductID_MissingProductID(t *testing.T) {
	tp := proxy.NewProxy("", "http://localhost:3001")
	h := &OrderHandler{proxy: tp}
	req := httptest.NewRequest("GET", "/orders/product/", nil)
	rr := httptest.NewRecorder()
	h.GetByProductID(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

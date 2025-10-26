package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mzulfanw/gateway-service/internal/proxy"
)

func TestCreateProduct(t *testing.T) {
	tp := proxy.NewProxy("http://product", "")
	h := &ProductHandler{proxy: tp}
	req := httptest.NewRequest("POST", "/products", nil)
	rr := httptest.NewRecorder()
	h.CreateProduct(rr, req)
	if rr.Code != http.StatusServiceUnavailable && rr.Code != http.StatusOK {
		t.Errorf("expected status 200 or 503, got %d", rr.Code)
	}
}

func TestGetProductDetail_Success(t *testing.T) {
	tp := proxy.NewProxy("http://product", "")
	h := &ProductHandler{proxy: tp}
	req := httptest.NewRequest("GET", "/products/04b18607-1a30-402b-8a7b-26dd5d5d6235", nil)
	vars := map[string]string{"id": "123"}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	h.GetProductDetail(rr, req)
	if rr.Code != http.StatusServiceUnavailable && rr.Code != http.StatusOK {
		t.Errorf("expected status 200 or 503, got %d", rr.Code)
	}
}

func TestGetProductDetail_MissingID(t *testing.T) {
	tp := proxy.NewProxy("http://product", "")
	h := &ProductHandler{proxy: tp}
	req := httptest.NewRequest("GET", "/products/", nil)
	rr := httptest.NewRecorder()
	h.GetProductDetail(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mzulfanw/gateway-service/internal/proxy"
	"github.com/mzulfanw/gateway-service/internal/response"
)

type ProductHandler struct {
	proxy *proxy.Proxy
}

func NewProductHandler(p *proxy.Proxy) *ProductHandler {
	return &ProductHandler{proxy: p}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyWithContext(w, r, h.proxy.ProductServiceUrl+"/products")
}

func (h *ProductHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, ok := vars["id"]
	if !ok {
		response.ErrorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	targetURL := h.proxy.ProductServiceUrl + "/products/" + productID
	h.proxy.ProxyWithContext(w, r, targetURL)
}

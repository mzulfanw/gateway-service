package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mzulfanw/gateway-service/internal/proxy"
	"github.com/mzulfanw/gateway-service/internal/response"
)

type OrderHandler struct {
	proxy *proxy.Proxy
}

func NewOrderHandler(p *proxy.Proxy) *OrderHandler {
	return &OrderHandler{proxy: p}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	h.proxy.ProxyWithContext(w, r, h.proxy.OrderServiceUrl+"/orders")
}

func (h *OrderHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, ok := vars["productId"]
	if !ok {
		response.ErrorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	targetUrl := h.proxy.OrderServiceUrl + "/orders/product/" + productId
	h.proxy.ProxyWithContext(w, r, targetUrl)
}

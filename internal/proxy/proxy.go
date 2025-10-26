package proxy

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mzulfanw/gateway-service/internal/response"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Proxy struct {
	ProductServiceUrl string
	OrderServiceUrl   string
	client            HTTPClient
}

func NewProxy(productServiceUrl, orderServiceUrl string) *Proxy {
	return &Proxy{
		ProductServiceUrl: productServiceUrl,
		OrderServiceUrl:   orderServiceUrl,
		client:            &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *Proxy) ProxyWithContext(w http.ResponseWriter, r *http.Request, targetUrl string) {
	req, err := http.NewRequestWithContext(r.Context(), r.Method, targetUrl, r.Body)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	req.Header = r.Header.Clone()

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = p.client.Do(req)
		if err == nil {
			break
		}
		log.Printf("proxy error (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		response.ErrorResponse(w, http.StatusServiceUnavailable, "Service unavailable")
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

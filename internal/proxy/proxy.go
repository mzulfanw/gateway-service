package proxy

import (
	"bytes"
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
	client            *http.Client
}

var defaultTransport = &http.Transport{
	MaxIdleConns:        1000,
	MaxIdleConnsPerHost: 1000,
	IdleConnTimeout:     90 * time.Second,
}

func NewProxy(productServiceUrl, orderServiceUrl string) *Proxy {
	return &Proxy{
		ProductServiceUrl: productServiceUrl,
		OrderServiceUrl:   orderServiceUrl,
		client:            &http.Client{Timeout: 10 * time.Second, Transport: defaultTransport},
	}
}

func (p *Proxy) ProxyWithContext(w http.ResponseWriter, r *http.Request, targetUrl string) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var resp *http.Response
	for i := 0; i < 3; i++ {
		req, reqErr := http.NewRequestWithContext(r.Context(), r.Method, targetUrl, io.NopCloser(bytes.NewReader(bodyBytes)))
		if reqErr != nil {
			response.ErrorResponse(w, http.StatusInternalServerError, reqErr.Error())
			return
		}
		req.Header = r.Header.Clone()

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
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("error copying response body: %v", err)
	}
}

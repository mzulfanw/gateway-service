package proxy

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Proxy struct {
	ProductServiceUrl string
	OrderServiceUrl   string
	client            *http.Client
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
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

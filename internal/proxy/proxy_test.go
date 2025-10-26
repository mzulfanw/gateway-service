package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockClient struct {
	resp *http.Response
	err  error
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func TestProxyWithContext_Success(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString("proxied response"))
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       body,
		Header:     http.Header{"X-Test": []string{"value"}},
	}
	p := &Proxy{
		ProductServiceUrl: "http://product",
		OrderServiceUrl:   "http://order",
		client:            &mockClient{resp: mockResp, err: nil},
	}

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	p.ProxyWithContext(rr, req, "http://target")

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	if rr.Body.String() != "proxied response" {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
	if rr.Header().Get("X-Test") != "value" {
		t.Errorf("expected header X-Test=value, got %s", rr.Header().Get("X-Test"))
	}
}

func TestProxyWithContext_Error(t *testing.T) {
	p := &Proxy{
		ProductServiceUrl: "http://product",
		OrderServiceUrl:   "http://order",
		client:            &mockClient{resp: nil, err: http.ErrHandlerTimeout},
	}

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	p.ProxyWithContext(rr, req, "http://target")

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", rr.Code)
	}
}

// ...existing code...

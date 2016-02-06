package app

import (
	"log"
	"net/http"
	"sync"
)

var HTTPClientPool httpClientPool

func init() {
	HTTPClientPool.Put(NewHTTPClient("Mozilla/5.0 (X11; Linux x86_64; rv:42.0) Gecko/20100101 Firefox/42.0"))

	HTTPClientPool.Put(NewHTTPClient("Mozilla/5.0 (X11; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0"))
	log.Println("info: pool created.")
}

// Pool
// ----------------------

type httpClientPool struct {
	pool sync.Pool
}

func (p *httpClientPool) Put(client *HTTPClient) {
	p.pool.Put(client)
}

func (p *httpClientPool) Get() *HTTPClient {
	return p.pool.Get().(*HTTPClient)
}

// HTTP Client
// ----------------------

type HTTPClient struct {
	useragent string
	client    *http.Client
}

func NewHTTPClient(useragent string) *HTTPClient {
	return &HTTPClient{
		client:    &http.Client{},
		useragent: useragent,
	}
}

// NOTE(remy): only supports GET request
func (c HTTPClient) NewRequest(url string) (*http.Request, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", c.useragent)
	return r, nil
}

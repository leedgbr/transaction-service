package app

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// NewHttpClient creates a new http client with configured timeouts and MaxConnsPerHost
func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxConnsPerHost:       100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// NewListener creates a new net.Listener on the supplied port
func NewListener(port int) (net.Listener, error) {
	addr := fmt.Sprintf(":%d", port)
	return net.Listen("tcp", addr)
}

// newServer creates a new http server with the supplied http handler
func newServer(router http.Handler) *http.Server {
	return &http.Server{
		Handler: router,
	}
}

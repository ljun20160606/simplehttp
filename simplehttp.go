package simplehttp

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	Verbose           = false
	logger            = log.New(os.Stdout, "[http] ", log.LstdFlags)
	DefaultHttpClient = &http.Client{
		Transport: HTTP1.RoundTripper(),
	}
	DefaultDialContext = DialContext(30*time.Second, 30*time.Second, 0)
)

func DialContext(connTimeout, KeepAlive, rwTimeout time.Duration) func(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   connTimeout,
		KeepAlive: KeepAlive,
		DualStack: true,
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		if rwTimeout > 0 {
			timeoutConn := &tcpConn{
				TCPConn: c.(*net.TCPConn),
				timeout: rwTimeout,
			}
			return timeoutConn, nil
		}
		return c, nil
	}
}

// quick
func CookieJar() http.CookieJar {
	return DefaultHttpClient.Jar
}

func Get() *Request {
	return &Request{Header: http.Header{}, Method: http.MethodGet, Client: DefaultClient}
}

func Post() *Request {
	return &Request{Header: http.Header{}, Method: http.MethodPost, Client: DefaultClient}
}

func Delete() *Request {
	return &Request{Header: http.Header{}, Method: http.MethodDelete, Client: DefaultClient}
}

func Put() *Request {
	return &Request{Header: http.Header{}, Method: http.MethodPut, Client: DefaultClient}
}

func Patch() *Request {
	return &Request{Header: http.Header{}, Method: http.MethodPatch, Client: DefaultClient}
}

func Head() *Request {
	return &Request{Header: http.Header{}, Method: http.MethodHead, Client: DefaultClient}
}

func Options() *Request {
	return &Request{Header: http.Header{}, Method: http.MethodOptions, Client: DefaultClient}
}

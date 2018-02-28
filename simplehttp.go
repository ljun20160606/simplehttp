package simplehttp

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
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

func Get(url string) *Request {
	return &Request{Header: http.Header{}, Url: newStringBuilder(url), Method: http.MethodGet, Client: DefaultClient}
}

func Post(url string) *Request {
	return &Request{Header: http.Header{}, Url: newStringBuilder(url), Method: http.MethodPost, Client: DefaultClient}
}

func Delete(url string) *Request {
	return &Request{Header: http.Header{}, Url: newStringBuilder(url), Method: http.MethodDelete, Client: DefaultClient}
}

func Put(url string) *Request {
	return &Request{Header: http.Header{}, Url: newStringBuilder(url), Method: http.MethodPut, Client: DefaultClient}
}

func Patch(url string) *Request {
	return &Request{Header: http.Header{}, Url: newStringBuilder(url), Method: http.MethodPatch, Client: DefaultClient}
}

func Head(url string) *Request {
	return &Request{Header: http.Header{}, Url: newStringBuilder(url), Method: http.MethodHead, Client: DefaultClient}
}

func Options(url string) *Request {
	return &Request{Header: http.Header{}, Url: newStringBuilder(url), Method: http.MethodOptions, Client: DefaultClient}
}

func newStringBuilder(str string) *strings.Builder {
	builder := strings.Builder{}
	builder.WriteString(str)
	return &builder
}

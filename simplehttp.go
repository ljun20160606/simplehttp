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
	return Method(http.MethodGet)
}

func Post() *Request {
	return Method(http.MethodPost)
}

func Delete() *Request {
	return Method(http.MethodDelete)
}

func Put() *Request {
	return Method(http.MethodPut)
}

func Patch() *Request {
	return Method(http.MethodPatch)
}

func Head() *Request {
	return Method(http.MethodHead)
}

func Options() *Request {
	return Method(http.MethodOptions)
}

func Method(m string) *Request {
	return NewPureRequest().SetClient(DefaultClient).SetMethod(m)
}

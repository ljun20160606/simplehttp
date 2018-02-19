package simplehttp

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Verbose            = false
	logger             = log.New(os.Stdout, "[http] ", log.LstdFlags)
	DefaultClient      = NewNoSSLVerify()
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

type tcpConn struct {
	*net.TCPConn
	timeout time.Duration
}

func (c *tcpConn) Read(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Read(b)
}

func (c *tcpConn) Write(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Write(b)
}

func NewNoSSLVerify() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           DefaultDialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// quick
func CookieJar() http.CookieJar {
	return DefaultClient.Jar
}

func Get(url string) *Request {
	return &Request{header: http.Header{}, url: newStringBuilder(url), method: http.MethodGet, Driver: DefaultDriver}
}

func Post(url string) *Request {
	return &Request{header: http.Header{}, url: newStringBuilder(url), method: http.MethodPost, Driver: DefaultDriver}
}

func Delete(url string) *Request {
	return &Request{header: http.Header{}, url: newStringBuilder(url), method: http.MethodDelete, Driver: DefaultDriver}
}

func Put(url string) *Request {
	return &Request{header: http.Header{}, url: newStringBuilder(url), method: http.MethodPut, Driver: DefaultDriver}
}

func Patch(url string) *Request {
	return &Request{header: http.Header{}, url: newStringBuilder(url), method: http.MethodPatch, Driver: DefaultDriver}
}

func Head(url string) *Request {
	return &Request{header: http.Header{}, url: newStringBuilder(url), method: http.MethodHead, Driver: DefaultDriver}
}

func Options(url string) *Request {
	return &Request{header: http.Header{}, url: newStringBuilder(url), method: http.MethodOptions, Driver: DefaultDriver}
}

func newStringBuilder(str string) *strings.Builder {
	builder := strings.Builder{}
	builder.WriteString(str)
	return &builder
}

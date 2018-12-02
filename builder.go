package simplehttp

import (
	"encoding/json"
	"github.com/ljun20160606/cookiejar"
	"github.com/ljun20160606/simplehttp/cache"
	"golang.org/x/net/http2"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultSessionTimeout = 10 * time.Minute
)

type (
	StoreCookie func(jar http.CookieJar)

	Proxy func(*http.Request) (*url.URL, error)

	Builder struct {
		// Timeout specifies a time limit for requests made by this
		// Client. The timeout includes connection time, any
		// redirects, and reading the response body. The timer remains
		// running after Get, Head, Post, or Do return and will
		// interrupt reading of the Response.Body.
		//
		// A Timeout of zero means no timeout.
		//
		// The Client cancels requests to the underlying Transport
		// using the Request.Cancel mechanism. Requests passed
		// to Client.Do may still set Request.Cancel; both will
		// cancel the request.
		//
		// For compatibility, the Client will also use the deprecated
		// CancelRequest method on Transport if found. New
		// RoundTripper implementations should use Request.Cancel
		// instead of implementing CancelRequest.
		Timeout time.Duration

		// Would don't use cache If nil
		Cache cache.Cache

		// Default Http1
		ProtoMajor ProtoMajor

		Proxy          Proxy
		SessionID      string
		SessionTimeout time.Duration
		Client         *http.Client
	}
)

func (b *Builder) Build() Client {
	c := b.client()
	// timeout
	c.Timeout = b.Timeout
	// proxy
	proxy := b.proxy()
	if proxy != nil {
		if tr, ok := c.Transport.(*http.Transport); ok {
			tr.Proxy = proxy
			tr.ExpectContinueTimeout = 0
		}
	}
	// protoMajor
	if b.ProtoMajor == HTTP2 {
		if tr, ok := c.Transport.(*http.Transport); ok {
			err := http2.ConfigureTransport(tr)
			if err != nil {
				logger.Fatal("[client-err]", err)
			}
		}
	}
	// httpClient
	client := &HttpClient{Client: c}
	// cookie
	if b.Cache != nil {
		c.Jar = b.loadCookie()
		client.StoreCookie = b.storeCookie
	}
	return client
}

func (b *Builder) client() *http.Client {
	if b.Client != nil {
		return b.Client
	}
	return &http.Client{
		Transport: HTTP1.RoundTripper(),
	}
}

func (b *Builder) proxy() Proxy {
	if b.Proxy != nil {
		return b.Proxy
	}
	return http.ProxyFromEnvironment
}

func (b *Builder) loadCookie() (jar http.CookieJar) {
	cookieJarBytes, _ := b.Cache.Get(b.SessionID)
	options := &cookiejar.Options{PublicSuffixList: publicsuffix.List}
	var err error
	if cookieJarBytes == nil {
		jar, err = cookiejar.New(options)
	} else {
		jar, err = cookiejar.LoadFromJson(options, cookieJarBytes)
	}
	if err != nil {
		logger.Fatal("[cookie-jar-err]", err)
	}
	return jar
}

func (b *Builder) storeCookie(cookieJar http.CookieJar) {
	data, _ := json.Marshal(cookieJar)
	b.Cache.Set(b.SessionID, data, b.sessionTimeout())
}

func (b *Builder) sessionTimeout() time.Duration {
	if b.SessionTimeout == 0 {
		return defaultSessionTimeout
	}
	return b.SessionTimeout
}

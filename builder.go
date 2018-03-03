package simplehttp

import (
	"encoding/json"
	"github.com/ljun20160606/cookiejar"
	"github.com/ljun20160606/simplehttp/cache"
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
	c.Timeout = b.Timeout
	proxy := b.proxy()
	if proxy != nil {
		if tr, ok := c.Transport.(*http.Transport); ok {
			tr.Proxy = proxy
			tr.ExpectContinueTimeout = 0
		}
	}
	client := &HttpClient{Client: c}
	if b.Cache != nil {
		b.loadCookie(c)
		client.StoreCookie = b.storeCookie
	}
	return client
}

func (b *Builder) client() *http.Client {
	if b.Client != nil {
		return b.Client
	}
	if b.ProtoMajor == 0 {
		b.ProtoMajor = HTTP1
	}
	return &http.Client{
		Transport: b.ProtoMajor.RoundTripper(),
	}
}

func (b *Builder) proxy() Proxy {
	if b.Proxy != nil {
		return b.Proxy
	}
	return http.ProxyFromEnvironment
}

func (b *Builder) loadCookie(client *http.Client) {
	cookieJarBytes := b.loadCache()
	if cookieJarBytes == nil {
		Jars, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			logger.Fatal("[cookie-jar-err]", err)
		}
		client.Jar = Jars
	} else {
		Jars, err := cookiejar.LoadFromJson(&cookiejar.Options{PublicSuffixList: publicsuffix.List}, cookieJarBytes)
		if err != nil {
			logger.Fatal("[cookie-jar-err]", err)
		}
		client.Jar = Jars
	}
}

func (b *Builder) loadCache() (data []byte) {
	if cacheData, found := b.Cache.Get(b.SessionID); found && cacheData != nil {
		data = cacheData
		return
	}
	return
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

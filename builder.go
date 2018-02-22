package simplehttp

import (
	"crypto/tls"
	"encoding/json"
	"github.com/ljun20160606/cookiejar"
	"github.com/ljun20160606/simplehttp/cache"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultCheckRedirect = func(req *http.Request, via []*http.Request) error {
		if Verbose {
			logger.Println(req.URL)
		}
		return nil
	}
)

const (
	defaultSessionTimeout = 10 * time.Minute
)

type (
	StoreCookie func(jar http.CookieJar)

	Proxy func(*http.Request) (*url.URL, error)

	Builder struct {
		SessionTimeout time.Duration
		Timeout        time.Duration
		Proxy          Proxy
		Cache          cache.Cache
		SessionID      string
		Client         *http.Client
	}
)

func (b *Builder) Build() Client {
	c := b.client()
	c.Timeout = b.Timeout
	if tr, ok := c.Transport.(*http.Transport); ok {
		tr.Proxy = b.proxy()
		tr.ExpectContinueTimeout = 0
	} else {
		logger.Fatal("[builder-err]", "can't convert Transport to *http.Transport")
	}
	c.CheckRedirect = defaultCheckRedirect
	b.loadCookie(c)
	return &client{Client: c, StoreCookie: b.storeCookie}
}

func (b *Builder) client() *http.Client {
	if b.Client != nil {
		return b.Client
	}
	return newNoSSLVerifyClient()
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

func newNoSSLVerifyClient() *http.Client {
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

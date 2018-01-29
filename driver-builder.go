package simplehttp

import (
	"encoding/json"
	"github.com/ljun20160606/cookiejar"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/url"
	"time"
)

var (
	emptyDuration        time.Duration
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
	Proxy       func(*http.Request) (*url.URL, error)
)

type Builder struct {
	SessionTimeout time.Duration
	Timeout        time.Duration
	Proxy          Proxy
	Cache          Cache
	SessionID      string
	Client         *http.Client
}

func (b *Builder) sessionTimeout() time.Duration {
	if b.SessionTimeout == emptyDuration {
		return defaultSessionTimeout
	}
	return b.SessionTimeout
}

func (b *Builder) proxy() Proxy {
	if b.Proxy != nil {
		return b.Proxy
	}
	return http.ProxyFromEnvironment
}

func (b *Builder) client() *http.Client {
	if b.Client != nil {
		return b.Client
	}
	return NewNoSSLVerify()
}

func (b *Builder) saveCache(data []byte) {
	b.Cache.Set(b.SessionID, data, b.sessionTimeout())
	return
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
	b.saveCache(data)
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

func (b *Builder) Build() Driver {
	client := b.client()
	client.Timeout = b.Timeout
	if tr, ok := client.Transport.(*http.Transport); ok {
		tr.Proxy = b.proxy()
		tr.ExpectContinueTimeout = 0
	} else {
		logger.Fatal("[builder-err]", "can't convert Transport to *http.Transport")
	}
	client.CheckRedirect = defaultCheckRedirect
	b.loadCookie(client)
	return &HttpDriver{Client: client, StoreCookie: b.storeCookie}
}

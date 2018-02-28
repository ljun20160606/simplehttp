package simplehttp

import (
	"bytes"
	"encoding/json"
	"github.com/ljun20160606/cookiejar"
	"github.com/ljun20160606/simplehttp/simplehttputil"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"io/ioutil"
	"net/http"
	"strings"
)

var DefaultClient = &client{Client: DefaultHttpClient}

type Client interface {
	Send(request *Request) *Response
}

func NewClient(builder *Builder) Client {
	return builder.Build()
}

type client struct {
	*http.Client
	StoreCookie StoreCookie
}

func (h *client) Send(r *Request) (resp *Response) {
	var err error
	if r.Querys != nil {
		r.Url.WriteByte('?')
		r.Url.Write(simplehttputil.BuildQueryEncoded(r.Querys, r.Charset))
	}
	switch {
	case r.Body != nil:
	case r.Forms != nil:
		r.Body = bytes.NewReader(simplehttputil.BuildFormEncoded(r.Forms, r.Charset))
	case r.JsonData != nil:
		body, err := json.Marshal(r.JsonData)
		if err != nil {
			resp.err = err
			return
		}
		r.Body = bytes.NewReader(body)
	}
	resp = new(Response)
	realReq, err := http.NewRequest(r.Method, r.Url.String(), r.Body)
	if err != nil {
		resp.err = err
		return
	}
	realReq.Header = r.Header
	if r.IsClearCookie || h.Client.Jar == nil {
		h.Client.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	}
	if r.Cookies != nil {
		h.Client.Jar.SetCookies(realReq.URL, r.Cookies)
	}
	switch r.Retry {
	case 0:
		resp = h.send(realReq)
	default:
		for times := -1; times < r.Retry; times++ {
			resp = h.send(realReq)
			if resp.err == nil || !strings.Contains(resp.err.Error(), "request canceled") {
				break
			}
		}
	}
	return resp
}

func (h *client) send(realReq *http.Request) (resp *Response) {
	response, err := h.Client.Do(realReq)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	resp = new(Response)
	if err != nil {
		resp.err = err
		return
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp.err = err
		return
	}
	if h.StoreCookie != nil {
		h.StoreCookie(h.Client.Jar)
	}
	var enc encoding.Encoding
	contentType := response.Header.Get(ContentType)
	if len(contentType) > 0 {
		subMatch := ContentTypeMatchCharset.FindStringSubmatch(contentType)
		var name string
		if len(subMatch) == 2 {
			if Verbose {
				logger.Println("find html Encode ", subMatch[1])
			}
			name = subMatch[1]
		}
		if name != "" {
			enc, _ = htmlindex.Get(name)
		}
	}
	return &Response{body: data, encoding: enc, Response: response}
}

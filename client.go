package simplehttp

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"github.com/ljun20160606/cookiejar"
	"github.com/ljun20160606/simplehttp/simplehttputil"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var DefaultClient = &HttpClient{Client: DefaultHttpClient}

type Client interface {
	Send(request *Request) *Response
}

func NewClient(builder *Builder) Client {
	return builder.Build()
}

type HttpClient struct {
	*http.Client
	StoreCookie StoreCookie
}

func (h *HttpClient) Send(r *Request) (resp *Response) {
	// request
	realReq, err := castToHttpRequest(r)
	if err != nil {
		resp = failResponse(err)
		return
	}
	// cookie
	h.prepareCookie(realReq.URL, r)
	switch r.Config.Retry {
	case 0:
		resp = h.emit(realReq)
	default:
		for times := -1; times < r.Config.Retry; times++ {
			resp = h.emit(realReq)
			if resp.err == nil || !strings.Contains(resp.err.Error(), "request canceled") {
				break
			}
		}
	}
	return resp
}

// cast simplehttp.Request to http.Request
func castToHttpRequest(r *Request) (*http.Request, error) {
	// url
	buildUrl(r)
	// body
	switch {
	case r.Body != nil:
	case r.Forms != nil:
		r.Body = bytes.NewReader(simplehttputil.BuildQueryEncoded(r.Forms, r.Charset))
	case r.JsonData != nil:
		body, err := json.Marshal(r.JsonData)
		if err != nil {
			return nil, errors.Wrap(err, "http body json fail marshal")
		}
		r.Body = bytes.NewReader(body)
	}
	realReq, err := http.NewRequest(r.Method, r.Url.String(), r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to cast http.Request")
	}
	realReq.Header = r.Header
	return realReq, nil
}

func buildUrl(r *Request) string {
	if r.Querys != nil {
		r.Url.WriteByte('?')
		r.Url.Write(simplehttputil.BuildQueryEncoded(r.Querys, r.Charset))
	}
	return r.Url.String()
}

// according simplehttp.Request config
func (h *HttpClient) prepareCookie(URL *url.URL, request *Request) {
	if request.Config.IsClearCookie || h.Jar == nil {
		h.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	}
	if request.Cookies != nil {
		h.Jar.SetCookies(URL, request.Cookies)
	}
}

func (h *HttpClient) emit(realReq *http.Request) *Response {
	response, err := h.Do(realReq)
	if err != nil {
		return failResponse(err)
	}
	defer func() {
		response.Body.Close()
		if h.StoreCookie != nil {
			h.StoreCookie(h.Jar)
		}
	}()
	return castToResponse(response)
}

func castToResponse(response *http.Response) *Response {
	reader := deCompress(response.Header.Get("content-encoding"), response.Body)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return failResponse(err)
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

func failResponse(err error) *Response {
	resp := new(Response)
	resp.err = err
	return resp
}

func deCompress(contentEncoding string, reader io.ReadCloser) (r io.ReadCloser) {
	switch strings.ToLower(contentEncoding) {
	case "gzip":
		r, _ = gzip.NewReader(reader)
		return
	case "deflate":
		r = flate.NewReader(reader)
		return
	default:
		return reader
	}
}

package simplehttp

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"net/http"
	"net/url"
	"regexp"
	"unsafe"
)

var ContentTypeMatchCharset = regexp.MustCompile(`[cC]harset=([\w|\-]*)`)

type Response struct {
	code     int
	err      error
	header   http.Header
	body     []byte
	url      *url.URL
	encoding encoding.Encoding
}

func (r *Response) Document() (doc *goquery.Document, err error) {
	if r.err != nil {
		return nil, r.err
	}
	enc := r.encoding
	data := r.body
	if enc != nil {
		data, err = enc.NewDecoder().Bytes(r.body)
		if err != nil {
			return nil, err
		}
	}
	return goquery.NewDocumentFromReader(bytes.NewReader(data))
}

func (r *Response) DetectedEncode() (err error) {
	if r.err != nil {
		return r.err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.body))
	if err != nil {
		return
	}
	var name string
	selector := doc.Find("meta[Charset]")
	if selector != nil && selector.Size() > 0 {
		name, _ = selector.Attr("Charset")
		if Verbose {
			logger.Println("find html Encode ", name)
		}
	} else {
		selector = doc.Find(`meta[http-equiv="Content-Type"]`)
		if selector != nil && selector.Size() > 0 {
			attrContent, exists := selector.Attr("content")
			if Verbose {
				logger.Println("find html Encode from content ", attrContent)
			}
			if exists {
				subMatch := ContentTypeMatchCharset.FindStringSubmatch(attrContent)
				if len(subMatch) == 2 {
					if Verbose {
						logger.Println("find html Encode ", subMatch[1])
					}
					name = subMatch[1]
				}
			}
		}
	}
	if name != "" {
		enc, err := htmlindex.Get(name)
		if err != nil {
			return err
		}
		r.encoding = enc
	}
	return nil
}

func (r *Response) DocumentDetectedEncode() (doc *goquery.Document, err error) {
	if r.err != nil {
		return nil, r.err
	}
	enc := r.encoding
	if enc == nil {
		err = r.DetectedEncode()
		if err != nil {
			return nil, err
		}
	}
	data := r.body
	if enc != nil {
		data, err = enc.NewDecoder().Bytes(r.body)
		if err != nil {
			return nil, err
		}
	}
	return goquery.NewDocumentFromReader(bytes.NewReader(data))
}

func (r *Response) Code() (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.code, nil
}

func (r *Response) Body() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.body, nil
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) URL() *url.URL {
	return r.url
}

func (r *Response) String() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	if r.encoding == nil {
		return *((*string)(unsafe.Pointer(&r.body))), nil
	}
	data, err := r.encoding.NewDecoder().Bytes(r.body)
	if err != nil {
		return "", err
	}
	str := *((*string)(unsafe.Pointer(&data)))
	if Verbose {
		logger.Println(str)
	}
	return str, err
}

func (r *Response) JSON(data interface{}) error {
	if r.err != nil {
		return r.err
	}
	if Verbose {
		logger.Println(string(r.body))
	}
	return json.Unmarshal(r.body, data)
}

func (r *Response) Bytes() []byte {
	return r.body
}

// http://www.w3.org/TR/encoding
func (r *Response) Encode(name string) *Response {
	if r.err != nil {
		return r
	}
	enc, err := htmlindex.Get(name)
	if err != nil {
		r.err = err
		return r
	}
	r.encoding = enc
	return r
}

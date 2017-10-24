package simplehttputil

import (
	"bytes"
	"net/url"
)

func BuildQueryEncoded(source [][2]string, charset Charset) []byte {
	var buf bytes.Buffer
	switch length := len(source); {
	case length == 0:
	case length > 1:
		for _, kv := range source[:length-1] {
			k, v := kv[0], kv[1]
			buf.WriteString(k)
			buf.WriteByte('=')
			charset.Encode(&v)
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
		fallthrough
	default:
		kv := source[length-1]
		k, v := kv[0], kv[1]
		buf.WriteString(k)
		buf.WriteByte('=')
		charset.Encode(&v)
		buf.WriteString(url.QueryEscape(v))
	}
	return buf.Bytes()
}

func BuildFormEncoded(source map[string][]string, charset Charset) []byte {
	var buf bytes.Buffer
	for k, strs := range source {
		for _, v := range strs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			charset.Encode(&v)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.Bytes()
}

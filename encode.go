package simplehttp

import (
	"bytes"
	"net/url"
)

func BuildQueryEncoded(source [][2]string, charset Charset) ([]byte, error) {
	var buf bytes.Buffer
	switch length := len(source); {
	case length == 0:
	case length > 1:
		for _, kv := range source[:length-1] {
			k, v := kv[0], kv[1]
			buf.WriteString(k)
			buf.WriteByte('=')
			err := charset.Encode(&v)
			if err != nil {
				return nil, err
			}
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
		fallthrough
	default:
		kv := source[length-1]
		k, v := kv[0], kv[1]
		buf.WriteString(k)
		buf.WriteByte('=')
		err := charset.Encode(&v)
		if err != nil {
			return nil, err
		}
		buf.WriteString(url.QueryEscape(v))
	}
	return buf.Bytes(), nil
}

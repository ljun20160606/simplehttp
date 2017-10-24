package simplehttputil

import "golang.org/x/text/encoding/simplifiedchinese"

type (
	Charset    int
	encodeFunc func(s *string) error
)

const (
	UTF8 Charset = iota
	GB18030
)

var encodeArray = [...]encodeFunc{
	defaultEncode,
	gb18030Encode,
}

func (c Charset) Encode(s *string) error {
	return encodeArray[c](s)
}

func defaultEncode(s *string) error {
	return nil
}

func gb18030Encode(s *string) error {
	v, err := simplifiedchinese.GB18030.NewEncoder().String(*s)
	if err != nil {
		return err
	}
	*s = v
	return nil
}

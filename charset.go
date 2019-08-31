package simplehttp

import "golang.org/x/text/encoding/simplifiedchinese"

type Charset interface {
	Encode(s *string) error
}

type CharsetEncoderFunc func(s *string) error

func (c CharsetEncoderFunc) Encode(s *string) error {
	return c(s)
}

var (
	UTF8    Charset = CharsetEncoderFunc(func(s *string) error { return nil })
	GB18030 Charset = CharsetEncoderFunc(func(s *string) error {
		v, err := simplifiedchinese.GB18030.NewEncoder().String(*s)
		if err != nil {
			return err
		}
		*s = v
		return nil
	})
)

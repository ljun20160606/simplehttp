package cache

import "time"

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, exp time.Duration)
}

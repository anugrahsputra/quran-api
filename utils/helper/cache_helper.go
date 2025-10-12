package helper

import (
	"fmt"
	"time"

	"github.com/op/go-logging"
	"github.com/patrickmn/go-cache"
)

var logger = logging.MustGetLogger("repository")

func GetOrSetCache[T any](
	c *cache.Cache,
	key string,
	ttl time.Duration,
	fetch func() (T, error),
) (T, error) {
	var zero T

	if cached, found := c.Get(key); found {
		logger.Infof("[CACHE] hit for key: %s", key)
		if val, ok := cached.(T); ok {
			return val, nil
		}
	} else {
		logger.Infof("[CACHE] miss for key: %s", key)
	}

	value, err := fetch()
	if err != nil {
		if cached, found := c.Get(key); found {
			if val, ok := cached.(T); ok {
				return val, nil
			}
		}
		return zero, fmt.Errorf("fetch failed: %w", err)
	}

	c.Set(key, value, ttl)

	return value, nil
}

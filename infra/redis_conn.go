package infra

import (
	"github.com/redis/go-redis/v9"
)

func RedisConn(url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
}

package redis

import (
	"github.com/gomodule/redigo/redis"
	"os"
)

func NewRedis() (redis.Conn, error) {
	redisHost := os.Getenv("REDIS_HOST")
    redisPort := os.Getenv("REDIS_PORT")

	redis, err := redis.Dial("tcp", redisHost + ":" + redisPort)
	if err != nil {
		return nil, err
	}

	return redis, nil
}

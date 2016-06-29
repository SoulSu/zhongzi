package redis

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool

var (
	Server = "123.184.17.14:6379"
	Password = "c47d13bdcbb63cf0f38159f9e44c0f38b63efe86"
)

func init() {
	redisPool = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", Password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Get() redis.Conn {
	return redisPool.Get()
}


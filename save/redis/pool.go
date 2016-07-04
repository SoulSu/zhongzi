package redis

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool

func redis_init() {
	redisPool = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", *REDIS_HOST)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", *REDIS_PWD); err != nil {
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


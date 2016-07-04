package redis

import "flag"

var REDIS_HOST *string

var REDIS_PWD *string

func init() {

	REDIS_HOST = flag.String("rh", ":6379", "connect redis host ")

	REDIS_PWD = flag.String("rpwd", "", "redis pwd")

	redis_init()
}

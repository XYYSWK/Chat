package operate

import "github.com/go-redis/redis/v8"

type RDB struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *RDB {
	return &RDB{rdb: rdb}
}

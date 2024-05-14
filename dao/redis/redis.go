package redis

import (
	"Chat/dao/redis/operate"
	"context"
	"github.com/go-redis/redis/v8"
)

func Init(Addr, Password string, PoolSize, DB int) *operate.RDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Addr,     // host:port
		Password: Password, // 密码
		PoolSize: PoolSize, // 连接池
		DB:       DB,       // 默认连接数据库（0-15）
	})
	_, err := rdb.Ping(context.Background()).Result() //测试连接
	if err != nil {
		panic(err)
	}
	return operate.New(rdb)
}

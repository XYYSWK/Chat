package settings

import (
	"Chat/dao"
	"Chat/dao/postgresql"
	"Chat/dao/redis"
	"Chat/global"
)

type database struct {
}

// Init 数据库(持久化层)初始化
func (d database) Init() {
	// mysql 初始化
	dao.Database.DB = postgresql.Init(global.PrivateSetting.Postgresql.SourceName)
	// redis 初始化
	dao.Database.Redis = redis.Init(
		global.PrivateSetting.Redis.Address,
		global.PrivateSetting.Redis.Password,
		global.PrivateSetting.Redis.PoolSize,
		global.PrivateSetting.Redis.DB,
	)
}

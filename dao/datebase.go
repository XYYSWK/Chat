package dao

import (
	"Chat/dao/postgresql"
	"Chat/dao/redis/operate"
)

type database struct {
	DB    postgresql.DB
	Redis *operate.RDB
}

var Database = new(database)

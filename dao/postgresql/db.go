package postgresql

import (
	db "Chat/dao/postgresql/sqlc"
	"Chat/dao/postgresql/tx"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB interface {
	db.Querier
	tx.TXer
}

func Init(dataSourceName string) DB {
	//创建连接池
	pool, err := pgxpool.Connect(context.Background(), dataSourceName)
	if err != nil {
		panic(err)
	}
	return &tx.SqlStore{Queries: db.New(pool), Pool: pool}
}

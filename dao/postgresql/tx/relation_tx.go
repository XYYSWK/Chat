package tx

import (
	db "Chat/dao/postgresql/sqlc"
	"Chat/dao/redis/operate"
	"context"
)

// DeleteRelationWithTx 从数据库中删除关系并删除 redis 中的关系
func (store *SqlStore) DeleteRelationWithTx(ctx context.Context, rdb *operate.RDB, relationID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.DeleteRelation(ctx, relationID)
		if err != nil {
			return err
		}
		return rdb.DeleteRelations(ctx, relationID)
	})
}

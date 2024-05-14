package tx

import (
	db "Chat/dao/postgresql/sqlc"
	"Chat/dao/redis/operate"
	"Chat/pkg/tool"
	"context"
)

// AddSettingWithTx 向数据库和 redis 中同时添加群成员
func (store *SqlStore) AddSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64, isLeader bool) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.CreateSetting(ctx, &db.CreateSettingParams{
			AccountID:  accountID,
			RelationID: relationID,
			IsLeader:   isLeader,
			IsSelf:     false,
		})
		if err != nil {
			return err
		}
		return rdb.AddRelationAccount(ctx, relationID, accountID)
	})
}

// TransferGroupWithTx 转让群
func (store *SqlStore) TransferGroupWithTx(ctx context.Context, accountID, relationID, toAccountID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		err = tool.DoThat(err, func() error {
			// 将原群主的 isLeader 转换为 false
			return queries.TransferIsLeaderFalse(ctx, &db.TransferIsLeaderFalseParams{
				RelationID: relationID,
				AccountID:  accountID,
			})
		})
		err = tool.DoThat(err, func() error {
			// 将新群主的 isLeader 转换为 true
			return queries.TransferIsLeaderTrue(ctx, &db.TransferIsLeaderTrueParams{
				RelationID: relationID,
				AccountID:  toAccountID,
			})
		})
		return err
	})
}

// DeleteSettingWithTx 从数据库和 redis 中删除群员
func (store *SqlStore) DeleteSettingWithTx(ctx context.Context, rdb *operate.RDB, accountID, relationID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		err := queries.DeleteSetting(ctx, &db.DeleteSettingParams{
			AccountID:  accountID,
			RelationID: relationID,
		})
		if err != nil {
			return err
		}
		return rdb.DeleteRelationAccount(ctx, relationID, accountID)
	})
}

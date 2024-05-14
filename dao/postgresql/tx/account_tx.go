package tx

import (
	db "Chat/dao/postgresql/sqlc"
	"Chat/dao/redis/operate"
	"Chat/pkg/tool"
	"context"
	"errors"
)

var (
	ErrAccountOverNum     = errors.New("账户数量超过限制")
	ErrAccountNameExists  = errors.New("账户名已存在")
	ErrAccountGroupLeader = errors.New("账户是群主")
)

// CreateAccountWithTx 检查数量、账户名之后创建账户并建立和自己的关系
func (store *SqlStore) CreateAccountWithTx(ctx context.Context, rdb *operate.RDB, maxAccountNum int32, arg *db.CreateAccountParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		var accountNum int32
		// 检查数量
		err = tool.DoThat(err, func() error {
			accountNum, err = queries.CountAccountsByUserID(ctx, arg.UserID)
			return err
		})
		if accountNum >= maxAccountNum {
			return ErrAccountOverNum
		}
		// 检查账户名
		var exists bool
		err = tool.DoThat(err, func() error {
			exists, err = queries.ExistsAccountByNameAndUserID(ctx, &db.ExistsAccountByNameAndUserIDParams{
				UserID: arg.UserID,
				Name:   arg.Name,
			})
			return err
		})
		if exists {
			return ErrAccountNameExists
		}
		// 创建账户
		err = tool.DoThat(err, func() error {
			return queries.CreateAccount(ctx, arg)
		})
		// 建立关系(自己与自己的好友关系)
		var relationID int64
		err = tool.DoThat(err, func() error {
			relationID, err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: arg.ID,
				Account2ID: arg.ID,
			})
			return err
		})
		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  arg.ID,
				RelationID: relationID,
				IsSelf:     true,
			})
		})
		// 添加自己一个人的关系到 redis
		err = tool.DoThat(err, func() error {
			return rdb.AddRelationAccount(ctx, relationID, arg.ID)
		})
		return err
	})
}

func (store *SqlStore) DeleteAccountWithTx(ctx context.Context, rdb *operate.RDB, accountID int64) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		// 判断该账户是否是群主
		var isLeader bool
		err = tool.DoThat(err, func() error {
			isLeader, err = queries.ExistsGroupLeaderByAccountIDWithLock(ctx, accountID)
			return err
		})
		if isLeader {
			return ErrAccountGroupLeader
		}
		// 删除好友
		var friendRelationIDs []int64
		err = tool.DoThat(err, func() error {
			friendRelationIDs, err = queries.DeleteFriendRelationsByAccountID(ctx, accountID)
			return err
		})
		// 删除群
		var groupRelationIDs []int64
		err = tool.DoThat(err, func() error {
			groupRelationIDs, err = queries.DeleteSettingsByAccountID(ctx, accountID)
			return err
		})
		// 删除账户
		err = tool.DoThat(err, func() error {
			err = queries.DeleteAccount(ctx, accountID)
			return err
		})
		// 从 redis 中删除对应的关系
		// 从 redis 中删除该账户的好友关系
		err = tool.DoThat(err, func() error {
			return rdb.DeleteRelations(ctx, friendRelationIDs...)
		})
		// 在 redis 中删除该账户所在的群聊中的该账户
		err = tool.DoThat(err, func() error {
			return rdb.DeleteAccountFromRelations(ctx, accountID, groupRelationIDs...)
		})
		return err
	})
}

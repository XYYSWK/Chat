package tx

import (
	db "Chat/dao/postgresql/sqlc"
	"Chat/dao/redis/operate"
	"Chat/errcodes"
	"Chat/model"
	"Chat/pkg/tool"
	"context"
	"database/sql"
	"github.com/jackc/pgtype"
)

// CreateApplicationTx 使用事务先判断是否存在申请，不存在则创建申请
func (store *SqlStore) CreateApplicationTx(ctx context.Context, params *db.CreateApplicationParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		// 查看申请是否存在
		ok, err := queries.ExistsApplicationByIDWithLock(ctx, &db.ExistsApplicationByIDWithLockParams{
			Account1ID: params.Account1ID,
			Account2ID: params.Account2ID,
		})
		if err != nil {
			return err
		}
		if ok {
			return errcodes.ApplicationExists
		}
		// 创建申请
		return queries.CreateApplication(ctx, params)
	})
}

// AcceptApplicationTx 接收申请并建立好友关系和双方关系设置并添加到 redis 中
func (store *SqlStore) AcceptApplicationTx(ctx context.Context, rdb *operate.RDB, account1, account2 *db.GetAccountByIDRow) (*db.Message, error) {
	var result *db.Message
	err := store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		// 修改申请状态
		err = tool.DoThat(err, func() error {
			return queries.UpdateApplication(ctx, &db.UpdateApplicationParams{
				Status:     db.ApplicationstatusValue1,
				Account1ID: account1.ID,
				Account2ID: account2.ID,
			})
		})
		id1, id2 := account1.ID, account2.ID
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		// 建立好友关系
		var relationID int64
		err = tool.DoThat(err, func() error {
			relationID, err = queries.CreateFriendRelation(ctx, &db.CreateFriendRelationParams{
				Account1ID: id1,
				Account2ID: id2,
			})
			return err
		})
		// 建立双方的关系设置
		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  account1.ID,
				RelationID: relationID,
				IsLeader:   false,
				IsSelf:     false,
			})
		})
		err = tool.DoThat(err, func() error {
			return queries.CreateSetting(ctx, &db.CreateSettingParams{
				AccountID:  account2.ID,
				RelationID: relationID,
				IsLeader:   false,
				IsSelf:     false,
			})
		})
		// 新建一个系统通知消息作为好友的第一条消息
		err = tool.DoThat(err, func() error {
			arg := &db.CreateMessageParams{
				NotifyType: db.MsgnotifytypeCommon,
				MsgType:    string(model.MsgTypeText),
				MsgContent: "我们已经成为好友啦，现在可以开始聊天啦！",
				MsgExtend:  pgtype.JSON{Status: pgtype.Null},
				AccountID:  sql.NullInt64{Int64: account2.ID, Valid: true},
				RelationID: relationID,
			}
			msgInfo, err := queries.CreateMessage(ctx, arg)
			result = &db.Message{
				ID:         msgInfo.ID,
				NotifyType: arg.NotifyType,
				MsgType:    arg.MsgType,
				MsgContent: arg.MsgContent,
				RelationID: relationID,
				CreateAt:   msgInfo.CreateAt,
			}
			return err
		})
		// 添加关系到 redis
		err = tool.DoThat(err, func() error {
			return rdb.AddRelationAccount(ctx, relationID, account1.ID, account2.ID)
		})
		return err
	})
	return result, err
}

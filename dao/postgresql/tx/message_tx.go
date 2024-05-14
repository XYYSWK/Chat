package tx

import (
	db "Chat/dao/postgresql/sqlc"
	"Chat/pkg/tool"
	"context"
)

// RevokeMsgWithTx 撤回消息，如果消息 pin 或者置顶，则全部取消
func (store *SqlStore) RevokeMsgWithTx(ctx context.Context, msgID int64, isPin, isTop bool) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		err = tool.DoThat(err, func() error {
			return queries.UpdateMsgRevoke(ctx, &db.UpdateMsgRevokeParams{
				ID:       msgID,
				IsRevoke: true,
			})
		})
		if isPin {
			err = tool.DoThat(err, func() error {
				return queries.UpdateMsgPin(ctx, &db.UpdateMsgPinParams{
					ID:    msgID,
					IsPin: false,
				})
			})
		}
		if isTop {
			err = tool.DoThat(err, func() error {
				return queries.UpdateMsgTop(ctx, &db.UpdateMsgTopParams{
					ID:    msgID,
					IsTop: false,
				})
			})
		}
		return err
	})
}

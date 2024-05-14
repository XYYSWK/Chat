package tx

import (
	db "Chat/dao/postgresql/sqlc"
	"context"
	"database/sql"
)

// UploadGroupAvatarWithTx 创建群组头像文件
func (store *SqlStore) UploadGroupAvatarWithTx(ctx context.Context, arg db.CreateFileParams) error {
	return store.execTx(ctx, func(queries *db.Queries) error {
		var err error
		_, err = queries.GetGroupAvatar(ctx, arg.RelationID)
		if err != nil {
			// 如果没有设置过群头像
			if err.Error() == "no rows in result set" {
				_, err = queries.CreateFile(ctx, &db.CreateFileParams{
					FileName:   arg.FileName,
					FileType:   "img",
					FileSize:   arg.FileSize,
					Key:        arg.Key,
					Url:        arg.Url,
					RelationID: arg.RelationID,
					AccountID:  sql.NullInt64{},
				})
			} else {
				return err
			}
		} else {
			// 在 file 表中覆盖之前的群头像
			err = queries.UpdateGroupAvatar(ctx, &db.UpdateGroupAvatarParams{
				Url:        arg.Url,
				RelationID: arg.RelationID,
			})
		}
		data, err := queries.GetGroupRelationByID(ctx, arg.RelationID.Int64)
		if err != nil {
			return err
		}
		// 更新 relation 表中的头像数据
		return queries.UpdateGroupRelation(ctx, &db.UpdateGroupRelationParams{
			Name:        data.Name,
			Description: data.Description,
			Avatar:      arg.Url,
			ID:          arg.RelationID.Int64,
		})
	})
}

package request

import "Chat/model"

type ParamCreateNotify struct {
	RelationID int64            `json:"relation_id" form:"relation_id" binding:"required"`
	MsgContent string           `json:"msg_content" form:"msg_content" binding:"required"`
	MsgExtend  *model.MsgExtend `json:"msg_extend" form:"msg_extend"`
}

type ParamUpdateNotify struct {
	ID         int64            `json:"id" form:"id" binding:"required"` // 群通知 ID
	RelationID int64            `json:"relation_id" form:"relation_id" binding:"required"`
	MsgContent string           `json:"msg_content" form:"msg_content" binding:"required"`
	MsgExtend  *model.MsgExtend `json:"msg_extend" form:"msg_extend"`
}

type ParamGetNotifyByID struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required"`
}

type ParamDeleteNotify struct {
	ID         int64 `json:"id" form:"id" binding:"required"`                   // 群通知的 ID
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required"` // 群 ID
}

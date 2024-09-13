package request

import (
	"Chat/model/common"
	"mime/multipart"
)

type ParamCreateFileMsg struct {
	RelationID int64                 `form:"relation_id"` // 关系 ID
	File       *multipart.FileHeader `form:"file"`        // 文件
	RlyMsg     int64                 `form:"rly_msg"`     // 回复消息 ID
}

type ParamGetMsgsByRelationIDAndTime struct {
	RelationID int64 `form:"relation_id" binding:"required,gte=1"` // 关系 ID
	LastTime   int32 `form:"last_time" binding:"required,gte=1"`   // 拉取消息最晚的时间戳（精确到秒）
	common.Page
}

type ParamOfferMsgsByAccountIDAndTime struct {
	LastTime int32 `form:"last_time" binding:"required,gte=1"` // 拉取消息的最晚时间戳（精确到秒）
	common.Page
}

type ParamUpdateMsgPin struct {
	ID         int64 `json:"id" binding:"required,gte=1"`          // 消息 ID
	RelationID int64 `json:"relation_id" binding:"required,gte=1"` // 关系 ID
	IsPin      bool  `json:"is_pin" binding:"required"`            // 是否 pin
}

type ParamUpdateMsgTop struct {
	ID         int64 `json:"id" binding:"required,gte=1"`          // 消息 ID
	RelationID int64 `json:"relation_id" binding:"required,gte=1"` // 关系 ID
	IsTop      bool  `json:"is_top" binding:"required"`            // 是否置顶
}

type ParamRevokeMsg struct {
	ID int64 `json:"id" binding:"required,gte=1"` // 消息 ID
}

type ParamGetTopMsgByRelationID struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required,gte=1"`
}

type ParamGetPinMsgsByRelationID struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required,gte=1"`
	common.Page
}

type ParamGetRlyMsgsInfoByMsgID struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required,gte=1"`
	MsgID      int64 `json:"msg_id" form:"msg_id" binding:"required,gte=1"`
	common.Page
}

type ParamGetMsgsByContent struct {
	RelationID int64  `form:"relation_id"`
	Content    string `form:"content" binding:"required"`
	common.Page
}

package client

import (
	"Chat/model"
	"time"
)

/*
chat 中 client 端有关消息请求的结构
*/

type HandleSendMsgParams struct {
	RelationID int64            `json:"relation_id" validate:"required,gte=1"`          // 关系 ID
	MsgContent string           `json:"msg_content" validate:"required,gte=1,lte=1000"` // 消息内容
	MsgExtend  *model.MsgExtend `json:"msg_extend"`                                     // 消息扩展信息
	RlyMsgID   int64            `json:"rly_msg_id"`                                     // 回复消息 ID（如果是回复消息，则该字段大于 0）
}

type HandleSendMsgRly struct {
	MsgID    int64     `json:"msg_id"`    // 消息 ID
	CreateAt time.Time `json:"create_at"` // 创建时间
}

type HandleReadMsgParams struct {
	MsgIDs     []int64 `json:"msg_ids,omitempty" validate:"required,gte=1,lte=20"` // 消息 ID
	RelationID int64   `json:"relation_id,omitempty" validate:"required,gte=1"`    // 这些消息所属的关系 ID
}

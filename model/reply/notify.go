package reply

import (
	"Chat/model"
	"time"
)

type ParamGroupNotify struct {
	ID         int64            `json:"id"`
	RelationID int64            `json:"relation_id"`
	MsgContent string           `json:"msg_content"`
	MsgExtend  *model.MsgExtend `json:"msg_extend"`
	AccountID  int64            `json:"account_id"`
	CreateAt   time.Time        `json:"create_at"`
	ReadIDs    []int64          `json:"read_ids"`
}

type ParamGetNotifyByID struct {
	List  []ParamGroupNotify `json:"list"`
	Total int64              `json:"total"`
}

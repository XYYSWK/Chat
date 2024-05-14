package server

import "Chat/model"

type CreateNotify struct {
	EnToken    string           `json:"en_token"`
	AccountID  int64            `json:"account_id"`
	RelationID int64            `json:"relation_id"`
	MsgContent string           `json:"msg_content"`
	MsgExtend  *model.MsgExtend `json:"msg_extend"`
}

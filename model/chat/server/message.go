package server

import "Chat/model/reply"

/*
chat 中 server 端有关消息请求的结构
*/

type SendMsg struct {
	EnToken string                    `json:"en_token,omitempty"` // 加密后的 token
	MsgInfo reply.ParamMsgInfoWithRly `json:"msg_info"`
}

type ReadMsg struct {
	EnToken  string  `json:"en_token,omitempty"` // 加密后的 Token
	MsgIDs   []int64 `json:"msg_ids"`            // 已读消息 IDs
	ReaderID int64   `json:"reader_id"`          // 读者账号 ID
}

type MsgType string

const (
	MsgPin    MsgType = "pin"    // pin 消息
	MsgTop    MsgType = "top"    // 置顶消息
	MsgRevoke MsgType = "revoke" // 撤回消息
)

type UpdateMsgState struct {
	EnToken string  `json:"en_token,omitempty"` // 加密后的 token
	MsgType MsgType `json:"msg_type,omitempty"` // 消息类型 [pin, top, revoke]
	MsgID   int64   `json:"msg_id,omitempty"`   // 消息 ID
	State   bool    `json:"state"`              // 状态设置
}

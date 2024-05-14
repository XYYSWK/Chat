package reply

import (
	"Chat/model"
	"time"
)

type ParamCreateFileMsg struct {
	ID         int64     `json:"id,omitempty"`          // 消息 ID
	MsgContent string    `json:"msg_content,omitempty"` // 消息内容，文件则为 url
	FileID     int64     `json:"file_id,omitempty"`     // 文件 ID
	CreateAt   time.Time `json:"create_at"`             // 创建时间
}

type ParamRlyMsg struct {
	MsgID      int64            `json:"msg_id,omitempty"`      // 回复消息 ID
	MsgType    string           `json:"msg_type,omitempty"`    // 消息类型[text,file]
	MsgContent string           `json:"msg_content,omitempty"` // 消息内容
	MsgExtend  *model.MsgExtend `json:"msg_extend,omitempty"`  // 消息扩展信息(被@的信息)，可能为 null
	IsRevoked  bool             `json:"is_revoked,omitempty"`  // 是否撤回
}

type ParamMsgInfo struct {
	ID         int64            `json:"id,omitempty"`          // 消息 ID
	NotifyType string           `json:"notify_type,omitempty"` // 通知类型 [system,common]
	MsgType    string           `json:"msg_type,omitempty"`    // 消息类型 [text, file]
	MsgContent string           `json:"msg_content,omitempty"` // 消息内容 文件则为 url，文本则为文本内容，由拓展信息进行补充
	MsgExtend  *model.MsgExtend `json:"msg_extend"`            // 消息扩展信息，可能为 null
	FileID     int64            `json:"file_id"`               // 文件 ID，当消息类型为 file 时 > 0
	AccountID  int64            `json:"account_id"`            // 账号 ID，发送者的 ID
	RelationID int64            `json:"relation_id,omitempty"` // 关系 ID
	CreateAt   time.Time        `json:"create_at"`             // 创建时间
	IsRevoke   bool             `json:"is_revoke,omitempty"`   // 是否撤回
	IsTop      bool             `json:"is_top,omitempty"`      // 是否置顶
	IsPin      bool             `json:"is_pin,omitempty"`      // 是否 pin
	PinTime    time.Time        `json:"pin_time"`              // pin 时间
	ReadIds    []int64          `json:"read_ids,omitempty"`    // 已读的账号ID，当请求者为发送者时为空
	ReplyCount int64            `json:"reply_count,omitempty"` // 回复数
}

// ParamMsgInfoWithRly 完整的消息详情，包含回复消息
type ParamMsgInfoWithRly struct {
	ParamMsgInfo
	RlyMsg *ParamRlyMsg `json:"rly_msg"` // 回复消息详情，可能为 nil
}

// ParamMsgInfoWithRlyAndHasRead 完整的消息详情，包含回复消息，包含是否已读
type ParamMsgInfoWithRlyAndHasRead struct {
	ParamMsgInfoWithRly
	HasRead bool `json:"has_read"` // 是否已读
}

type ParamGetMsgsRelationIDAndTime struct {
	List  []*ParamMsgInfoWithRly `json:"list"`
	Total int64                  `json:"total"`
}

type ParamOfferMsgsByAccountIDAndTime struct {
	List  []*ParamMsgInfoWithRlyAndHasRead `json:"list,omitempty"`
	Total int64                            `json:"total,omitempty"`
}

type ParamGetTopMsgByRelationID struct {
	MsgInfo ParamMsgInfo `json:"msg_info"` // 置顶消息详情
}

type ParamGetPinMsgsByRelationID struct {
	List  []*ParamMsgInfo `json:"list"`
	Total int64           `json:"total"`
}

type ParamGetRlyMsgsInfoByMsgID struct {
	List  []*ParamMsgInfo `json:"list"`
	Total int64           `json:"total"`
}

type ParamGetMsgsByContent struct {
	List  []*ParamBriefMsgInfo `json:"list"`
	Total int64                `json:"total"`
}

type ParamBriefMsgInfo struct {
	ID         int64            `json:"id"`          // 消息ID
	NotifyType string           `json:"notify_type"` // 通知类型 [system,common]
	MsgType    string           `json:"msg_type"`    // 消息类型 [text,file]
	MsgContent string           `json:"msg_content"` // 消息内容 文件则为url，文本则为文本内容，由拓展信息进行补充
	Extend     *model.MsgExtend `json:"msg_extend"`  // 消息扩展信息 可能为null
	FileID     int64            `json:"file_id"`     // 文件ID 当消息类型为file时>0
	AccountID  int64            `json:"account_id"`  // 账号ID 发送者ID
	RelationID int64            `json:"relation_id"` // 关系ID
	CreateAt   time.Time        `json:"create_at"`   // 创建时间
}

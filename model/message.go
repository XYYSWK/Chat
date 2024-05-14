package model

import (
	"github.com/jackc/pgtype"
	"mime/multipart"
	"time"
)

type MsgType string

const (
	MsgTypeText MsgType = "text"
	MsgTypeFile MsgType = "file"
)

type CreateFileMsg struct {
	AccountID  int64
	RelationID int64
	RlyMsgID   int64
	File       *multipart.FileHeader
}

type Remind struct {
	Idx       int64 `json:"idx,omitempty" binding:"required,gte=1" validate:"required,gte=1"`        // 第几个 @
	AccountID int64 `json:"account_id,omitempty" binding:"required,gte=1" validate:"required,gte=1"` // 被 @ 的账号 ID
}

type MsgExtend struct {
	Remind []Remind `json:"remind"` // @ 的描述信息
}

type GetMsgsByRelationIDAndTime struct {
	AccountID  int64
	RelationID int64
	LastTime   time.Time
	Limit      int32
	Offset     int32
}

type OfferMsgsByAccountIDAndTime struct {
	AccountID int64
	LastTime  time.Time
	Limit     int32
	Offset    int32
}

// ExtendToJson 将 MsgExtend 转化为 pgtype.Json，可以是 nil
// 参数：消息扩展信息
// 返回：pgtype.Json 对象
func ExtendToJson(extend *MsgExtend) (pgtype.JSON, error) {
	data := pgtype.JSON{}
	err := data.Set(extend)
	return data, err
}

// JsonToExtend 将 pgtype.Json 转化为 MsgExtend
// 参数：pgtype.Json 对象（如果存的 json 为 nil 或未定义则返回 nil）
// 返回：解析后的消息扩展信息（可能为 nil）
func JsonToExtend(data pgtype.JSON) (*MsgExtend, error) {
	if data.Status != pgtype.Present { // 判断 data 是否是有效值，即是否为 nil
		return nil, nil
	}
	extend := &MsgExtend{}
	err := data.AssignTo(extend)
	return extend, err
}

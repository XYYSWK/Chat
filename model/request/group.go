package request

import (
	"Chat/model/common"
	"mime/multipart"
)

type ParamCreateGroup struct {
	Name        string `json:"name" form:"name" binding:"required"` // 群名称
	Description string `json:"description" form:"description"`      // 群描述
	//Avatar      *multipart.FileHeader `form:"file" binding:"required"`                           // 头像文件
}

type ParamTransferGroup struct {
	RelationID  int64 `json:"relation_id" form:"relation_id" binding:"required"`
	ToAccountID int64 `json:"to_account_id" form:"to_account_id" binding:"required"` // 群要转让给的账户ID
}

type ParamDissolveGroup struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required"` // 要解散的 relationID
}

type ParamUpdateGroup struct {
	RelationID  int64                 `json:"relation_id" form:"relation_id" binding:"required"`
	Name        string                `json:"name" form:"name"`
	Description string                `json:"description" form:"description"`
	Avatar      *multipart.FileHeader `form:"avatar"`
}

type ParamInviteAccount struct {
	AccountID  []int64 `json:"account_id" form:"account_id" binding:"required"`   // 要邀请的群成员
	RelationID int64   `json:"relation_id" form:"relation_id" binding:"required"` // 群聊 ID
}

type ParamQuitGroup struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required"`
}

type ParamGetGroupsByName struct {
	Name string `json:"name" form:"name" binding:"required"`
	common.Page
}

type ParamGetGroupMembers struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required"`
	common.Page
}

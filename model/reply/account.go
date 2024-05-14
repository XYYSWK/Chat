package reply

import (
	"Chat/model/common"
	"time"
)

/*
定义账户相关的响应参数结构体
*/

type ParamCreateAccount struct {
	ParamAccountInfo     ParamAccountInfo     `json:"param_account_info"`      // 账号信息
	ParamGetAccountToken ParamGetAccountToken `json:"param_get_account_token"` // 账号 token
}

type ParamAccountInfo struct {
	ID     int64  `json:"id,omitempty"`     // 账号 ID
	Name   string `json:"name,omitempty"`   // 昵称
	Avatar string `json:"avatar,omitempty"` // 头像
	Gender string `json:"gender,omitempty"` // 性别 [男，女，未知]
}

type ParamGetAccountToken struct {
	AccountToken common.Token `json:"account_token"` // 账号 token
}

type ParamGetAccountsByUserID struct {
	List  []ParamAccountInfo `json:"list,omitempty"`  // 账号列表
	Total int64              `json:"total,omitempty"` // 总数
}

type ParamFriendInfo struct {
	ParamAccountInfo
	RelationID int64 `json:"relation_id,omitempty"` // 好友关系 ID，0表示没有好友关系
}

type ParamGetAccountsByName struct {
	List  []*ParamFriendInfo `json:"list,omitempty"`  // 账号列表
	Total int64              `json:"total,omitempty"` // 总数
}

type ParamGetAccountByID struct {
	Info       ParamAccountInfo `json:"info"`                // 账号信息
	Signature  string           `json:"signature,omitempty"` // 个性签名
	CreateAt   time.Time        `json:"create_at,omitempty"` // 创建时间
	RelationID int64            `json:"relation_id"`         // 关系ID，如果不存在则为 0
}

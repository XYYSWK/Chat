package request

import "Chat/model/common"

/*
定义账户相关的请求参数结构体
*/

type ParamCreateAccount struct {
	Name      string `json:"name" binding:"required,gte=1,lte=20"`       // 账户名（唯一）
	Gender    string `json:"gender" binding:"required,oneof= 男 女 未知"`    // 性别
	Signature string `json:"signature" binding:"required,gte=0,lte=100"` // 签名
}

type ParamGetAccountToken struct {
	AccountID int64 `json:"account_id" form:"account_id" binding:"required,gte=1"` // 账号 ID
}

type ParamDeleteAccount struct {
	AccountID int64 `json:"account_id" binding:"required,gte=1"` // 账号 ID
}

type ParamUpdateAccount struct {
	Name      string `json:"name" binding:"required,gte=1,lte=20"`       // 名称
	Gender    string `json:"gender" binding:"required,oneof= 男 女 未知"`    // 性别
	Signature string `json:"signature" binding:"required,gte=0,lte=100"` // 个性签名
}

type ParamGetAccountsByName struct {
	Name        string `json:"name" form:"name" binding:"required,gte=1,lte=20"` // 搜索名称
	common.Page        // 分页
}

type ParamGetAccountByID struct {
	AccountID int64 `json:"account_id" form:"account_id" binding:"required,gte=1"` // 账号ID
}

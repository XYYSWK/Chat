package reply

import (
	"github.com/XYYSWK/Lutils/pkg/token"
	"time"
)

/*
定义用户相关的响应参数结构体
*/

type ParamUserInfo struct {
	ID       int64     `json:"id,omitempty"`    // user id
	Email    string    `json:"email,omitempty"` // 邮箱
	CreateAt time.Time `json:"create_at"`       // 创建时间
}

type ParamRegister struct {
	ParamUserInfo ParamUserInfo `json:"param_user_info"` // 用户信息
	Token         ParamToken    `json:"token"`           // 用户令牌
}

type ParamLogin struct {
	ParamUserInfo ParamUserInfo `json:"param_user_info"` // 用户信息
	Token         ParamToken    `json:"token"`           // 用户令牌
}

type ParamToken struct {
	AccessToken   string         `json:"access_token"`
	AccessPayload *token.Payload `json:"access_payload"`
	RefreshToken  string         `json:"refresh_token"`
}

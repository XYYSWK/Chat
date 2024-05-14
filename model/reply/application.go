package reply

import "time"

type ParamApplicationInfo struct {
	AccountID1 int64     `json:"account_id_1,omitempty"` // 申请者账号 ID
	AccountID2 int64     `json:"account_id_2,omitempty"` // 被申请者账号 ID
	ApplyMsg   string    `json:"apply_msg,omitempty"`    // 申请信息
	Refuse     string    `json:"refuse,omitempty"`       // 拒绝信息
	Status     string    `json:"status,omitempty"`       // 状态 [已申请，已拒绝，已同意]
	CreateAt   time.Time `json:"create_at"`              // 创建时间
	UpdateAt   time.Time `json:"update_at"`              // 更新时间
	Name       string    `json:"name,omitempty"`         // 对方账号的名字
	Avatar     string    `json:"avatar,omitempty"`       // 对方账号的头像
}

type ParamListApplications struct {
	List  []*ParamApplicationInfo `json:"list,omitempty"`
	Total int64                   `json:"total,omitempty"`
}

package request

type ParamCreateApplication struct {
	AccountID      int64  `json:"account_id" binding:"required,gte=1"` // 被申请者的 ID
	ApplicationMsg string `json:"application_msg" binding:"lte=200"`   // 申请信息
}

type ParamDeleteApplication struct {
	AccountID int64 `json:"account_id" binding:"required,gte=1"` // 被申请者的 ID
}

type ParamRefuseApplication struct {
	AccountID int64  `json:"account_id" binding:"required,gte=1"` // 申请者的 ID
	RefuseMsg string `json:"refuse_msg" binding:"lte=200"`        // 拒绝信息
}

type ParamAcceptApplication struct {
	AccountID int64 `json:"account_id" binding:"required,gte=1"` // 申请者的 ID
}

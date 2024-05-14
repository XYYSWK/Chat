package request

/*
定义邮箱相关的请求参数结构体
*/

type ParamExistEmail struct {
	Email string `form:"email" binding:"required, email,lte=50"` // 邮箱
}

type ParamSendEmail struct {
	Email string `json:"email" binding:"required,lte=50"` // 邮箱
}

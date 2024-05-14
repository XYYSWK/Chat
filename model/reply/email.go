package reply

/*
定义邮件相关的响应参数结构体
*/

type ParamExistEmail struct {
	Exist bool `json:"exist,omitempty"` // 是否存在
}

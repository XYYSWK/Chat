package common

import (
	"Chat/global"
	"encoding/json"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/go-playground/validator/v10"
	"sync"
	"time"
)

/*
通用的结构体
*/

// Page 分页
type Page struct {
	Page     int32 `json:"page" form:"page"`           // 第几页
	PageSize int32 `json:"page_size" form:"page_size"` // 每页的大小
}

// State 自制的标准回复格式的结构体
type State struct {
	Code int         `json:"code,omitempty"` // 状态码，0：成功；否则失败
	Msg  string      `json:"msg,omitempty"`  // 状态的具体描述
	Data interface{} `json:"data,omitempty"` // 数据，失败时返回空
}

// NewState 创建一个自己的标准恢复格式
// 参数：err 错误信息（可为 nil）；datas 数据（存在只选择第一个值）
// 返回：自制的标准回复格式的结构体
func NewState(err errcode.Err, datas ...interface{}) *State {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	if err == nil {
		err = errcode.StatusOk
	} else {
		data = nil
	}
	return &State{
		Code: err.ECode(),
		Msg:  err.Error(),
		Data: data,
	}
}

// Json 将 state 结构体转换为 json 格式的数据
func (s *State) Json() ([]byte, error) {
	return json.Marshal(s)
}

// MustJson 将 state 结构体转换为 json 格式的数据，如果出错，则返回
func (s *State) MustJson() string {
	v, err := s.Json()
	if err != nil {
		global.Logger.Logger.Error(err.Error())
		return ""
	}
	return string(v)
}

// List 列表
type List struct {
	List  interface{} `json:"list,omitempty"`  // 列表数据，可以是任意类型的集合
	Total int64       `json:"total,omitempty"` // 列表中数据项的总数
}

// Token token
type Token struct {
	Token    string    `json:"token,omitempty"` // token
	ExpireAt time.Time `json:"expire_at"`       // token 过期时间
}

var validate *validator.Validate
var validateOnce sync.Once

// Decode 将 json 格式的数据解析到结构体，并进行校验
// 参数：data：json 格式的数据；v：将要转化为的结构体
// 返回：错误信息
func Decode(data string, v interface{}) error {
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return err
	}
	// 确保只创建一次 validator 实例
	validateOnce.Do(func() {
		validate = validator.New()
	})
	// 结构体校验
	return validate.Struct(v)
}

package server

/*
chat 中 server 端有关处理结构
*/

type AccountLogin struct {
	EnToken string `json:"enToken"` // 加密之后的 token
	Address string `json:"address"` // 地址
}

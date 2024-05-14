package server

type UpdateEmail struct {
	EnToken string `json:"en_token,omitempty"`
	Email   string `json:"email,omitempty"` // 更新后的邮箱
}

type UpdateAccount struct {
	EnToken   string `json:"en_token,omitempty"`  // 加密后的 token
	Name      string `json:"name,omitempty"`      // 昵称
	Gender    string `json:"gender,omitempty"`    // 性别
	Signature string `json:"signature,omitempty"` // 签名
}

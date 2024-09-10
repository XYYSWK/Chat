package model

import (
	"encoding/json"
	"github.com/XYYSWK/Lutils/pkg/token"
)

type TokenType string

type Content struct {
	TokenType TokenType `json:"token_type,omitempty"` // token 类型，用户token\账户token
	ID        int64     `json:"id,omitempty"`
}

// Token 结合 Token、token.Payload 和 Content
type Token struct {
	AccessToken string
	Payload     *token.Payload
	Content     *Content
}

const (
	UserToken    TokenType = "user"
	AccountToken TokenType = "account"
)

// NewTokenContent 新建一种类型的 token
func NewTokenContent(t TokenType, ID int64) *Content {
	return &Content{
		TokenType: t,
		ID:        ID,
	}
}

// Marshal 将 Content 结构体序列化为 json 序列
func (c *Content) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

// Unmarshal 将 json 序列解析为 Content 结构体
func (c *Content) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	return nil
}

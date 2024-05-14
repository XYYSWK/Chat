package chat

import (
	"Chat/dao"
	"Chat/errcodes"
	"Chat/global"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/chat"
	"Chat/model/common"
	"context"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
	"time"
)

// CheckAuth 检查 token 是否有效，有效返回 token，否则断开连接
func CheckAuth(s socketio.Conn) (*model.Token, bool) {
	token, myErr := CheckConnCtxToken(s.Context())
	if myErr != nil {
		s.Emit(chat.ServerError, common.NewState(myErr))
		_ = s.Close()
		return nil, false
	}
	return token, true
}

// CheckConnCtxToken 检查连接上下文中的 token 是否有效，有效则返回 token
// 参数：连接上下文
// 成功：上下文中包含 *model.Token 且有效
// 失败：返回 errcodes.AuthenticationFailed,errcodes.AuthOverTime
func CheckConnCtxToken(v interface{}) (*model.Token, errcode.Err) {
	token, ok := v.(*model.Token)
	if !ok {
		return nil, errcodes.AuthenticationFailed
	}
	if token.Payload.ExpiredAt.Before(time.Now()) {
		return nil, errcodes.AuthOverTime
	}
	return token, nil
}

// MustAccount 解析 token 并判断是否是账户，返回 token
// 参数：accessToken
// 成功：解析 token 的 content 并进行校验返回 *model.Token, nil
// 失败：返回 errcodes.AuthenticationFailed，errcodes.AccountNotFound，errcode.ErrServer
func MustAccount(accessToken string) (*model.Token, errcode.Err) {
	payload, _, myErr := middlewares.ParseToken(accessToken)
	if myErr != nil {
		return nil, myErr
	}
	content := new(model.Content)
	if err := content.Unmarshal(payload.Content); err != nil {
		return nil, errcodes.AuthenticationFailed
	}
	if content.TokenType != model.AccountToken {
		return nil, errcodes.AuthenticationFailed
	}
	ok, err := dao.Database.DB.ExistsAccountByID(context.Background(), content.ID)
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	if !ok {
		return nil, errcodes.AccountNotFound
	}
	return &model.Token{
		AccessToken: accessToken,
		Payload:     payload,
		Content:     content,
	}, nil
}

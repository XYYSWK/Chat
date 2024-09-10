package middlewares

import (
	"Chat/dao"
	"Chat/errcodes"
	"Chat/global"
	"Chat/model"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/XYYSWK/Lutils/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/*
用户验证（paseTo 生成 Token）
*/

// GetToken 从当前请求头获取 token
func GetToken(header http.Header) (string, errcode.Err) {
	//本项目 Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
	authorizationHeader := header.Get(global.PrivateSetting.Token.AuthorizationKey)
	if len(authorizationHeader) == 0 {
		return "", errcodes.AuthNotExist
	}
	//按空格切割（切割为：Bearer 和 token）
	parts := strings.SplitN(authorizationHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == global.PrivateSetting.Token.AuthorizationType) {
		return "", errcodes.AuthenticationFailed
	}
	return parts[1], nil
}

// ParseToken 解析 header 中的 token。返回 payload，token，err
func ParseToken(accessToken string) (*token.Payload, string, errcode.Err) {
	//解析 token
	payload, err := global.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		if err.Error() == "超时错误" {
			return nil, "", errcodes.AuthOverTime
		}
		return nil, "", errcodes.AuthenticationFailed
	}
	return payload, accessToken, nil
}

// PasetoAuth 鉴权中间件，用于解析并写入 Token
func PasetoAuth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		accessToken, err := GetToken(ctx.Request.Header)
		if err != nil {
			ctx.Next()
			return
		}
		payload, _, err := ParseToken(accessToken)
		if err != nil {
			ctx.Next()
			return
		}
		content := &model.Content{}
		if err := content.Unmarshal(payload.Content); err != nil {
			ctx.Next()
			return
		}
		//将当前请求头中的 Content（token 类型和 id）信息保存到请求的上下文 ctx 上
		ctx.Set(global.PrivateSetting.Token.AuthorizationKey, content)
		ctx.Next() //后续的处理请求的函数可以通过 ctx.Get(global.PrivateSetting.Token.AuthorizationKey) 来获取当前请求的用户信息
	}
}

// MustUser 必须是用户
func MustUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reply := app.NewResponse(ctx)
		val, ok := ctx.Get(global.PrivateSetting.Token.AuthorizationKey)
		if !ok {
			reply.Reply(errcodes.AuthNotExist)
			ctx.Abort()
			return
		}
		data := val.(*model.Content)
		if data.TokenType != model.UserToken {
			reply.Reply(errcodes.AuthenticationFailed)
			ctx.Abort()
			return
		}
		ok, err := dao.Database.DB.ExistsUserByID(ctx, data.ID)
		if err != nil {
			global.Logger.Error(err.Error(), ErrLogMsg(ctx)...)
			reply.Reply(errcode.ErrServer)
			ctx.Abort()
			return
		}

		if !ok {
			reply.Reply(errcodes.UserNotFound)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// MustAccount 必须是账号
func MustAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reply := app.NewResponse(ctx)
		val, ok := ctx.Get(global.PrivateSetting.Token.AuthorizationKey)
		if !ok {
			reply.Reply(errcodes.AuthNotExist)
			ctx.Abort()
			return
		}
		data := val.(*model.Content)
		if data.TokenType != model.AccountToken {
			reply.Reply(errcodes.AuthenticationFailed)
			ctx.Abort()
			return
		}
		ok, err := dao.Database.DB.ExistsAccountByID(ctx, data.ID)
		if err != nil {
			global.Logger.Error(err.Error(), ErrLogMsg(ctx)...)
			reply.Reply(errcode.ErrServer)
			ctx.Abort()
			return
		}
		if !ok {
			reply.Reply(errcodes.AccountNotFound)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// GetTokenContent 从当前上下文中获取保存的 Content 内容
func GetTokenContent(ctx *gin.Context) (*model.Content, bool) {
	value, ok := ctx.Get(global.PrivateSetting.Token.AuthorizationKey)
	if !ok {
		return nil, false
	}
	return value.(*model.Content), true
}

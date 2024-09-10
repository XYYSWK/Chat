package api

import (
	"Chat/errcodes"
	"Chat/global"
	"Chat/logic"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/request"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type account struct {
}

// CreateAccount 创建账号
// @Tags account
// @Summary 创建账号
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data body request.ParamCreateAccount true "创建账号信息"
// @Success 200 {object} common.State{data=reply.ParamCreateAccount} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2012:账户名已经存在"
// @Router /api/account/create [post]
func (account) CreateAccount(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateAccount)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.CreateAccount(ctx, content.ID, params.Name, global.PublicSetting.Rules.DefaultAvatarURL, params.Gender, params.Signature)
	reply.Reply(err, result)
}

// GetAccountToken 获取账号令牌
// @Tags account
// @Summary 获取账号令牌
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data query request.ParamGetAccountToken true "账号ID"
// @Success 200 {object} common.State{data=reply.ParamGetAccountToken} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 2010:账号不存在"
// @Router /api/account/token [get]
func (account) GetAccountToken(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountToken)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok && content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.GetAccountToken(ctx, content.ID, params.AccountID)
	reply.Reply(err, result)
}

// DeleteAccount 删除账号
// @Tags account
// @Summary 删除账号
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param data query request.ParamDeleteAccount true "账号ID"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 2010:账号不存在"
// @Router /api/account/delete [delete]
func (account) DeleteAccount(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteAccount)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Account.DeleteAccount(ctx, content.ID, params.AccountID)
	reply.Reply(err)
}

// GetAccountsByUserID 获取用户的所有账号
// @Tags account
// @Summary 获取用户的所有账号
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} common.State{data=reply.ParamGetAccountsByUserID} "1001:参数有误 1003:系统错误 2008:身份验证失败 2010:账号不存在"
// @Router /api/account/infos/account [get]
func (account) GetAccountsByUserID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.UserToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.GetAccountsByUserID(ctx, content.ID)
	reply.Reply(err, result)
}

// UpdateAccount 更改账号信息
// @Tags account
// @Summary 更改账号信息
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账号令牌"
// @Param data body request.ParamUpdateAccount true "账号信息"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败"
// @Router /api/account/infos/account [get]
func (account) UpdateAccount(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateAccount)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Account.UpdateAccount(ctx, content.ID, params.Name, params.Gender, params.Signature)
	reply.Reply(err)
}

// GetAccountsByName 根据昵称模糊查找账号
// @Tags account
// @Summary 根据昵称模糊查找账号
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账号令牌"
// @Param data query request.ParamGetAccountsByName true "账号信息"
// @Success 200 {object} common.State{data=reply.ParamGetAccountsByName} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在"
// @Router /api/account/infos/name [get]
func (account) GetAccountsByName(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountsByName)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Account.GetAccountsByName(ctx, content.ID, params.Name, limit, offset)
	reply.Reply(err, result)
}

// GetAccountByID 根据ID查找账户
// @Tags account
// @Summary 根据ID查找账户
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账号令牌"
// @Param data query request.ParamGetAccountByID true "要查询的账号ID"
// @Success 200 {object} common.State{data=reply.ParamGetAccountByID} "1001:参数有误 1003:系统错误 2009:权限不足 2007:身份不存在 2008:身份验证失败 2010:账号不存在"
// @Router /api/account/info [get]
func (account) GetAccountByID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetAccountByID)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Account.GetAccountByID(ctx, params.AccountID, content.ID)
	reply.Reply(err, result)
}

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

type application struct {
}

// CreateApplication 创建好友申请
// @Tags application
// @Summary 创建好友申请
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Param data body request.ParamCreateApplication true "申请信息"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 3001:申请已经存在 3003:申请不合法 4001:关系已经存在"
// @Router /api/application/create [post]
func (application) CreateApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateApplication)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Application.CreateApplication(ctx, content.ID, params.AccountID, params.ApplicationMsg)
	reply.Reply(err)
}

// DeleteApplication 申请者删除好友申请
// @Tags application
// @Summary 申请者删除好友申请
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Param data body request.ParamDeleteApplication true "需要删除的申请"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 2009:权限不足 3002:申请不存在 3003:申请不合法"
// @Router /api/application/delete [delete]
func (application) DeleteApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteApplication)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Application.DeleteApplication(ctx, content.ID, params.AccountID)
	reply.Reply(err)
}

// RefuseApplication 被申请者拒绝好友申请
// @Tags application
// @Summary 被申请者拒绝好友申请
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Param data body request.ParamRefuseApplication true "需要拒绝的申请"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 3002:申请不存在 3004:重复操作申请"
// @Router /api/application/refuse [put]
func (application) RefuseApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamRefuseApplication)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Application.RefuseApplication(ctx, params.AccountID, content.ID, params.RefuseMsg)
	reply.Reply(err)
}

// AcceptApplication 被申请者接受好友申请
// @Tags application
// @Summary 被申请者接受好友申请
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Param data body request.ParamAcceptApplication true "需要接受的申请"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 3002:申请不存在 3004:重复操作申请"
// @Router /api/application/accept [put]
func (application) AcceptApplication(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamAcceptApplication)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Application.AcceptApplication(ctx, params.AccountID, content.ID)
	reply.Reply(err)
}

// ListApplications 账号查看和自身相关的好友申请（不论是申请者还是被申请者）
// @Tags application
// @Summary 账号查看和自身相关的好友申请（不论是申请者还是被申请者）
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Success 200 {object} common.State{data=reply.ParamListApplications} "1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在"
// @Router /api/application/list [get]
func (application) ListApplications(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Application.ListApplications(ctx, content.ID, limit, offset)
	reply.Reply(err, result)
}

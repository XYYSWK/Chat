package api

import (
	"Chat/errcodes"
	"Chat/logic"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/request"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type notify struct {
}

// CreateNotify 创建群通知
// @Tags     notify
// @Summary  创建群通知
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamCreateNotify                  true  "请求信息"
// @Success  200            {object}  common.State{data=reply.ParamGroupNotify}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 7001:非群主 7003:非群成员"
// @Router   /api/notify/create [post]
func (notify) CreateNotify(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateNotify)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Notify.CreateNotify(ctx, content.ID, params)
	reply.Reply(err, result)
}

// UpdateNotify 更新群通知
// @Tags     notify
// @Summary  更新群通知
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamUpdateNotify                  true  "请求信息"
// @Success  200            {object}  common.State{}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 7001:非群主 7003:非群成员"
// @Router   /api/notify/update [post]
func (notify) UpdateNotify(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateNotify)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Notify.UpdateNotify(ctx, content.ID, params)
	reply.Reply(err)
}

// GetNotifyByID 更具群ID获取群通知
// @Tags     notify
// @Summary  更具群ID获取群通知
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           query     request.ParamGetNotifyByID                 true  "请求信息"
// @Success  200            {object}  common.State{data=reply.ParamGetNotifyByID}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 7003:非群成员"
// @Router   /api/notify/get [get]
func (notify) GetNotifyByID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetNotifyByID)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Notify.GetNotifyByID(ctx, content.ID, params.RelationID)
	reply.ReplyList(err, result.Total, result.List)
}

// DeleteNotify
// @Tags     notify
// @Summary  获取群通知
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           query     request.ParamDeleteNotify                 true  "请求信息"
// @Success  200            {object}  common.State{}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 7001:非群主 7003:非群成员"
// @Router   /api/notify/delete [delete]
func (notify) DeleteNotify(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteNotify)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Notify.DeleteNotify(ctx, content.ID, params.ID, params.RelationID)
	reply.Reply(err)
}

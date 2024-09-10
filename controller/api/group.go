package api

import (
	"Chat/errcodes"
	"Chat/global"
	"Chat/logic"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/reply"
	"Chat/model/request"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type group struct {
}

// CreateGroup 创建群
// @Tags     group
// @Summary  创建群
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamCreateGroup                   true  "请求信息"
// @Success  200            {object}  common.State{data=reply.ParamCreateGroup}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/create [post]
func (group) CreateGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamCreateGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	relationID, err := logic.Logics.Group.CreateGroup(ctx, content.ID, params.Name, params.Description)
	if err != nil {
		rly.Reply(err)
		return
	}
	result, err := logic.Logics.File.UploadGroupAvatar(ctx, nil, content.ID, relationID)
	rly.Reply(err, reply.ParamCreateGroup{
		Name:        params.Name,
		AccountID:   content.ID,
		RelationID:  relationID,
		Description: params.Description,
		Avatar:      result.URL,
	})
}

// TransferGroup 转让群
// @Tags     group
// @Summary  转让群
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamTransferGroup                   true  "请求信息"
// @Success  200            {object}  common.State{}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/transfer [post]
func (group) TransferGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamTransferGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Group.TransferGroup(ctx, content.ID, params.RelationID, params.ToAccountID)
	rly.Reply(err)
}

// DissolveGroup 解散群
// @Tags     group
// @Summary  解散群
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamDissolveGroup                   true  "请求信息"
// @Success  200            {object}  common.State{}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/dissolve [post]
func (group) DissolveGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamDissolveGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Group.DissolveGroup(ctx, content.ID, params.RelationID)
	rly.Reply(err)
}

// UpdateGroup 更新群信息
// @Tags     group
// @Summary  更新群信息
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamUpdateGroup                   true  "请求信息"
// @Success  200            {object}  common.State{data=reply.ParamUpdateGroup}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/update [post]
func (group) UpdateGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamUpdateGroup)
	if err := ctx.ShouldBind(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Group.UpdateGroup(ctx, content.ID, params.RelationID, params.Name, params.Description)
	if err != nil {
		rly.Reply(err, result)
	}
	avatar, err := logic.Logics.File.UploadGroupAvatar(ctx, params.Avatar, content.ID, params.RelationID)
	result.Avatar = avatar.URL
	rly.Reply(err, result)
}

// InviteAccount 邀请成员入群
// @Tags     group
// @Summary  邀请成员入群
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamInviteAccount                   true  "请求信息"
// @Success  200            {object}  common.State{data=reply.ParamInviteAccount}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/invite [post]
func (group) InviteAccount(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamInviteAccount)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Group.InviteAccount(ctx, content.ID, params.RelationID, params.AccountID)
	rly.Reply(err, result)
}

// GetGroupList 获取该账号所有的群聊信息列表
// @Tags     group
// @Summary  获取该账号所有的群聊信息列表
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Success  200            {object}  common.State{data=reply.ParamGetGroupList}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/list [get]
func (group) GetGroupList(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Group.GetGroupList(ctx, content.ID)
	rly.Reply(err, result)
}

// QuitGroup 退出群聊
// @Tags     group
// @Summary  退出群聊
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamQuitGroup                true  "请求信息"
// @Success  200            {object}  common.State{}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/quit [post]
func (group) QuitGroup(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamQuitGroup)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Group.QuitGroup(ctx, content.ID, params.RelationID)
	rly.Reply(err)
}

// GetGroupsByName 通过群名模糊查找群聊
// @Tags     group
// @Summary  通过群名模糊查找群聊
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamGetGroupsByName          true  "请求信息"
// @Success  200            {object}  common.State{data=reply.ParamGetGroupsByName}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 7003:非群员"
// @Router   /api/group/name [get]
func (group) GetGroupsByName(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamGetGroupsByName)
	if err := ctx.ShouldBindJSON(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Group.GetGroupsByName(ctx, content.ID, params.Name, limit, offset)
	rly.Reply(err, result)
}

// GetGroupMembers 查看所有群成员
// @Tags     group
// @Summary  查看所有群成员
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                                true  "Bearer 账户令牌"
// @Param    data           body      request.ParamGetGroupMembers          true  "请求信息"
// @Success  200            {object}  common.State{data=reply.ParamGetGroupMembers}  "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router   /api/group/members [get]
func (group) GetGroupMembers(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	params := new(request.ParamGetGroupMembers)
	if err := ctx.ShouldBindQuery(params); err != nil {
		rly.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		rly.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Group.GetGroupMembers(ctx, content.ID, params.RelationID, limit, offset)
	rly.Reply(err, result)
}

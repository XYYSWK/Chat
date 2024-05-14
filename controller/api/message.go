package api

import (
	"Chat/errcodes"
	"Chat/global"
	"Chat/logic"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/request"
	"github.com/XYYSWK/Rutils/pkg/app"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"time"
)

type message struct {
}

// CreateFileMsg 发布文件类型的消息
// @Tags      message
// @Summary   发布文件类型的消息
// @Security  BasicAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     Authorization  header    string                                  true   "Bearer 账户令牌"
// @Param     file           formData  file                                    true   "文件"
// @Param     relation_id    formData  int64                                   true   "关系id"
// @Param     rly_msg_id     formData  int64                                   false  "回复消息id"
// @Success   200            {object}  common.State{data=reply.ParamCreateFileMsg}  ""
// @Router    /api/message/file [post]
func (message) CreateFileMsg(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamCreateFileMsg)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Message.CreateFileMsg(ctx, model.CreateFileMsg{
		AccountID:  content.ID,
		RelationID: params.RelationID,
		RlyMsgID:   params.RlyMsg,
		File:       params.File,
	})
	reply.Reply(err, result)
}

// GetMsgsByRelationIDAndTime 获取指定关系指定时间戳之前的信息，获取的消息按照发布时间先后排序
// @Tags      message
// @Summary   获取指定关系指定时间戳之前的信息，获取的消息按照发布时间先后排序
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                                  true   "Bearer 账户令牌"
// @Param     data           query     request.ParamGetMsgsByRelationIDAndTime   true   "请求信息"
// @Success   200            {object}  common.State{data=reply.ParamGetMsgsRelationIDAndTime}  "完整的消息详情，包含回复消息"
// @Router    /api/message/list/time [get]
func (message) GetMsgsByRelationIDAndTime(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetMsgsByRelationIDAndTime)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.GetMsgsByRelationIDAndTime(ctx, model.GetMsgsByRelationIDAndTime{
		AccountID:  content.ID,
		RelationID: params.RelationID,
		LastTime:   time.Unix(int64(params.LastTime), 0),
		Limit:      limit,
		Offset:     offset,
	})
	reply.Reply(err, result)
}

// OfferMsgsByAccountIDAndTime 获取所有关系指定时间戳之后的信息，获取的消息按照发布时间先后排序，同时包含是否已读的标识
// @Tags      message
// @Summary   获取所有关系指定时间戳之后的信息，获取的消息按照发布时间先后排序，同时包含是否已读的标识
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                                  true   "Bearer 账户令牌"
// @Param     data           query     request.ParamOfferMsgsByAccountIDAndTime   true   "请求信息"
// @Success   200            {object}  common.State{data=reply.ParamOfferMsgsByAccountIDAndTime}  "完整的消息详情，包含回复消息"
// @Router    /api/message/list/offer [get]
func (message) OfferMsgsByAccountIDAndTime(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamOfferMsgsByAccountIDAndTime)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.OfferMsgsByAccountIDAndTime(ctx, model.OfferMsgsByAccountIDAndTime{
		AccountID: content.ID,
		LastTime:  time.Unix(int64(params.LastTime), 0),
		Limit:     limit,
		Offset:    offset,
	})
	reply.Reply(err, result)
}

// UpdateMsgPin 更改消息的 pin 状态
// @Tags      message
// @Summary   更改消息的 pin 状态
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                      true   "Bearer 账户令牌"
// @Param     data           query     request.ParamUpdateMsgPin   true   "请求信息"
// @Success   200            {object}  common.State{}              "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 2010:账号不存在 5001:消息不存在"
// @Router    /api/message/update/pin [put]
func (message) UpdateMsgPin(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateMsgPin)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.UpdateMsgPin(ctx, content.ID, params)
	reply.Reply(err)
}

// UpdateMsgTop 更改消息的置顶状态
// @Tags      message
// @Summary   更改消息的置顶状态
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                      true   "Bearer 账户令牌"
// @Param     data           query     request.ParamUpdateMsgTop   true   "请求信息"
// @Success   200            {object}  common.State{}              "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 2010:账号不存在 5001:消息不存在"
// @Router    /api/message/update/top [put]
func (message) UpdateMsgTop(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateMsgTop)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.UpdateMsgTop(ctx, content.ID, params)
	reply.Reply(err)
}

// RevokeMsg 撤回消息
// @Tags      message
// @Summary   撤回消息
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                      true   "Bearer 账户令牌"
// @Param     data           query     request.ParamRevokeMsg      true   "请求信息"
// @Success   200            {object}  common.State{}              "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 2010:账号不存在 5001:消息不存在"
// @Router    /api/message/update/revoke [put]
func (message) RevokeMsg(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamRevokeMsg)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Message.RevokeMsg(ctx, content.ID, params.ID)
	reply.Reply(err)
}

// GetTopMsgByRelationID 获取指定关系中的置顶消息，如果不存在则为 null
// @Tags      message
// @Summary   获取指定关系中的置顶消息，如果不存在则为 null
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                                               true   "Bearer 账户令牌"
// @Param     data           query     request.ParamGetTopMsgByRelationID                   true   "请求信息"
// @Success   200            {object}  common.State{data=reply.ParamGetTopMsgByRelationID}  "完整的消息详情,但不存在则为null"
// @Router    /api/message/info/top [get]
func (message) GetTopMsgByRelationID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetTopMsgByRelationID)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Message.GetTopMsgByRelationID(ctx, content.ID, params.RelationID)
	reply.Reply(err, result)
}

// GetPinMsgsByRelationID 获取指定关系的 pin 消息，按照pin时间倒序排序
// @Tags      message
// @Summary   获取指定关系的 pin 消息，按照pin时间倒序排序
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                                                true   "Bearer 账户令牌"
// @Param     data           query     request.ParamGetPinMsgsByRelationID                   true   "请求信息"
// @Success   200            {object}  common.State{data=reply.ParamGetPinMsgsByRelationID}  "完整的消息详情，包含回复消息"
// @Router    /api/message/list/pin [get]
func (message) GetPinMsgsByRelationID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetPinMsgsByRelationID)
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
	result, err := logic.Logics.Message.GetPinMsgsByRelationID(ctx, content.ID, params.RelationID, limit, offset)
	reply.ReplyList(err, result.Total, result.List)
}

// GetRlyMsgsInfoByMsgID 获取指定消息的所有回复消息，按照回复时间先后排序
// @Tags      message
// @Summary   获取指定消息的所有回复消息，按照回复时间先后排序
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                                               true   "Bearer 账户令牌"
// @Param     data           query     request.ParamGetRlyMsgsInfoByMsgID                   true   "请求信息"
// @Success   200            {object}  common.State{data=reply.ParamGetRlyMsgsInfoByMsgID}  "完整的消息详情，包含回复消息"
// @Router    /api/message/list/reply [get]
func (message) GetRlyMsgsInfoByMsgID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetRlyMsgsInfoByMsgID)
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
	result, err := logic.Logics.Message.GetRlyMsgsInfoByMsgID(ctx, content.ID, params.RelationID, params.MsgID, limit, offset)
	reply.ReplyList(err, result.Total, result.List)
}

// GetMsgsByContent 通过内容模糊查找指定或所有关系中的消息
// @Tags      message
// @Summary   通过内容模糊查找指定或所有关系中的消息，按照时间先后顺序倒序排序，不会查询撤回的消息（指定关系 ID < 0 则查询所有关系中的消息）
// @accept    application/json
// @Produce   application/json
// @Param     Authorization  header    string                                          true   "Bearer 账户令牌"
// @Param     data           query     request.ParamGetMsgsByContent                   true   "请求信息"
// @Success   200            {object}  common.State{data=reply.ParamGetMsgsByContent}  "完整的消息详情，包含回复消息"
// @Router    /api/message/list/content [get]
func (message) GetMsgsByContent(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetMsgsByContent)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Message.GetMsgsByContent(ctx, content.ID, params.RelationID, params.Content, limit, offset)
	reply.ReplyList(err, result.Total, result.List)
}

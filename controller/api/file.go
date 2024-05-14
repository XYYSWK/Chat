package api

import (
	"Chat/errcodes"
	"Chat/global"
	"Chat/logic"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/request"
	"Chat/pkg/gtype"
	"github.com/XYYSWK/Rutils/pkg/app"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type file struct {
}

// PublishFile 上传文件
// @Tags file
// @Summary 上传文件
// @accept multipart/form-data
// @Param Authorization header string true "Bearer 账户令牌"
// @Param file formData request.ParamPublishFile true "文件"
// @Success 200 {object} common.State{data=reply.ParamPublishFile} "1001:参数有误 1003:系统错误 8001:存储失败"
// @Router /api/file/publish [post]
func (file) PublishFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamPublishFile)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	fileType, myErr := gtype.GetFileType(params.File)
	if myErr != nil {
		global.Logger.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		reply.Reply(errcode.ErrServer)
		return
	}
	if fileType != "img" && fileType != "png" && fileType != "jpg" {
		fileType = "file"
	}
	result, err := logic.Logics.File.PublishFile(ctx, model.PublishFile{
		File:       params.File,
		RelationID: params.RelationID,
		AccountID:  params.AccountID,
	})
	reply.Reply(err, result)
}

// DeleteFile 删除文件
// @Tags file
// @Summary 删除文件
// @accept application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Param data body request.ParamDeleteFile true "文件ID"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 8002:文件不存在 8003文件删除失败"
// @Router /api/file/delete [delete]
func (file) DeleteFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteFile)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	err := logic.Logics.File.DeleteFile(ctx, params.FileID)
	reply.Reply(err)
}

// GetRelationFile 获取关系文件列表
// @Tags file
// @Summary 获取关系文件列表
// @accept application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Param data body request.ParamGetRelationFile true "关系ID"
// @Success 200 {object} common.State{data=reply.ParamGetRelationFile} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 8001:存储失败"
// @Router /api/file/getFiles [post]
func (file) GetRelationFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetRelationFile)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.File.GetRelationFile(ctx, params.RelationID)
	reply.Reply(err, result)
}

// UploadAccountAvatar 更新账户头像
// @Tags file
// @Summary 更新账户头像
// @accept multipart/form-data
// @Param Authorization header string true "Bearer 账户令牌"
// @Param file formData file true "文件"
// @Param data body request.ParamUploadAccountAvatar true "文件及账户信息"
// @Success 200 {object} common.State{data=reply.ParamUploadAvatar} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足 8001:存储失败"
// @Router /api/file/avatar/account [post]
func (file) UploadAccountAvatar(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUploadAccountAvatar)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.File.UploadAccountAvatar(ctx, content.ID, params.File)
	reply.Reply(err, result)
}

// GetFileDetailsByID 获取文件详情
// @Tags file
// @Summary 获取文件详情
// @accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 账户令牌"
// @Param data body request.ParamGetFileDetailsByID true "文件及账户信息"
// @Success 200 {object} common.State{data=reply.ParamFile} "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2009:权限不足"
// @Router /api/file/details [post]
func (file) GetFileDetailsByID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetFileDetailsByID)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.File.GetFileDetailsByID(ctx, params.FileID)
	reply.Reply(err, result)
}

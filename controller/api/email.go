package api

import (
	"Chat/logic"
	"Chat/model/request"
	"github.com/XYYSWK/Rutils/pkg/app"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type email struct {
}

// ExistEmail 是否已经注册过该 email
// @Tags email
// @Summary 是否已经注册过该 email
// @accept application/json
// @Produce application/json
// @Param data query request.ParamExistEmail true "email 是否存在"
// @Success 200 {object} common.State{Data=reply.ParamExistEmail} "1001:参数有误 1003:系统错误"
// @Router /api/email/exist [get]
func (email) ExistEmail(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamExistEmail{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.Email.ExistEmail(ctx, params.Email)
	reply.Reply(err, result)
}

// SendMark 发送验证码(邮件)
// @Tags email
// @Summary 发送验证码(邮件)
// @accept application/json
// @Produce application/json
// @Param data body request.ParamSendEmail true "email"
// @Success 200 {object} common.State{} "1001:参数有误 1003:系统错误 2006:邮箱已经注册 2003:邮件发送过于频繁，请稍后再试"
// @Router /api/email/send [post]
func (email) SendMark(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := &request.ParamSendEmail{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	err := logic.Logics.Email.SendMark(ctx, params.Email)
	reply.Reply(err)
}

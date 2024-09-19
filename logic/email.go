package logic

import (
	"Chat/dao"
	"Chat/errcodes"
	"Chat/global"
	"Chat/middlewares"
	"Chat/model/reply"
	"Chat/pkg/emailMark"
	"errors"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/XYYSWK/Lutils/pkg/utils"
	"github.com/gin-gonic/gin"
)

type email struct {
}

// ExistEmail 是否存在 email
func (email) ExistEmail(ctx *gin.Context, emailStr string) (*reply.ParamExistEmail, errcode.Err) {
	// 先在 redis 缓存中查找
	ok, err := dao.Database.Redis.ExistEmail(ctx, emailStr)
	if err == nil {
		return &reply.ParamExistEmail{Exist: ok}, nil
	}
	global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
	// 如果在 redis 中没找到，再到 PostgreSQL 数据库中查找
	ok, err = dao.Database.DB.ExistEmail(ctx, emailStr)
	if err != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamExistEmail{Exist: ok}, nil
}

// CheckEmailNotExists 判断邮箱是否存在
// 先去缓存中查询，如果不在缓存中，再去数据库中查询，如果存在，将邮箱写入缓存中，返回邮箱已注册的错误。（旁路缓存中的读策略）
func CheckEmailNotExists(ctx *gin.Context, emailStr string) errcode.Err {
	result, err := email{}.ExistEmail(ctx, emailStr)
	if err != nil {
		return err
	}
	if result.Exist {
		return errcodes.EmailExists
	}
	exist, myErr := dao.Database.DB.ExistEmail(ctx, emailStr)
	if myErr != nil {
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if exist {
		myErr = dao.Database.Redis.AddEmails(ctx, emailStr)
		if myErr != nil {
			global.Logger.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		return errcodes.EmailExists
	}
	return nil
}

// SendMark 发送验证码(邮件)
func (email) SendMark(emailStr string) errcode.Err {
	// 判断发送邮件的频率
	if global.EmailMark.CheckUserExist(emailStr) {
		return errcodes.EmailSendMany
	}
	// 异步发送邮件(使用工作池)
	global.Worker.SendTask(func() {
		code := utils.RandomString(global.PublicSetting.Rules.CodeLength)
		if err := global.EmailMark.SendMark(emailStr, code); err != nil && !errors.Is(err, emailMark.ErrSendTooMany) {
			global.Logger.Error(err.Error())
		}
	})
	return nil
}

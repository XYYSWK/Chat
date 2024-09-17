package logic

import (
	"Chat/dao"
	db "Chat/dao/postgresql/sqlc"
	"Chat/errcodes"
	"Chat/global"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/reply"
	"Chat/task"
	"errors"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/XYYSWK/Lutils/pkg/password"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type user struct{}

func (user) Register(ctx *gin.Context, emailStr, pwd, code string) (*reply.ParamRegister, errcode.Err) {
	// 判断邮箱是否已经注册过了
	if err := CheckEmailNotExists(ctx, emailStr); err != nil {
		return nil, err // 已经注册过了
	}
	// 校验验证码
	if !global.EmailMark.CheckCode(emailStr, code) {
		return nil, errcodes.EmailCodeNotValid
	}
	// 密码加密
	hashPassword, err := password.HashPassword(pwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 将 user 写入数据库并返回 UserInfo
	userInfo, err := dao.Database.DB.CreateUser(ctx, &db.CreateUserParams{
		Email:    emailStr,
		Password: hashPassword,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 添加邮箱到 redis
	err = dao.Database.Redis.AddEmails(ctx, emailStr)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 创建 token
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	refreshToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}

	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, refreshToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}

	return &reply.ParamRegister{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreateAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  refreshToken,
		},
	}, nil
}

func (user) Login(ctx *gin.Context, emailStr, pwd string) (*reply.ParamLogin, errcode.Err) {
	userInfo, myerr := getUserInfoByEmail(ctx, emailStr)
	if myerr != nil {
		return nil, myerr
	}
	if err := password.CheckPassword(pwd, userInfo.Password); err != nil {
		return nil, errcodes.PasswordNotValid
	}
	// 创建 token
	accessToken, accessPayload, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.AccessTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	refreshToken, _, err := newUserToken(model.UserToken, userInfo.ID, global.PrivateSetting.Token.RefreshTokenExpire)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if err = dao.Database.Redis.SaveUserToken(ctx, userInfo.ID, []string{accessToken, refreshToken}); err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	return &reply.ParamLogin{
		ParamUserInfo: reply.ParamUserInfo{
			ID:       userInfo.ID,
			Email:    userInfo.Email,
			CreateAt: userInfo.CreateAt,
		},
		Token: reply.ParamToken{
			AccessToken:   accessToken,
			AccessPayload: accessPayload,
			RefreshToken:  refreshToken,
		},
	}, nil
}

// UpdateUserPassword 更新用户密码
// 参数：code 验证码 newPwd 新密码
func (user) UpdateUserPassword(ctx *gin.Context, userID int64, code, newPwd string) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		global.Logger.Error(myerr.Error(), middlewares.ErrLogMsg(ctx)...)
		return myerr
	}
	// 校验验证码
	if !global.EmailMark.CheckCode(userInfo.Email, code) {
		return errcodes.EmailCodeNotValid
	}
	// 密码加密
	hashPassword, err := password.HashPassword(newPwd)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 更新密码
	if err = dao.Database.DB.UpdateUser(ctx, &db.UpdateUserParams{
		Email:    userInfo.Email,
		Password: hashPassword,
		ID:       userInfo.ID,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}

	// 清除用户的 token
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, userID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

// UpdateUserEmail 修改用户邮箱
func (user) UpdateUserEmail(ctx *gin.Context, userID int64, emailStr, code string) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		return myerr
	}
	// 判断邮箱不是之前那个邮箱
	if userInfo.Email == emailStr {
		return errcodes.EmailSame
	}
	// 判断邮箱没有被注册过
	if err := CheckEmailNotExists(ctx, emailStr); err != nil {
		return err
	}
	// 校验验证码
	if !global.EmailMark.CheckCode(emailStr, code) {
		return errcodes.EmailCodeNotValid
	}
	if err := dao.Database.DB.UpdateUser(ctx, &db.UpdateUserParams{
		Email:    emailStr,
		Password: userInfo.Password,
		ID:       userInfo.ID,
	}); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 更新 redis 中的用户邮箱
	if err := dao.Database.Redis.UpdateEmail(ctx, userInfo.Email, emailStr); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送更改邮箱通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateEmail(accessToken, userID, emailStr))
	return nil
}

// Logout 退出登录
func (user) Logout(ctx *gin.Context) errcode.Err {
	Token, payload, err := GetTokenAndPayload(ctx)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcodes.AuthenticationFailed
	}
	content := &model.Content{}
	_ = content.Unmarshal(payload.Content)
	//先判断用户在redis中是否存在
	if ok := dao.Database.Redis.CheckUserTokenValid(ctx, content.ID, Token); !ok {
		return errcodes.UserNotFound
	}
	//先将token从redis中清除
	if err := dao.Database.Redis.DeleteAllTokenByUser(ctx, content.ID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

// DeleteUser 删除用户
func (user) DeleteUser(ctx *gin.Context, userID int64) errcode.Err {
	userInfo, myerr := getUserInfoByID(ctx, userID)
	if myerr != nil {
		return myerr
	}
	// 查询 user 的 账户(account)
	accountNum, err := dao.Database.DB.CountAccountsByUserID(ctx, userID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if accountNum > 0 {
		return errcodes.UserHasAccount
	}
	// 从 postgresql 中删除 user
	if err := dao.Database.DB.DeleteUser(ctx, userID); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 从 redis 中删除 user 的 email
	if err := dao.Database.Redis.DeleteEmail(ctx, userInfo.Email); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		reTry("deleteEmail:"+userInfo.Email, func() error {
			return dao.Database.Redis.DeleteEmail(ctx, userInfo.Email)
		})
	}
	return nil
}

// getUserInfoByID 通过用户 ID 获取用户信息
// 参数：userID 用户ID
// 成功：用户信息，nil
// 失败：打印日志 errcodes.UserNotFound, errcode.ErrServer
func getUserInfoByID(ctx *gin.Context, userID int64) (*db.User, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return userInfo, nil
}

// getUserInfoByEmail 通过邮箱获取用户信息
// 参数：emailStr 邮箱
// 成功：用户信息，nil
// 失败：打印日志 errcodes.UserNotFound, errcode.ErrServer
func getUserInfoByEmail(ctx *gin.Context, emailStr string) (*db.User, errcode.Err) {
	userInfo, err := dao.Database.DB.GetUserByEmail(ctx, emailStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errcodes.UserNotFound
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return userInfo, nil
}

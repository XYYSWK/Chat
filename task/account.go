package task

import (
	"Chat/dao"
	"Chat/global"
	"Chat/model/chat"
	"Chat/model/chat/server"
	"github.com/XYYSWK/Lutils/pkg/utils"
)

/*
有关 account 的 任务
*/

// UpdateEmail 更新邮箱时，通知用户(user)的每个账户(account)
func UpdateEmail(accessToken string, userID int64, email string) func() {
	return func() {
		ctx, cancel := global.DefaultContextWithTimeout()
		defer cancel()
		accountIDs, err := dao.Database.DB.GetAcountIDsByUserID(ctx, userID)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
		global.ChatMap.SendMany(accountIDs, chat.ServerUpdateEmail, server.UpdateEmail{
			EnToken: utils.EncodeMD5(accessToken),
			Email:   email,
		})
	}
}

// UpdateAccount 更新账号信息时，通知该账号成功更新后的信息
func UpdateAccount(accessToken string, accountID int64, name, gender, signature string) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerUpdateAccount, server.UpdateAccount{
			EnToken:   utils.EncodeMD5(accessToken),
			Name:      name,
			Gender:    gender,
			Signature: signature,
		})
	}
}

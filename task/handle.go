package task

import (
	"Chat/global"
	"Chat/model/chat"
	"Chat/model/chat/server"
	"github.com/XYYSWK/Lutils/pkg/utils"
)

// AccountLogin 发送账户上线的通知
func AccountLogin(accessToken, address string, accountID int64) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerAccountLogin, server.AccountLogin{
			EnToken: utils.EncodeMD5(accessToken),
			Address: address,
		})
	}
}

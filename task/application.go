package task

import (
	"Chat/global"
	"Chat/model/chat"
)

func Application(accountID int64) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerApplication)
	}
}

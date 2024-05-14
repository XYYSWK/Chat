package settings

import (
	"Chat/global"
	"Chat/manager"
)

type chat struct {
}

func (chat) Init() {
	global.ChatMap = manager.NewChatMap()
}

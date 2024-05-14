package task

import (
	"Chat/global"
	"Chat/model/chat"
	"Chat/model/chat/server"
	"github.com/XYYSWK/Rutils/pkg/utils"
)

// UpdateNickName 更新昵称的通知
func UpdateNickName(accessToken string, accountID, relationID int64, nickName string) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerUpdateNickName, server.UpdateNickName{
			EnToken:    utils.EncodeMD5(accessToken),
			RelationID: relationID,
			NickName:   nickName,
		})
	}
}

// UpdateSettingState 更新 relation 状态的通知
func UpdateSettingState(accessToken string, settingType server.SettingType, accountID, relationID int64, state bool) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerUpdateSettingState, server.UpdateSettingState{
			EnToken:    utils.EncodeMD5(accessToken),
			RelationID: relationID,
			Type:       settingType,
			State:      state,
		})
	}
}

// DeleteRelation 删除关系的通知
func DeleteRelation(accessToken string, accountID, relationID int64) func() {
	return func() {
		global.ChatMap.Send(accountID, chat.ServerDeleteRelation, server.DeleteRelation{
			EnToken:    utils.EncodeMD5(accessToken),
			RelationID: relationID,
		})
	}
}

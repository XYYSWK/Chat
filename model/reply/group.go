package reply

import "Chat/model"

type ParamCreateGroup struct {
	Name        string `json:"name"`
	AccountID   int64  `json:"account_id"`
	RelationID  int64  `json:"relation_id"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type ParamUpdateGroup struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	Avatar      string `json:"avatar" form:"avatar" binding:"required"`
}

type ParamInviteAccount struct {
	InviteMember []int64 `json:"invite_member"`
}

type ParamGetGroupList struct {
	List  []*model.SettingGroup
	Total int64
}

type ParamGetGroupsByName struct {
	List  []model.SettingGroup
	Total int64
}

type ParamGetGroupMembers struct {
	List  []ParamGroupMemberInfo
	Total int64
}

type ParamGroupMemberInfo struct {
	AccountID int64  `json:"account_id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	IsLeader  bool   `json:"is_leader"`
}

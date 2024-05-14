package request

import "Chat/model/common"

type ParamUpdateNickName struct {
	RelationID int64  `json:"relation_id" binding:"required,gte=1"`      // 关系 ID
	NickName   string `json:"nick_name" binding:"required,gte=1,tle=20"` // 昵称
}

type ParamUpdateSettingPin struct {
	RelationID int64 `json:"relation_id" binding:"required,gte=1"` // 关系 ID
	IsPin      *bool `json:"is_pin" binding:"required"`            // 是否 pin（置顶）
}

type ParamUpdateSettingDisturb struct {
	RelationID   int64 `json:"relation_id" binding:"required,gte=1"` // 关系 ID
	IsNotDisturb *bool `json:"is_not_disturb" binding:"required"`    // 是否免打扰
}

type ParamUpdateSettingShow struct {
	RelationID int64 `json:"relation_id" binding:"required,gte=1"` // 关系 ID
	IsShow     *bool `json:"is_show" binding:"required"`           // 是否展示
}

type ParamDeleteFriend struct {
	RelationID int64 `json:"relation_id" binding:"required,gte=1"` // 要删除的关系的 ID
}

type ParamGetFriendsByName struct {
	Name string `json:"name" binding:"required,gte=1"` // 要查询的名称
	common.Page
}

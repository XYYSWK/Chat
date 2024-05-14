package server

type SettingType string

const (
	SettingPin        = "pin"
	SettingShow       = "show"
	SettingNotDisturb = "not_disturb"
)

type UpdateNickName struct {
	EnToken    string `json:"en_token"`            // 加密过的 token
	NickName   string `json:"nick_name,omitempty"` // 更改过的新昵称
	RelationID int64  `json:"relation_id,omitempty"`
}

type UpdateSettingState struct {
	EnToken    string      `json:"en_token"`              // 加密过的 token
	RelationID int64       `json:"relation_id,omitempty"` // 关系 ID
	Type       SettingType `json:"type"`                  // 通知类型[pin, show, not_disturb]
	State      bool        `json:"state"`                 // 状态设置
}

type DeleteRelation struct {
	EnToken    string `json:"en_token"`
	RelationID int64  `json:"relation_id,omitempty"`
}

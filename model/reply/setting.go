package reply

import "Chat/model"

type ParamGetPins struct {
	List  []*model.SettingPin `json:"list"`
	Total int64               `json:"total"`
}

type ParamGetShows struct {
	List  []*model.Setting `json:"list,omitempty"`
	Total int64            `json:"total,omitempty"`
}

type ParamGetFriends struct {
	List  []*model.SettingFriend `json:"list,omitempty"`
	Total int64                  `json:"total,omitempty"`
}

type ParamGetFriendsByName struct {
	List  []*model.SettingFriend `json:"list,omitempty"`
	Total int64                  `json:"total,omitempty"`
}

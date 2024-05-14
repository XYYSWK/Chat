package routers

import (
	"Chat/controller/api"
	"Chat/middlewares"
	"github.com/gin-gonic/gin"
)

type setting struct {
}

func (setting) Init(router *gin.RouterGroup) {
	r := router.Group("setting", middlewares.MustAccount())
	{
		updateGroup := r.Group("update")
		{
			updateGroup.PUT("nick_name", api.Apis.Setting.UpdateNickName)
			updateGroup.PUT("pin", api.Apis.Setting.UpdateSettingPin)
			updateGroup.PUT("disturb", api.Apis.Setting.UpdateSettingDisturb)
			updateGroup.PUT("show", api.Apis.Setting.UpdateSettingShow)
		}
		r.GET("pins", api.Apis.Setting.GetPins)
		r.GET("shows", api.Apis.Setting.GetShows)
		friendGroup := r.Group("friend")
		{
			friendGroup.GET("list", api.Apis.Setting.GetFriends)
			friendGroup.DELETE("delete", api.Apis.Setting.DeleteFriend)
			friendGroup.GET("name", api.Apis.Setting.GetFriendsByName)
		}
	}
}

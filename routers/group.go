package routers

import (
	"Chat/controller/api"
	"Chat/middlewares"
	"github.com/gin-gonic/gin"
)

type group struct {
}

func (group) Init(router *gin.RouterGroup) {
	r := router.Group("group").Use(middlewares.MustAccount())
	{
		r.POST("create", api.Apis.Group.CreateGroup)
		r.POST("transfer", api.Apis.Group.TransferGroup)
		r.POST("dissolve", api.Apis.Group.DissolveGroup)
		r.POST("update", api.Apis.Group.UpdateGroup)
		r.POST("invite", api.Apis.Group.InviteAccount)
		r.GET("list", api.Apis.Group.GetGroupList)
		r.POST("quit", api.Apis.Group.QuitGroup)
		r.GET("name", api.Apis.Group.GetGroupsByName)
		r.GET("members", api.Apis.Group.GetGroupMembers)
	}
}

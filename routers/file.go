package routers

import (
	"Chat/controller/api"
	"Chat/middlewares"
	"github.com/gin-gonic/gin"
)

type file struct {
}

func (file) Init(router *gin.RouterGroup) {
	r := router.Group("file", middlewares.MustAccount())
	{
		r.POST("publish", api.Apis.File.PublishFile)
		r.DELETE("delete", api.Apis.File.DeleteFile)
		r.POST("getFiles", api.Apis.File.GetRelationFile)
		avatarGroup := r.Group("avatar")
		{
			avatarGroup.PUT("account", api.Apis.File.UploadAccountAvatar)

		}
		r.POST("details", api.Apis.File.GetFileDetailsByID)
	}
}

package routers

import (
	"Chat/controller/api"
	"github.com/gin-gonic/gin"
)

type email struct {
}

func (email) Init(router *gin.RouterGroup) {
	r := router.Group("email")
	{
		r.GET("exist", api.Apis.Email.ExistEmail)
		r.POST("send", api.Apis.Email.SendMark)
	}
}

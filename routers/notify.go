package routers

import (
	"Chat/controller/api"
	"Chat/middlewares"
	"github.com/gin-gonic/gin"
)

type notify struct {
}

func (notify) Init(router *gin.RouterGroup) {
	r := router.Group("notify").Use(middlewares.MustAccount())
	{
		r.POST("create", api.Apis.Notify.CreateNotify)
		r.PUT("update", api.Apis.Notify.UpdateNotify)
		r.GET("get", api.Apis.Notify.GetNotifyByID)
		r.DELETE("delete", api.Apis.Notify.DeleteNotify)
	}
}

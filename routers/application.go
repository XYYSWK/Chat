package routers

import (
	"Chat/controller/api"
	"Chat/middlewares"
	"github.com/gin-gonic/gin"
)

type application struct {
}

func (application) Init(router *gin.RouterGroup) {
	r := router.Group("/application").Use(middlewares.MustAccount())
	{
		r.POST("create", api.Apis.Application.CreateApplication)
		r.DELETE("delete", api.Apis.Application.DeleteApplication)
		r.PUT("refuse", api.Apis.Application.RefuseApplication)
		r.PUT("accept", api.Apis.Application.AcceptApplication)
		r.GET("list", api.Apis.Application.ListApplications)
	}
}

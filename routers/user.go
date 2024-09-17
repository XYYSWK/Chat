package routers

import (
	"Chat/controller/api"
	"Chat/middlewares"
	"github.com/gin-gonic/gin"
)

type user struct {
}

func (user) Init(router *gin.RouterGroup) {
	r := router.Group("user")
	{
		r.POST("register", api.Apis.User.Register)
		r.POST("login", api.Apis.User.Login)
		updateGroup := r.Group("update").Use(middlewares.MustUser()) // 添加鉴权中间件,以及 MustUser 中间件
		{
			updateGroup.PUT("pwd", api.Apis.User.UpdateUserPassword)
			updateGroup.PUT("email", api.Apis.User.UpdateUserEmail)
			updateGroup.GET("/logout", api.Apis.User.Logout)
		}
		r.DELETE("deleteUser", middlewares.MustUser(), api.Apis.User.DeleteUser)
	}
}

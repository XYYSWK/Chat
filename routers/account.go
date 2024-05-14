package routers

import (
	"Chat/controller/api"
	"Chat/middlewares"
	"github.com/gin-gonic/gin"
)

type account struct {
}

func (account) Init(router *gin.RouterGroup) {
	r := router.Group("account")
	{
		userGroup := r.Group("").Use(middlewares.MustUser())
		{
			userGroup.POST("create", api.Apis.Account.CreateAccount)
			userGroup.GET("token", api.Apis.Account.GetAccountToken)
			userGroup.DELETE("delete", api.Apis.Account.DeleteAccount)
			userGroup.GET("infos/account", api.Apis.Account.GetAccountsByUserID)
		}
		accountGroup := r.Group("").Use(middlewares.MustAccount())
		{
			accountGroup.PUT("update", api.Apis.Account.UpdateAccount)
			accountGroup.GET("infos/name", api.Apis.Account.GetAccountsByName)
			accountGroup.GET("info", api.Apis.Account.GetAccountByID)
		}
	}
}

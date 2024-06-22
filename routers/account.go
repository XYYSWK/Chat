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
			userGroup.POST("token", api.Apis.Account.GetAccountToken)
			userGroup.DELETE("delete", api.Apis.Account.DeleteAccount)
			userGroup.POST("infos/account", api.Apis.Account.GetAccountsByUserID)
		}
		accountGroup := r.Group("").Use(middlewares.MustAccount())
		{
			accountGroup.PUT("update", api.Apis.Account.UpdateAccount)
			accountGroup.POST("infos/name", api.Apis.Account.GetAccountsByName)
			accountGroup.POST("info", api.Apis.Account.GetAccountByID)
		}
	}
}

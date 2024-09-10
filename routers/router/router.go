package router

import (
	_ "Chat/docs"
	"Chat/global"
	"Chat/middlewares"
	"Chat/routers"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func NewRouter() (*gin.Engine, *socketio.Server) {
	//创建一个新的路由
	r := gin.New()
	r.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.Recovery(true))

	root := r.Group("api", middlewares.LogBody(), middlewares.PasetoAuth())
	{
		root.GET("swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
		// 响应测试 Ping-Pong
		root.GET("ping", func(ctx *gin.Context) {
			reply := app.NewResponse(ctx)
			global.Logger.Info("ping", middlewares.ErrLogMsg(ctx)...)
			reply.Reply(nil, "pong")
		})

		rg := routers.Routers
		rg.User.Init(root)
		rg.Email.Init(root)
		rg.Account.Init(root)
		rg.Application.Init(root)
		rg.File.Init(root)
		rg.Message.Init(root)
		rg.Setting.Init(root)
		rg.Group.Init(root)
		rg.Notify.Init(root)
	}
	return r, routers.Routers.Chat.Init(r)
}

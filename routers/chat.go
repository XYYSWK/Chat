package routers

import (
	"Chat/controller/api"
	chat2 "Chat/model/chat"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ws struct {
}

func (ws) Init(router *gin.Engine) *socketio.Server {
	server := socketio.NewServer(nil) // 创建一个 socketIO 服务器对象
	{
		// 定义处理客户端事件的函数
		server.OnConnect("/", api.Apis.Chat.Handle.OnConnect)
		server.OnError("/", api.Apis.Chat.Handle.OnError)
		server.OnDisconnect("/", api.Apis.Chat.Handle.OnDisconnect)
	}
	chatHandle(server)
	// 将 server（Socket.IO 服务器）绑定到 Gin 框架的路由上，以处理客户端发起的 WebSocket 连接请求。
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	return server
}

func chatHandle(server *socketio.Server) {
	namespace := "/chat"
	server.OnEvent(namespace, chat2.ClientSendMsg, api.Apis.Chat.Message.SendMsg)
	server.OnEvent(namespace, chat2.ClientReadMsg, api.Apis.Chat.Message.ReadMsg)
	server.OnEvent(namespace, chat2.ClientAuth, api.Apis.Chat.Handle.Auth) // 账户登录
	server.OnEvent(namespace, chat2.ClientTest, api.Apis.Chat.Handle.Test)
}

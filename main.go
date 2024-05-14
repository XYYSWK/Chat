package main

import (
	"Chat/global"
	"Chat/routers/router"
	"Chat/settings"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//1. 初始化项目（配置加载，日志、数据库，雪花算法...初始化等等）
	settings.Inits()
	//设置 Gin 框架为 Release（生产）模式
	if global.PublicSetting.Server.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	//2.注册路由，返回 路由 和 Socket.IO 服务器实例
	r, ws := router.NewRouter()

	//3.启动服务（优雅关机）
	//http.Server 内置的 Shutdown() 方法支持优雅关机
	sever := http.Server{
		Addr:           global.PublicSetting.Server.HttpPort, //端口号
		Handler:        r,                                    //路由处理器
		MaxHeaderBytes: 1 << 20,                              //最大请求头大小（1MB）
		//设置合适的 MaxHeaderBytes 值，可以确保服务器能够有效地处理请求头，避免不必要的资源浪费或潜在的安全风险
	}
	global.Logger.Info("Server started!") //输出日志，服务器已启动
	fmt.Println("AppName:", global.PublicSetting.App.Name, "Version:", global.PublicSetting.App.Version, "Address:", global.PublicSetting.Server.HttpPort, "RunMode:", global.PublicSetting.Server.RunMode)

	errChan := make(chan error, 1)
	defer close(errChan) //延迟关闭错误通道

	go func() {
		//启动 HTTP 服务器
		err := sever.ListenAndServe()
		if err != nil {
			errChan <- err //将错误发送到错误通道
		}
	}()

	// 启动 Socket.IO 服务器
	go func() {
		defer ws.Close()
		// 接收并处理网络连接
		if err := ws.Serve(); err != nil {
			errChan <- err
		}
	}()

	//创建一个接收信号的通道
	quit := make(chan os.Signal, 1)                      //os.Signal 表示操作系统的信号，比如中断信号、终止信号等。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit

	select {
	case err := <-errChan:
		global.Logger.Error(err.Error())
	case <-quit:
		global.Logger.Info("Shutdown Server.")
		///创建一个带超时的上下文（给几秒完成还未处理完的请求）
		ctx, cancel := context.WithTimeout(context.Background(), global.PublicSetting.Server.DefaultContextTimeout)
		defer cancel() //延迟取消上下文

		//上下文超时时间内优雅关机（将未处理完的请求处理完再关闭服务），超过超时时间时退出
		if err := sever.Shutdown(ctx); err != nil {
			global.Logger.Error("Server forced to Shutdown, err:" + err.Error())
		}
	}

	global.Logger.Info("Server exited!")
}

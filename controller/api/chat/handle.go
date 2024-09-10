package chat

import (
	"Chat/global"
	"Chat/model/chat/client"
	"Chat/model/common"
	"Chat/task"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"time"
)

type handle struct {
}

const AuthLimitTimeout = 10 * time.Second

// OnConnect
// 当客户端连接时触发
func (handle) OnConnect(s socketio.Conn) error {
	log.Println("connected:", s.RemoteAddr().String(), s.ID()) // s.RemoteAddr()获取客户端的 IP 地址和端口号信息。
	// 一定时间内需要进行 AUTH 认证，否则断开连接
	time.AfterFunc(AuthLimitTimeout, func() {
		if !global.ChatMap.HasSID(s.ID()) {
			global.Logger.Info(fmt.Sprintln("auth failed:", s.RemoteAddr().String(), s.ID()))
			_ = s.Close()
		}
	})
	return nil
}

// OnError
// 当发生错误时触发
func (handle) OnError(s socketio.Conn, err error) {
	log.Println("on error:", err)
	if s == nil {
		return
	}
	// 从在线中退出
	global.ChatMap.Leave(s)
	log.Println("disconnected:", s.RemoteAddr().String(), s.ID())
	_ = s.Close()
}

// OnDisconnect
// 当客户端断开连接时触发
func (handle) OnDisconnect(s socketio.Conn, _ string) {
	// 从在线中退出
	global.ChatMap.Leave(s)
	log.Println("disconnected:", s.RemoteAddr().String(), s.ID())
}

// Test 测试
func (handle) Test(s socketio.Conn, msg string) string {
	_, ok := CheckAuth(s)
	if !ok {
		return ""
	}
	params := new(client.TestParams)
	log.Println(msg)
	if err := common.Decode(msg, params); err != nil {
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	result := common.NewState(nil, client.TestRly{
		Name:    params.Name,
		Age:     params.Age,
		Address: s.RemoteAddr().String(),
		ID:      s.ID(),
	}).MustJson()
	// test
	s.Emit("test", "test")
	return result
}

// Auth 身份验证（account 上线）
func (handle) Auth(s socketio.Conn, accessToken string) string {
	token, myErr := MustAccount(accessToken)
	if myErr != nil {
		return common.NewState(myErr).MustJson()
	}
	s.SetContext(token)
	// 加入在线群组
	global.ChatMap.Link(s, token.Content.ID)
	// 通知其他设备
	global.Worker.SendTask(task.AccountLogin(accessToken, s.RemoteAddr().String(), token.Content.ID))
	log.Println("auth accept:", s.RemoteAddr().String())
	return common.NewState(nil).MustJson()
}

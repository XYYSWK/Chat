package chat

import (
	"Chat/chat"
	"Chat/global"
	"Chat/model"
	"Chat/model/chat/client"
	"Chat/model/common"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	socketio "github.com/googollee/go-socket.io"
)

// 用于处理客户端发送的 event
type message struct {
}

// SendMsg 发送消息
// 参数：client.HandleSendMsgParams
// 返回：client.HandleSendMsgRly
func (message) SendMsg(s socketio.Conn, msg string) string {
	token, ok := CheckAuth(s)
	if !ok {
		return ""
	}
	params := new(client.HandleSendMsgParams)
	if err := common.Decode(msg, params); err != nil {
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	result, myErr := chat.Group.Message.SendMsg(ctx, &model.HandleSendMsg{
		AccessToken: token.AccessToken,
		RelationID:  params.RelationID,
		AccountID:   token.Content.ID,
		MsgContent:  params.MsgContent,
		MsgExtend:   params.MsgExtend,
		RlyMsgID:    params.RlyMsgID,
	})
	return common.NewState(myErr, result).MustJson()
}

// ReadMsg 已读消息
// 参数：client.HandleReadMsgParams
// 返回：无
func (message) ReadMsg(s socketio.Conn, msg string) string {
	token, ok := CheckAuth(s)
	if !ok {
		return ""
	}
	params := new(client.HandleReadMsgParams)
	if err := common.Decode(msg, params); err != nil {
		return common.NewState(errcode.ErrParamsNotValid.WithDetails(err.Error())).MustJson()
	}
	ctx, cancel := global.DefaultContextWithTimeout()
	defer cancel()
	myErr := chat.Group.Message.ReadMsg(ctx, &model.HandleReadMsg{
		AccessToken: token.AccessToken,
		MsgIDs:      params.MsgIDs,
		RelationID:  params.RelationID,
		ReaderID:    token.Content.ID,
	})
	return common.NewState(myErr).MustJson()
}

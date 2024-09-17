package chat

import (
	"Chat/dao"
	db "Chat/dao/postgresql/sqlc"
	"Chat/errcodes"
	"Chat/global"
	"Chat/logic"
	"Chat/model"
	"Chat/model/chat/client"
	"Chat/model/reply"
	"Chat/task"
	"context"
	"database/sql"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
)

type message struct {
}

func (message) SendMsg(ctx context.Context, params *model.HandleSendMsg) (*client.HandleSendMsgRly, errcode.Err) {
	// 判断权限
	ok, myErr := logic.ExistsSetting(ctx, params.AccountID, params.RelationID)
	if myErr != nil {
		return nil, myErr
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	var rlyMsgID int64
	var rlyMsg *reply.ParamRlyMsg
	// 判断并回复消息
	if params.RlyMsgID > 0 {
		rlyInfo, myErr := logic.GetMsgInfoByID(ctx, params.RlyMsgID)
		if myErr != nil {
			return nil, myErr
		}
		// 不能回复别的群的消息
		if rlyInfo.RelationID != params.RelationID {
			return nil, errcodes.RlyMsgNotOneRelation
		}
		// 不能回复已经撤回的消息
		if rlyInfo.IsRevoke {
			return nil, errcodes.RlyMsgHasRevoked
		}
		rlyMsgID = params.RlyMsgID
		rlyMsgExtend, err := model.JsonToExtend(rlyInfo.MsgExtend)
		if err != nil {
			global.Logger.Error(err.Error())
			return nil, errcode.ErrServer
		}
		rlyMsg = &reply.ParamRlyMsg{
			MsgID:      rlyInfo.ID,
			MsgType:    rlyInfo.MsgType,
			MsgContent: rlyInfo.MsgContent,
			MsgExtend:  rlyMsgExtend,
			IsRevoked:  rlyInfo.IsRevoke,
		}
	}
	msgExtend, err := model.ExtendToJson(params.MsgExtend)
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	result, err := dao.Database.DB.CreateMessage(ctx, &db.CreateMessageParams{
		NotifyType: db.MsgnotifytypeCommon,
		MsgType:    string(model.MsgTypeText),
		MsgContent: params.MsgContent,
		MsgExtend:  msgExtend,
		AccountID:  sql.NullInt64{Int64: params.AccountID, Valid: true},
		RlyMsgID:   sql.NullInt64{Int64: rlyMsgID, Valid: rlyMsgID > 0},
		RelationID: params.RelationID,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}

	// 推送消息
	global.Worker.SendTask(task.PublishMsg(reply.ParamMsgInfoWithRly{
		ParamMsgInfo: reply.ParamMsgInfo{
			ID:         result.ID,
			NotifyType: string(db.MsgnotifytypeCommon),
			MsgType:    string(model.MsgTypeText),
			MsgContent: result.MsgContent,
			MsgExtend:  params.MsgExtend,
			AccountID:  params.AccountID,
			RelationID: params.RelationID,
			CreateAt:   result.CreateAt,
		},
		RlyMsg: rlyMsg,
	}))
	return &client.HandleSendMsgRly{
		MsgID:    result.ID,
		CreateAt: result.CreateAt,
	}, nil
}

func (message) ReadMsg(ctx context.Context, params *model.HandleReadMsg) errcode.Err {
	// 判断权限
	ok, myErr := logic.ExistsSetting(ctx, params.ReaderID, params.RelationID)
	if myErr != nil {
		return myErr
	}
	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	readMsgs, err := dao.Database.DB.UpdateMsgReads(ctx, &db.UpdateMsgReadsParams{
		RelationID: params.RelationID,
		Accountid:  params.ReaderID,
		Msgids:     params.MsgIDs,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return errcode.ErrServer
	}
	msgMap := make(map[int64][]int64)
	for _, v := range readMsgs {
		msgID, accountID := v.ID, v.AccountID // 消息ID，发消息者的ID
		msgMap[accountID] = append(msgMap[accountID], msgID)
	}
	// 推送消息已经被读取
	global.Worker.SendTask(task.ReadMsg(params.AccessToken, params.ReaderID, msgMap, params.MsgIDs))
	return nil
}

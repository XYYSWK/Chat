package logic

import (
	"Chat/dao"
	db "Chat/dao/postgresql/sqlc"
	"Chat/errcodes"
	"Chat/global"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/chat/server"
	"Chat/model/format"
	"Chat/model/reply"
	"Chat/model/request"
	"Chat/task"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type message struct{}

func (message) CreateFileMsg(ctx *gin.Context, params model.CreateFileMsg) (*reply.ParamCreateFileMsg, errcode.Err) {
	// 检查权限
	ok, myErr := ExistsSetting(ctx, params.AccountID, params.RelationID)
	if myErr != nil {
		return nil, myErr
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	// 上传文件
	fileInfo, myErr := Logics.File.PublishFile(ctx, model.PublishFile{
		File:       params.File,
		RelationID: params.RelationID,
		AccountID:  params.AccountID,
	})
	if myErr != nil {
		return nil, myErr
	}
	var isRly bool  // 是否是回复别人的消息
	var rlyID int64 // 回复 ID 为 rlyID 的消息
	var rlyMsg *reply.ParamRlyMsg
	if params.RlyMsgID > 0 { // 如果是回复别人的消息
		rltInfo, myErr := GetMsgInfoByID(ctx, params.RlyMsgID)
		if myErr != nil {
			return nil, myErr
		}
		if rltInfo.IsRevoke {
			return nil, errcodes.RlyMsgHasRevoked
		}
		isRly = true
		rlyID = params.RlyMsgID
		rlyMsgExtend, err := model.JsonToExtend(rltInfo.MsgExtend)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
		rlyMsg = &reply.ParamRlyMsg{
			MsgID:      rltInfo.ID,
			MsgType:    rltInfo.MsgType,
			MsgContent: rltInfo.MsgContent,
			MsgExtend:  rlyMsgExtend,
			IsRevoked:  rltInfo.IsRevoke,
		}
	}
	extend, _ := model.ExtendToJson(nil)
	result, err := dao.Database.DB.CreateMessage(ctx, &db.CreateMessageParams{
		NotifyType: db.MsgnotifytypeCommon,
		MsgType:    string(model.MsgTypeFile),
		MsgContent: fileInfo.Url,
		MsgExtend:  extend,
		FileID:     sql.NullInt64{Int64: fileInfo.ID, Valid: true},
		AccountID:  sql.NullInt64{Int64: params.AccountID, Valid: true},
		RlyMsgID:   sql.NullInt64{Int64: rlyID, Valid: isRly},
		RelationID: params.RelationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 获取 token
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	// 推送消息
	global.Worker.SendTask(task.PublishMsg(accessToken, reply.ParamMsgInfo{
		ID:         result.ID,
		NotifyType: string(db.MsgnotifytypeCommon),
		MsgType:    string(model.MsgTypeFile),
		MsgContent: fileInfo.Url,
		MsgExtend:  nil,
		FileID:     fileInfo.ID,
		AccountID:  params.AccountID,
		RelationID: params.RelationID,
		CreateAt:   fileInfo.CreateAt,
	}, rlyMsg))
	return &reply.ParamCreateFileMsg{
		ID:         result.ID,
		MsgContent: result.MsgContent,
		FileID:     result.FileID.Int64,
		CreateAt:   result.CreateAt,
	}, nil
}

// GetMsgInfoByID 获取消息详情
// 参数：msgID 消息ID
// 成功：消息详情，nil
// 失败：打印错误日志 errcode.ErrServer,errcodes.MsgNotExists
func GetMsgInfoByID(ctx context.Context, msgID int64) (*db.GetMessageByIDRow, errcode.Err) {
	result, err := dao.Database.DB.GetMessageByID(ctx, msgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.MsgNotExists
		}
		global.Logger.Error(err.Error())
		return nil, errcode.ErrServer
	}
	return result, nil
}

func (message) GetMsgsByRelationIDAndTime(ctx *gin.Context, params model.GetMsgsByRelationIDAndTime) (*reply.ParamGetMsgsRelationIDAndTime, errcode.Err) {
	// 权限验证
	ok, myErr := ExistsSetting(ctx, params.AccountID, params.RelationID)
	if myErr != nil {
		return nil, myErr
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	data, err := dao.Database.DB.GetMsgsByRelationIDAndTime(ctx, &db.GetMsgsByRelationIDAndTimeParams{
		RelationID: params.RelationID,
		CreateAt:   params.LastTime,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetMsgsRelationIDAndTime{List: []*reply.ParamMsgInfoWithRly{}}, nil
	}
	result := make([]*reply.ParamMsgInfoWithRly, 0, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke { // 该消息没有被撤回
			content = v.MsgContent
			extend, err = model.JsonToExtend(v.MsgExtend)
			if err != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				continue
			}
		}
		var readIDs []int64
		if params.AccountID == v.AccountID.Int64 {
			readIDs = v.ReadIds
		}
		var rlyMsg *reply.ParamRlyMsg
		if v.RlyMsgID.Valid { // 该 ID 有意义
			rlyMsgInfo, myErr := GetMsgInfoByID(ctx, v.RlyMsgID.Int64)
			if myErr != nil {
				continue
			}
			var rlyContent string
			var rlyExtend *model.MsgExtend
			if !rlyMsgInfo.IsRevoke { // 回复消息没有撤回
				rlyContent = rlyMsgInfo.MsgContent
				rlyExtend, err = model.JsonToExtend(rlyMsgInfo.MsgExtend)
				if err != nil {
					global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
					continue
				}
			}
			rlyMsg = &reply.ParamRlyMsg{
				MsgID:      v.RlyMsgID.Int64,
				MsgType:    rlyMsgInfo.MsgType,
				MsgContent: rlyContent,
				MsgExtend:  rlyExtend,
				IsRevoked:  rlyMsgInfo.IsRevoke,
			}
		}
		result = append(result, &reply.ParamMsgInfoWithRly{
			ParamMsgInfo: reply.ParamMsgInfo{
				ID:         v.ID,
				NotifyType: string(v.NotifyType),
				MsgType:    v.MsgType,
				MsgContent: content,
				MsgExtend:  extend,
				FileID:     v.FileID.Int64,
				AccountID:  v.AccountID.Int64,
				RelationID: v.RelationID,
				CreateAt:   v.CreateAt,
				IsRevoke:   v.IsRevoke,
				IsTop:      v.IsTop,
				IsPin:      v.IsPin,
				PinTime:    v.PinTime,
				ReadIds:    readIDs,
				ReplyCount: v.ReplyCount,
			},
			RlyMsg: rlyMsg,
		})
	}
	return &reply.ParamGetMsgsRelationIDAndTime{List: result, Total: data[0].Total}, nil
}

func (message) OfferMsgsByAccountIDAndTime(ctx *gin.Context, params model.OfferMsgsByAccountIDAndTime) (*reply.ParamOfferMsgsByAccountIDAndTime, errcode.Err) {
	data, err := dao.Database.DB.OfferMsgsByAccountIDAndTime(ctx, &db.OfferMsgsByAccountIDAndTimeParams{
		CreateAt:  params.LastTime,
		Limit:     params.Limit,
		Offset:    params.Offset,
		Accountid: params.AccountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamOfferMsgsByAccountIDAndTime{List: []*reply.ParamMsgInfoWithRlyAndHasRead{}}, nil
	}
	result := make([]*reply.ParamMsgInfoWithRlyAndHasRead, 0, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, err = model.JsonToExtend(v.MsgExtend)
			if err != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				continue
			}
		}
		var readIDs []int64
		if params.AccountID == v.AccountID.Int64 {
			readIDs = v.ReadIds
		}
		var rlyMsg *reply.ParamRlyMsg
		if v.RlyMsgID.Valid {
			rlyMsgInfo, myErr := GetMsgInfoByID(ctx, v.RlyMsgID.Int64)
			if myErr != nil {
				continue
			}
			var rlyContent string
			var rlyExtend *model.MsgExtend
			if !rlyMsgInfo.IsRevoke {
				rlyContent = rlyMsgInfo.MsgContent
				rlyExtend, err = model.JsonToExtend(rlyMsgInfo.MsgExtend)
				if err != nil {
					global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
					continue
				}
			}
			rlyMsg = &reply.ParamRlyMsg{
				MsgID:      rlyMsgInfo.ID,
				MsgType:    rlyMsgInfo.MsgType,
				MsgContent: rlyContent,
				MsgExtend:  rlyExtend,
				IsRevoked:  rlyMsgInfo.IsRevoke,
			}
		}
		result = append(result, &reply.ParamMsgInfoWithRlyAndHasRead{
			ParamMsgInfoWithRly: reply.ParamMsgInfoWithRly{
				ParamMsgInfo: reply.ParamMsgInfo{
					ID:         v.ID,
					NotifyType: string(v.NotifyType),
					MsgType:    v.MsgType,
					MsgContent: content,
					MsgExtend:  extend,
					FileID:     v.FileID.Int64,
					AccountID:  v.AccountID.Int64,
					RelationID: v.RelationID,
					CreateAt:   v.CreateAt,
					IsRevoke:   v.IsRevoke,
					IsTop:      v.IsTop,
					IsPin:      v.IsPin,
					PinTime:    v.PinTime,
					ReadIds:    readIDs,
					ReplyCount: v.ReplyCount,
				},
				RlyMsg: rlyMsg,
			},
			HasRead: false, //v.HasRead,
		})
	}
	return &reply.ParamOfferMsgsByAccountIDAndTime{List: result, Total: data[0].Total}, nil
}

func (message) UpdateMsgPin(ctx *gin.Context, accountID int64, params *request.ParamUpdateMsgPin) errcode.Err {
	ok, err := ExistsSetting(ctx, accountID, params.RelationID)
	if err != nil {
		return err
	}
	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	msgInfo, err := GetMsgInfoByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if msgInfo.IsPin == params.IsPin {
		return nil
	}
	myErr := dao.Database.DB.UpdateMsgPin(ctx, &db.UpdateMsgPinParams{
		ID:    params.ID,
		IsPin: params.IsPin,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送 pin 通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, params.RelationID, params.ID, server.MsgPin, params.IsPin))
	return nil
}

func (message) UpdateMsgTop(ctx *gin.Context, accountID int64, params *request.ParamUpdateMsgTop) errcode.Err {
	ok, err := ExistsSetting(ctx, accountID, params.RelationID)
	if err != nil {
		return err
	}
	if !ok {
		return errcodes.AuthPermissionsInsufficient
	}
	msgInfo, err := GetMsgInfoByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if msgInfo.IsTop == params.IsTop {
		return nil
	}
	myErr := dao.Database.DB.UpdateMsgTop(ctx, &db.UpdateMsgTopParams{
		ID:    params.ID,
		IsTop: params.IsTop,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送 置顶 消息
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, params.RelationID, params.ID, server.MsgTop, params.IsTop))
	return nil
}

func (message) RevokeMsg(ctx *gin.Context, accountID, msgID int64) errcode.Err {
	msgInfo, err := GetMsgInfoByID(ctx, msgID)
	if err != nil {
		return err
	}
	// 检查权限(是不是本人)
	if msgInfo.AccountID.Int64 != accountID {
		return errcodes.AuthPermissionsInsufficient
	}
	if msgInfo.IsRevoke {
		return errcodes.MsgAlreadyRevoke
	}
	myErr := dao.Database.DB.RevokeMsgWithTx(ctx, msgID, msgInfo.IsPin, msgInfo.IsTop)
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateMsgState(accessToken, msgInfo.RelationID, msgID, server.MsgRevoke, true))
	if msgInfo.IsTop {
		// 推送 top 通知
		global.Worker.SendTask(task.UpdateMsgState(accessToken, msgInfo.RelationID, msgID, server.MsgTop, false))
		// 创建并推送 top 消息
		f := func() error {
			arg := &db.CreateMessageParams{
				NotifyType: db.MsgnotifytypeSystem,
				MsgType:    string(model.MsgTypeText),
				MsgContent: fmt.Sprintf(format.UnTopMessage, accountID),
				MsgExtend:  pgtype.JSON{Status: pgtype.Null},
				RelationID: msgInfo.RelationID,
			}
			msgRly, err := dao.Database.DB.CreateMessage(ctx, arg)
			if err != nil {
				return err
			}
			global.Worker.SendTask(task.PublishMsg(accessToken, reply.ParamMsgInfo{
				ID:         msgRly.ID,
				NotifyType: string(arg.NotifyType),
				MsgType:    arg.MsgType,
				MsgContent: arg.MsgContent,
				RelationID: arg.RelationID,
				CreateAt:   msgRly.CreateAt,
			}, nil))
			return nil
		}
		if err := f(); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			reTry("UpdateMsgTop", f)
		}
	}
	return nil
}

func (message) GetTopMsgByRelationID(ctx *gin.Context, accountID, relationID int64) (*reply.ParamGetTopMsgByRelationID, errcode.Err) {
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errcodes.AuthPermissionsInsufficient
	}
	data, myErr := dao.Database.DB.GetTopMsgByRelationID(ctx, relationID)
	if myErr != nil {
		if errors.Is(myErr, pgx.ErrNoRows) {
			return nil, nil
		}
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	var content string
	var extend *model.MsgExtend
	if !data.IsRevoke {
		content = data.MsgContent
		extend, myErr = model.JsonToExtend(data.MsgExtend)
		if myErr != nil {
			global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
	}
	var readIDs []int64
	if accountID == data.AccountID.Int64 {
		readIDs = data.ReadIds
	}
	return &reply.ParamGetTopMsgByRelationID{MsgInfo: reply.ParamMsgInfo{
		ID:         data.ID,
		NotifyType: string(data.NotifyType),
		MsgType:    data.MsgType,
		MsgContent: content,
		MsgExtend:  extend,
		FileID:     data.FileID.Int64,
		AccountID:  data.AccountID.Int64,
		RelationID: data.RelationID,
		CreateAt:   data.CreateAt,
		IsRevoke:   data.IsRevoke,
		IsTop:      data.IsTop,
		IsPin:      data.IsPin,
		PinTime:    data.PinTime,
		ReadIds:    readIDs,
		ReplyCount: data.ReplyCount,
	}}, nil
}

func (message) GetPinMsgsByRelationID(ctx *gin.Context, accountID, relationID int64, limit, offset int32) (*reply.ParamGetPinMsgsByRelationID, errcode.Err) {
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, err
	}
	if !ok {
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcodes.AuthPermissionsInsufficient
	}
	data, myErr := dao.Database.DB.GetPinMsgsByRelationID(ctx, &db.GetPinMsgsByRelationIDParams{
		RelationID: relationID,
		Limit:      limit,
		Offset:     offset,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetPinMsgsByRelationID{List: []*reply.ParamMsgInfo{}}, nil
	}
	result := make([]*reply.ParamMsgInfo, 0, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, myErr = model.JsonToExtend(v.MsgExtend)
			if myErr != nil {
				global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
				return &reply.ParamGetPinMsgsByRelationID{Total: 0}, errcode.ErrServer
			}
		}
		var readIDs []int64
		if accountID == v.AccountID.Int64 {
			readIDs = v.ReadIds
		}
		result = append(result, &reply.ParamMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    v.MsgType,
			MsgContent: content,
			MsgExtend:  extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
			IsRevoke:   v.IsRevoke,
			IsTop:      v.IsTop,
			IsPin:      v.IsPin,
			PinTime:    v.PinTime,
			ReadIds:    readIDs,
			ReplyCount: v.ReplyCount,
		})
	}
	return &reply.ParamGetPinMsgsByRelationID{
		List:  result,
		Total: data[0].Total,
	}, nil
}

func (message) GetRlyMsgsInfoByMsgID(ctx *gin.Context, accountID, relationID, msgID int64, limit, offset int32) (*reply.ParamGetRlyMsgsInfoByMsgID, errcode.Err) {
	ok, err := ExistsSetting(ctx, accountID, relationID)
	if err != nil {
		return &reply.ParamGetRlyMsgsInfoByMsgID{}, err
	}
	if !ok {
		return &reply.ParamGetRlyMsgsInfoByMsgID{}, errcodes.AuthPermissionsInsufficient
	}
	data, myErr := dao.Database.DB.GetRlyMsgsInfoByMsgID(ctx, &db.GetRlyMsgsInfoByMsgIDParams{
		RelationID: relationID,
		Limit:      limit,
		Offset:     offset,
		RlyMsgID:   msgID,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetRlyMsgsInfoByMsgID{List: []*reply.ParamMsgInfo{}}, nil
	}
	result := make([]*reply.ParamMsgInfo, 0, len(data))
	for _, v := range data {
		var content string
		var extend *model.MsgExtend
		if !v.IsRevoke {
			content = v.MsgContent
			extend, myErr = model.JsonToExtend(v.MsgExtend)
			if myErr != nil {
				global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
				return &reply.ParamGetRlyMsgsInfoByMsgID{Total: 0}, errcode.ErrServer
			}
		}
		var readIDs []int64
		if accountID == v.AccountID.Int64 {
			readIDs = v.ReadIds
		}
		result = append(result, &reply.ParamMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    v.MsgType,
			MsgContent: content,
			MsgExtend:  extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
			IsRevoke:   v.IsRevoke,
			IsTop:      v.IsTop,
			IsPin:      v.IsPin,
			PinTime:    v.PinTime,
			ReadIds:    readIDs,
			ReplyCount: v.ReplyCount,
		})
	}
	return &reply.ParamGetRlyMsgsInfoByMsgID{
		List:  result,
		Total: data[0].Total,
	}, nil
}

// 从指定关系中模糊查找指定内容的信息
func getMsgsByContentAndRelation(ctx *gin.Context, params *db.GetMsgsByContentAndRelationParams) (*reply.ParamGetMsgsByContent, errcode.Err) {
	ok, err := ExistsSetting(ctx, params.AccountID, params.RelationID)
	if err != nil {
		return &reply.ParamGetMsgsByContent{}, err
	}
	if !ok {
		return &reply.ParamGetMsgsByContent{}, errcodes.AuthPermissionsInsufficient
	}
	data, myErr := dao.Database.DB.GetMsgsByContentAndRelation(ctx, params)
	if myErr != nil {
		global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetMsgsByContent{}, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetMsgsByContent{List: []*reply.ParamBriefMsgInfo{}}, nil
	}
	result := make([]*reply.ParamBriefMsgInfo, 0, len(data))
	for _, v := range data {
		var extend *model.MsgExtend
		extend, myErr = model.JsonToExtend(v.MsgExtend)
		if myErr != nil {
			global.Logger.Error(myErr.Error(), zap.Any("msgExtend", v.MsgExtend))
			continue
		}
		result = append(result, &reply.ParamBriefMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    v.MsgType,
			MsgContent: v.MsgContent,
			Extend:     extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
		})
	}
	return &reply.ParamGetMsgsByContent{List: result, Total: data[0].Total}, nil
}

func (message) GetMsgsByContent(ctx *gin.Context, accountID, relationID int64, content string, limit, offset int32) (*reply.ParamGetMsgsByContent, errcode.Err) {
	if relationID >= 0 {
		return getMsgsByContentAndRelation(ctx, &db.GetMsgsByContentAndRelationParams{
			RelationID: relationID,
			AccountID:  accountID,
			Limit:      limit,
			Offset:     offset,
			Content:    content,
		})
	}
	data, err := dao.Database.DB.GetMsgsByContent(ctx, &db.GetMsgsByContentParams{
		AccountID: accountID,
		Limit:     limit,
		Offset:    offset,
		Content:   content,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamGetMsgsByContent{}, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetMsgsByContent{List: []*reply.ParamBriefMsgInfo{}}, nil
	}
	result := make([]*reply.ParamBriefMsgInfo, 0, len(data))
	for _, v := range data {
		extend, myErr := model.JsonToExtend(v.MsgExtend)
		if myErr != nil {
			global.Logger.Error(myErr.Error(), zap.Any("msgExtend", v.MsgExtend))
			continue
		}
		result = append(result, &reply.ParamBriefMsgInfo{
			ID:         v.ID,
			NotifyType: string(v.NotifyType),
			MsgType:    v.MsgType,
			MsgContent: v.MsgContent,
			Extend:     extend,
			FileID:     v.FileID.Int64,
			AccountID:  v.AccountID.Int64,
			RelationID: v.RelationID,
			CreateAt:   v.CreateAt,
		})
	}
	return &reply.ParamGetMsgsByContent{List: result, Total: data[0].Total}, nil
}

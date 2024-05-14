package logic

import (
	"Chat/dao"
	db "Chat/dao/postgresql/sqlc"
	"Chat/errcodes"
	"Chat/global"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/reply"
	"Chat/model/request"
	"Chat/task"
	"database/sql"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
	"time"
)

type notify struct {
}

func (notify) CreateNotify(ctx *gin.Context, accountID int64, params *request.ParamCreateNotify) (*reply.ParamGroupNotify, errcode.Err) {
	ok, err := dao.Database.DB.ExistsIsLeader(ctx, &db.ExistsIsLeaderParams{
		RelationID: params.RelationID,
		AccountID:  accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !ok {
		return nil, errcodes.NotLeader
	}
	extend, _ := model.ExtendToJson(params.MsgExtend)
	data, err := dao.Database.DB.CreateGroupNotify(ctx, &db.CreateGroupNotifyParams{
		RelationID: sql.NullInt64{Int64: params.RelationID, Valid: true},
		MsgContent: params.MsgContent,
		MsgExpand:  extend,
		AccountID:  sql.NullInt64{Int64: accountID, Valid: true},
		CreateAt:   time.Now(),
		ReadIds:    []int64{accountID},
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	msgExtend, err := model.JsonToExtend(data.MsgExpand)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	// 推送创建群通知成功的消息
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.CreateNotify(accessToken, accountID, params.RelationID, data.MsgContent, msgExtend))
	return &reply.ParamGroupNotify{
		ID:         data.ID,
		RelationID: data.RelationID.Int64,
		MsgContent: data.MsgContent,
		MsgExtend:  msgExtend,
		AccountID:  data.AccountID.Int64,
		CreateAt:   data.CreateAt,
		ReadIDs:    data.ReadIds,
	}, nil
}

func (notify) UpdateNotify(ctx *gin.Context, accountID int64, params *request.ParamUpdateNotify) errcode.Err {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: params.RelationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcodes.NotGroupMember
	}
	extend, _ := model.ExtendToJson(params.MsgExtend)
	_, err = dao.Database.DB.UpdateGroupNotify(ctx, &db.UpdateGroupNotifyParams{
		RelationID: sql.NullInt64{Int64: params.RelationID, Valid: true},
		MsgContent: params.MsgContent,
		MsgExpand:  extend,
		AccountID:  sql.NullInt64{Int64: accountID, Valid: true},
		CreateAt:   time.Now(),
		ReadIds:    []int64{accountID},
		ID:         params.ID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送更改群通知的消息
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.UpdateNotify(accessToken, accountID, params.RelationID, params.MsgContent, params.MsgExtend))
	return nil
}

func (notify) GetNotifyByID(ctx *gin.Context, accountID, relationID int64) (*reply.ParamGetNotifyByID, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !ok {
		return nil, errcodes.NotGroupMember
	}
	data, err := dao.Database.DB.GetGroupNotifyByID(ctx, sql.NullInt64{Int64: relationID, Valid: true})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]reply.ParamGroupNotify, 0, len(data))
	for _, v := range data {
		msgExtend, err := model.JsonToExtend(v.MsgExpand)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
		result = append(result, reply.ParamGroupNotify{
			ID:         v.ID,
			RelationID: v.RelationID.Int64,
			MsgContent: v.MsgContent,
			MsgExtend:  msgExtend,
			AccountID:  v.AccountID.Int64,
			CreateAt:   v.CreateAt,
			ReadIDs:    v.ReadIds,
		})
	}
	return &reply.ParamGetNotifyByID{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (notify) DeleteNotify(ctx *gin.Context, accountID, id, relationID int64) errcode.Err {
	ok, err := dao.Database.DB.ExistsIsLeader(ctx, &db.ExistsIsLeaderParams{
		RelationID: relationID,
		AccountID:  accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcodes.NotLeader
	}
	err = dao.Database.DB.DeleteGroupNotify(ctx, id)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	return nil
}

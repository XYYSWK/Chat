package logic

import (
	"Chat/dao"
	db "Chat/dao/postgresql/sqlc"
	"Chat/errcodes"
	"Chat/global"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/reply"
	"Chat/task"
	"github.com/XYYSWK/Rutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type group struct {
}

func (group) CreateGroup(ctx *gin.Context, accountID int64, name string, description string) (relationID int64, err errcode.Err) {
	relationID, myErr := dao.Database.DB.CreateGroupRelation(ctx, &db.CreateGroupRelationParams{
		Name:        name,
		Description: description,
		Avatar:      global.PublicSetting.Rules.DefaultAvatarURL,
	})
	if myErr != nil {
		global.Logger.Error(myErr.Error())
		return 0, errcode.ErrServer
	}
	myErr = dao.Database.DB.AddSettingWithTx(ctx, dao.Database.Redis, accountID, relationID, true)
	if myErr != nil {
		global.Logger.Error(myErr.Error())
		return 0, errcode.ErrServer
	}
	return relationID, nil
}

func (group) TransferGroup(ctx *gin.Context, accountID, relationID, toAccountID int64) errcode.Err {
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
	ok, err = dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcodes.NotGroupMember
	}
	err = dao.Database.DB.TransferGroupWithTx(ctx, accountID, relationID, toAccountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送群主更换的通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.TransferGroup(accessToken, accountID, relationID))
	return nil
}

func (group) DissolveGroup(ctx *gin.Context, accountID, relationID int64) errcode.Err {
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
	err = dao.Database.DB.DeleteRelationWithTx(ctx, dao.Database.Redis, relationID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送群解散的消息
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.DissolveGroup(accessToken, relationID))
	return nil
}

func (group) UpdateGroup(ctx *gin.Context, accountID, relationID int64, name, description string) (*reply.ParamUpdateGroup, errcode.Err) {
	ok, err := dao.Database.DB.ExistsIsLeader(ctx, &db.ExistsIsLeaderParams{
		RelationID: relationID,
		AccountID:  accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !ok {
		return nil, errcodes.NotLeader
	}
	data, err := dao.Database.DB.GetGroupRelationByID(ctx, relationID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	err = dao.Database.DB.UpdateGroupRelation(ctx, &db.UpdateGroupRelationParams{
		Name:        name,
		Description: description,
		ID:          relationID,
		Avatar:      data.Avatar,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamUpdateGroup{
		Name:        name,
		Description: description,
	}, nil
}

func (group) InviteAccount(ctx *gin.Context, accountID, relationID int64, members []int64) (*reply.ParamInviteAccount, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		RelationID: relationID,
		AccountID:  accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !ok {
		return nil, errcodes.NotGroupMember
	}
	result := make([]int64, 0, len(members))
	for _, v := range members {
		ok1, err1 := dao.Database.DB.ExistsFriendSetting(ctx, &db.ExistsFriendSettingParams{
			Account1ID: accountID,
			Account2ID: v,
		})
		if err1 != nil {
			global.Logger.Error(err1.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
		if !ok1 {
			return nil, errcodes.RelationNotExists
		}
		ok2, err2 := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
			AccountID:  v,
			RelationID: relationID,
		})
		if err2 != nil {
			global.Logger.Error(err2.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcode.ErrServer
		}
		if !ok2 {
			err = dao.Database.DB.AddSettingWithTx(ctx, dao.Database.Redis, v, relationID, false)
			if err != nil {
				global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
				return nil, errcode.ErrServer
			}
			result = append(result, v)
		}
	}
	// 推送邀请进群的消息
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.InviteGroup(accessToken, accountID, relationID))

	return &reply.ParamInviteAccount{InviteMember: result}, nil
}

func (group) GetGroupList(ctx *gin.Context, accountID int64) (*reply.ParamGetGroupList, errcode.Err) {
	data, err := dao.Database.DB.GetGroupList(ctx, accountID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.SettingGroup, 0, len(data))
	for _, v := range data {
		result = append(result, &model.SettingGroup{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: "group",
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			GroupInfo: &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        v.GroupName.(string),
				Description: v.Discription.(string),
				Avatar:      v.GroupAvatar.(string),
			},
		})
	}
	return &reply.ParamGetGroupList{
		List:  result,
		Total: data[0].Total,
	}, nil
}

func (group) QuitGroup(ctx *gin.Context, accountID, relationID int64) errcode.Err {
	ok, err := dao.Database.DB.ExistsIsLeader(ctx, &db.ExistsIsLeaderParams{
		RelationID: relationID,
		AccountID:  accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if ok {
		return errcodes.IsLeader
	}
	ok, err = dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	if !ok {
		return errcodes.NotGroupMember
	}
	err = dao.Database.DB.DeleteSettingWithTx(ctx, dao.Database.Redis, accountID, relationID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送退群通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.QuitGroup(accessToken, accountID, relationID))
	return nil
}

func (group) GetGroupsByName(ctx *gin.Context, accountID int64, name string, limit, offset int32) (*reply.ParamGetGroupsByName, errcode.Err) {
	data, err := dao.Database.DB.GetGroupSettingsByName(ctx, &db.GetGroupSettingsByNameParams{
		AccountID: accountID,
		Limit:     limit,
		Offset:    offset,
		Name:      name,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]model.SettingGroup, 0, len(data))
	for _, v := range data {
		result = append(result, model.SettingGroup{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: "group",
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			GroupInfo: &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        v.GroupName.(string),
				Description: v.Description.(string),
				Avatar:      v.GroupAvatar.(string),
			},
		})
	}
	return &reply.ParamGetGroupsByName{
		List:  result,
		Total: data[0].Total,
	}, nil
}

func (group) GetGroupMembers(ctx *gin.Context, accountID, relationID int64, limit, offset int32) (*reply.ParamGetGroupMembers, errcode.Err) {
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
	data, err := dao.Database.DB.GetGroupMembersByID(ctx, &db.GetGroupMembersByIDParams{
		RelationID: relationID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]reply.ParamGroupMemberInfo, 0, len(data))
	for _, v := range data {
		result = append(result, reply.ParamGroupMemberInfo{
			AccountID: v.ID,
			Name:      v.Name,
			Avatar:    v.Avatar,
			Nickname:  v.NickName.String,
			IsLeader:  v.IsLeader.Bool,
		})
	}
	return &reply.ParamGetGroupMembers{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

package logic

import (
	"Chat/dao"
	db "Chat/dao/postgresql/sqlc"
	"Chat/errcodes"
	"Chat/global"
	"Chat/middlewares"
	"Chat/model"
	"Chat/model/reply"
	"Chat/pkg/gtype"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	upload "github.com/XYYSWK/Lutils/pkg/upload/obs"
	"github.com/XYYSWK/Lutils/pkg/upload/obs/huawei_cloud"
	"github.com/gin-gonic/gin"
	obs2 "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"mime/multipart"
)

type file struct {
}

// PublishFile 上传文件，传出 context 与 relationID，accountID，file(*multipart.FileHeader)，返回 model.PublishFileRe
func (file) PublishFile(ctx *gin.Context, params model.PublishFile) (model.PublishFileReply, errcode.Err) {
	fileType, myErr := gtype.GetFileType(params.File)
	if myErr != nil {
		return model.PublishFileReply{}, errcode.ErrServer
	}
	if fileType == "file" {
		if params.File.Size > global.PublicSetting.Rules.BiggestFileSize {
			return model.PublishFileReply{}, errcodes.FileTooBig
		}
	} else {
		fileType = "img"
	}
	input := new(obs2.PutObjectInput)
	url, key, err := global.OBS.UploadFile(params.File, input)
	if err != nil {
		fmt.Println("-------------------------------------------------------------------------------")
		global.Logger.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return model.PublishFileReply{}, errcode.ErrServer
	}
	r, err := dao.Database.DB.CreateFile(ctx, &db.CreateFileParams{
		FileName: params.File.Filename,
		FileType: db.Filetype(fileType),
		FileSize: params.File.Size,
		Key:      key,
		Url:      url,
		RelationID: sql.NullInt64{
			Int64: params.RelationID,
			Valid: true,
		},
		AccountID: sql.NullInt64{
			Int64: params.AccountID,
			Valid: true,
		},
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return model.PublishFileReply{}, errcode.ErrServer
	}
	return model.PublishFileReply{
		ID:       r.ID,
		FileType: fileType,
		FileSize: r.FileSize,
		Url:      r.Url,
		CreateAt: r.CreateAt,
	}, nil
}

// DeleteFile 删除文件
func (file) DeleteFile(ctx context.Context, fileID int64) errcode.Err {
	key, myErr := dao.Database.DB.GetFileKeyByID(ctx, fileID)
	if myErr != nil {
		if errors.Is(myErr, sql.ErrNoRows) {
			return errcodes.FileNotExist
		}
		global.Logger.Error(myErr.Error())
		return errcode.ErrServer
	}
	if key != "" {
		_, err := global.OBS.DeleteFile(key)
		if err != nil {
			global.Logger.Error(err.Error())
			return errcodes.FileDeleteFailed
		}
	}
	err := dao.Database.DB.DeleteFileByID(ctx, fileID)
	if err != nil {
		global.Logger.Error(err.Error())
		return errcode.ErrServer
	}
	return nil
}

func (file) GetRelationFile(ctx *gin.Context, relationID int64) (*reply.ParamGetRelationFile, errcode.Err) {
	list, err := dao.Database.DB.GetFileByRelationID(ctx, sql.NullInt64{Int64: relationID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.FileNotExist
		}
	}
	data := make([]*reply.ParamFile, len(list))
	for i, f := range list {
		data[i] = &reply.ParamFile{
			FileID:    f.ID,
			FileName:  f.FileName,
			FileType:  string(f.FileType),
			FileSize:  f.FileSize,
			Url:       f.Url,
			AccountID: f.AccountID.Int64,
			CreateAt:  f.CreateAt,
		}
	}
	return &reply.ParamGetRelationFile{FileList: data}, nil
}

// UploadAccountAvatar 更新账户头像
func (file) UploadAccountAvatar(ctx *gin.Context, accountID int64, fileInfo *multipart.FileHeader) (*reply.ParamUploadAvatar, errcode.Err) {
	relationID, err := dao.Database.DB.GetRelationIDByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.RelationNotExists
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	exists, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if !exists {
		return nil, errcodes.AuthenticationFailed
	}
	var url string
	if fileInfo != nil {
		input := new(obs2.PutObjectInput)
		obs := initOBS(huawei_cloud.AccountAvatarType)
		url, _, err = obs.UploadFile(fileInfo, input)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return nil, errcodes.FailedStore
		}
	}
	err = dao.Database.DB.UpdateAccountAvatar(ctx, &db.UpdateAccountAvatarParams{
		Avatar: url,
		ID:     accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if fileInfo == nil {
		return &reply.ParamUploadAvatar{URL: global.PublicSetting.Rules.DefaultAvatarURL}, nil
	}
	return &reply.ParamUploadAvatar{URL: url}, nil
}

func initOBS(avatarType string) upload.OBS {
	if avatarType == huawei_cloud.AccountAvatarType {
		return huawei_cloud.Init(huawei_cloud.Config{
			Location:         global.PrivateSetting.HuaWeiOBS.Location,
			BucketName:       global.PrivateSetting.HuaWeiOBS.BucketName,
			BucketUrl:        global.PrivateSetting.HuaWeiOBS.BucketUrl,
			Endpoint:         global.PrivateSetting.HuaWeiOBS.Endpoint,
			BasePath:         global.PrivateSetting.HuaWeiOBS.BasePath,
			AvatarType:       huawei_cloud.AccountAvatarType,
			AccountAvatarUrl: global.PrivateSetting.HuaWeiOBS.AccountAvatarUrl,
			GroupAvatarUrl:   global.PrivateSetting.HuaWeiOBS.GroupAvatarUrl,
		})
	} else if avatarType == huawei_cloud.GroupAvatarType {
		return huawei_cloud.Init(huawei_cloud.Config{
			Location:         global.PrivateSetting.HuaWeiOBS.Location,
			BucketName:       global.PrivateSetting.HuaWeiOBS.BucketName,
			BucketUrl:        global.PrivateSetting.HuaWeiOBS.BucketUrl,
			Endpoint:         global.PrivateSetting.HuaWeiOBS.Endpoint,
			BasePath:         global.PrivateSetting.HuaWeiOBS.BasePath,
			AvatarType:       huawei_cloud.GroupAvatarType,
			AccountAvatarUrl: global.PrivateSetting.HuaWeiOBS.AccountAvatarUrl,
			GroupAvatarUrl:   global.PrivateSetting.HuaWeiOBS.GroupAvatarUrl,
		})
	}
	return global.OBS
}

func (file) GetFileDetailsByID(ctx *gin.Context, fileID int64) (*reply.ParamFile, errcode.Err) {
	result, err := dao.Database.DB.GetFileDetailsByID(ctx, fileID)
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return &reply.ParamFile{
		FileID:    result.ID,
		FileName:  result.FileName,
		FileType:  string(result.FileType),
		FileSize:  result.FileSize,
		Url:       result.Url,
		AccountID: result.AccountID.Int64,
		CreateAt:  result.CreateAt,
	}, nil
}

func (file) UploadGroupAvatar(ctx *gin.Context, file *multipart.FileHeader, accountID, relationID int64) (*reply.ParamUploadAvatar, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamUploadAvatar{URL: ""}, errcode.ErrServer
	}
	if !ok {
		return &reply.ParamUploadAvatar{URL: ""}, errcodes.NotGroupMember
	}
	obs := initOBS(huawei_cloud.GroupAvatarType)
	var url, key string
	input := new(obs2.PutObjectInput)
	if file != nil {
		url, key, err = obs.UploadFile(file, input)
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return &reply.ParamUploadAvatar{URL: ""}, errcode.ErrServer
		}
	}
	if file == nil {
		url = global.PublicSetting.Rules.DefaultAvatarURL
	}
	err = dao.Database.DB.UploadGroupAvatarWithTx(ctx, db.CreateFileParams{
		FileName:   "groupAvatar",
		FileType:   "",
		FileSize:   0,
		Key:        key,
		Url:        url,
		RelationID: sql.NullInt64{Int64: relationID, Valid: true},
		AccountID:  sql.NullInt64{},
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return &reply.ParamUploadAvatar{URL: ""}, errcode.ErrServer
	}
	if file == nil {
		return &reply.ParamUploadAvatar{URL: global.PublicSetting.Rules.DefaultAvatarURL}, nil
	}
	return &reply.ParamUploadAvatar{URL: url}, nil
}

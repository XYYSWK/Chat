package request

import "mime/multipart"

type ParamPublishFile struct {
	File       *multipart.FileHeader `form:"file" binding:"required" swaggerignore:"true"`
	RelationID int64                 `form:"relation_id" binding:"required"`
	AccountID  int64                 `form:"account_id" binding:"required"`
}

type ParamDeleteFile struct {
	FileID int64 `json:"file_id" form:"file_id" binding:"required"`
}

type ParamGetRelationFile struct {
	RelationID int64 `json:"relation_id" form:"relation_id" binding:"required"`
}

type ParamUploadAccountAvatar struct {
	File *multipart.FileHeader `form:"file" binding:"required" swaggerignore:"true"`
}

type ParamGetFileDetailsByID struct {
	FileID int64 `json:"file_id" form:"file_id" binding:"required"`
}

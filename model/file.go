package model

import (
	"mime/multipart"
	"time"
)

type PublishFile struct {
	File       *multipart.FileHeader `form:"file" binding:"required" swaggerignore:"true"`
	RelationID int64                 `form:"relation_id" binding:"required"`
	AccountID  int64                 `form:"account_id" binding:"required"`
}

type PublishFileReply struct {
	ID       int64     `json:"id,omitempty"`
	FileType string    `json:"file_type,omitempty"`
	FileSize int64     `json:"file_size,omitempty"`
	Url      string    `json:"url,omitempty"`
	CreateAt time.Time `json:"create_at"`
}

package reply

import "time"

type ParamPublishFile struct {
	ID       int64     `json:"id,omitempty"`
	FileType string    `json:"file_type,omitempty"`
	FileSize int64     `json:"file_size,omitempty"`
	Url      string    `json:"url,omitempty"`
	CreateAt time.Time `json:"create_at"`
}

type ParamFile struct {
	FileID    int64     `json:"file_id,omitempty"`
	FileName  string    `json:"file_name,omitempty"`
	FileType  string    `json:"file_type,omitempty"`
	FileSize  int64     `json:"file_size,omitempty"`
	Url       string    `json:"url,omitempty"`
	AccountID int64     `json:"account_id,omitempty"`
	CreateAt  time.Time `json:"create_at"`
}

type ParamGetRelationFile struct {
	FileList []*ParamFile `json:"file_list,omitempty"`
}

type ParamUploadAvatar struct {
	URL string `json:"url,omitempty"`
}

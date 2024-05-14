package settings

import (
	"Chat/global"
	"github.com/XYYSWK/Rutils/pkg/upload/obs/huawei_cloud"
)

type obs struct {
}

func (obs) Init() {
	global.OBS = huawei_cloud.Init(huawei_cloud.Config{
		Location:         global.PrivateSetting.HuaWeiOBS.Location,
		BucketName:       global.PrivateSetting.HuaWeiOBS.BucketName,
		BucketUrl:        global.PrivateSetting.HuaWeiOBS.BucketUrl,
		Endpoint:         global.PrivateSetting.HuaWeiOBS.Endpoint,
		BasePath:         global.PrivateSetting.HuaWeiOBS.BasePath,
		AvatarType:       global.PrivateSetting.HuaWeiOBS.AvatarType,
		AccountAvatarUrl: global.PrivateSetting.HuaWeiOBS.AccountAvatarUrl,
		GroupAvatarUrl:   global.PrivateSetting.HuaWeiOBS.GroupAvatarUrl,
	})
}

package settings

import (
	"Chat/dao"
	"Chat/global"
	"Chat/pkg/tool"
	"context"
)

/*
所有需要再启动时初始化的配置
*/

type load struct {
}

func (load) Init() {
	var err error
	// 加载所有邮箱到 redis
	err = tool.DoThat(err, LoadAllEmailsToRedis)
	// 加载所有关系名单到 redis
	err = tool.DoThat(err, LoadAllGroupRelationToRedis)
	if err != nil {
		panic(err)
	}
}

// LoadAllEmailsToRedis 加载所有邮箱到 redis
func LoadAllEmailsToRedis() error {
	emails, err := dao.Database.DB.GetAllEmail(context.Background())
	if err != nil {
		return err
	}
	err = dao.Database.Redis.ReloadEmails(context.Background(), emails...)
	if err != nil {
		return err
	}
	global.Logger.Info("邮箱加载完成")
	return nil
}

// LoadAllGroupRelationToRedis 加载所有关系名单到 redis
func LoadAllGroupRelationToRedis() error {
	// 群 ID 和成员 IDs
	relations := make(map[int64][]int64)
	relationIDs, err := dao.Database.DB.GetAllRelationIDs(context.Background())
	if err != nil {
		return err
	}
	for _, relationID := range relationIDs {
		accountIDs, err := dao.Database.DB.GetAccountIDsByRelationID(context.Background(), relationID)
		if err != nil {
			return err
		}
		relations[relationID] = accountIDs
	}
	err = dao.Database.Redis.ReloadRelationIDs(context.Background(), relations)
	if err != nil {
		return err
	}
	global.Logger.Info("关系名单加载完成")
	return nil
}

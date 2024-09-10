package logic

import (
	"Chat/dao"
	"Chat/global"
	"context"
	"github.com/XYYSWK/Lutils/pkg/goroutine/task"
)

type auto struct {
}

const taskName = "deleteExpiredFile"

// Work 初始化并启动一个删除过期文件的任务
func (auto) Work() {
	ctx := context.Background()
	deleteExpiredFileTask := task.Task{
		Name:            taskName,
		Ctx:             ctx,
		TaskDuration:    global.PublicSetting.Auto.DeleteExpiredFileDuration,
		TimeoutDuration: global.PublicSetting.Server.DefaultContextTimeout,
		F:               DeleteExpiredFile(),
	}
	startTask(deleteExpiredFileTask)
}

// startTask 启动多个任务，并将它们转化为定时任务以供执行
func startTask(tasks ...task.Task) {
	for i := range tasks {
		task.NewTickerTask(tasks[i])
	}
}

// DeleteExpiredFile 定时删除没有 relation 的文件
func DeleteExpiredFile() task.DoFunc {
	return func(parentCtx context.Context) {
		global.Logger.Info("auto task run: deleteExpiredFile")
		ctx, cancel := context.WithTimeout(parentCtx, global.PublicSetting.Server.DefaultContextTimeout)
		defer cancel()
		data, myErr := dao.Database.DB.GetFileByRelationIDIsNULL(ctx)
		if myErr != nil {
			global.Logger.Error(myErr.Error())
			return
		}
		for _, v := range data {
			err := Logics.File.DeleteFile(ctx, v.ID)
			if err != nil {
				global.Logger.Error(err.Error())
				return
			}
		}
	}
}

package settings

import (
	"Chat/global"
	"github.com/XYYSWK/Lutils/pkg/goroutine/work"
)

type worker struct {
}

func (worker) Init() {
	global.Worker = work.Init(work.Config{
		TaskChanCapacity:   global.PublicSetting.Worker.TaskChanCapacity,
		WorkerChanCapacity: global.PublicSetting.Worker.WorkerChanCapacity,
		WorkerNum:          global.PublicSetting.Worker.WorkerNum,
	})
}

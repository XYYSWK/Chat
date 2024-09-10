package settings

import (
	"Chat/global"
	"github.com/XYYSWK/Lutils/pkg/logger"
)

type log struct {
}

// Init 初始化日志
func (log) Init() {
	global.Logger = logger.NewLogger(&logger.InitStruct{
		LogSavePath:   global.PublicSetting.Log.LogSavePath,
		LogFileExt:    global.PublicSetting.Log.LogFileExt,
		MaxSize:       global.PublicSetting.Log.MaxSize,
		MaxBackups:    global.PublicSetting.Log.MaxBackups,
		MaxAge:        global.PublicSetting.Log.MaxAge,
		Compress:      global.PublicSetting.Log.Compress,
		LowLevelFile:  global.PublicSetting.Log.LowLevelFile,
		HighLevelFile: global.PublicSetting.Log.HighLevelFile,
	}, global.PublicSetting.Log.Level)
}

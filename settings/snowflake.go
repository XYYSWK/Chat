package settings

import (
	"Chat/global"
	"github.com/XYYSWK/Lutils/pkg/generateID/snowflake"
	"time"
)

type generateID struct {
}

func (generateID) Init() {
	var err error
	global.GenerateID, err = snowflake.Init(time.Now(), global.PublicSetting.App.MachineID)
	if err != nil {
		panic(err)
	}
}

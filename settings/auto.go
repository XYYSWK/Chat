package settings

import "Chat/logic"

type auto struct {
}

func (auto) Init() {
	logic.Logics.Auto.Work()
}

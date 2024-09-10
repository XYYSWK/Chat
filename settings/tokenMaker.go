package settings

import (
	"Chat/global"
	"github.com/XYYSWK/Lutils/pkg/token"
)

type tokenMaker struct {
}

func (tokenMaker) Init() {
	var err error
	global.TokenMaker, err = token.NewPasetoMaker([]byte(global.PrivateSetting.Token.Key))
	if err != nil {
		panic(err)
	}
}

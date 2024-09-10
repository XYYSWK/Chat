package logic

import (
	"Chat/global"
	"Chat/model"
	"Chat/pkg/retry"
	"github.com/XYYSWK/Lutils/pkg/token"
)

// 尝试重试
// 失败：打印日志
func reTry(name string, f func() error) {
	go func() {
		report := <-retry.NewTry(name, f, global.PublicSetting.Auto.Retry.Duration, global.PublicSetting.Auto.Retry.MaxTimes).Run()
		global.Logger.Error(report.Error())
	}()
}

// newToken token
// 成功：返回 token，*token.Payload
// 失败：返回 nil, error
func newToken(t model.TokenType, id int64) (string, *token.Payload, error) {
	duration := global.PrivateSetting.Token.UserTokenDuration
	if t == model.AccountToken {
		duration = global.PrivateSetting.Token.AccountTokenDuration
	}
	data, err := model.NewTokenContent(t, id).Marshal()
	if err != nil {
		return "", nil, err
	}
	result, payload, err := global.TokenMaker.CreateToken(data, duration)
	if err != nil {
		return "", nil, err
	}
	return result, payload, nil
}

// 将 id 从小到大排序返回
func sortID(id1, id2 int64) (_, _ int64) {
	if id1 > id2 {
		return id2, id1
	}
	return id1, id2
}

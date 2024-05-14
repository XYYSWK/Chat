package global

import "context"

// DefaultContextWithTimeout 获取默认时间限制连接上下文
// 返回：上下文，终止函数
func DefaultContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), PublicSetting.Server.DefaultContextTimeout)
}

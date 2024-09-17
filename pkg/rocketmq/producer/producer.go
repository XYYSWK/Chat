package producer

import (
	"Chat/global"
	"Chat/model/reply"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// SendMsgToMQ 创建一个生产者，通过 RocketMQ 发送一条系统消息，告知指定用户他们已加入聊天室。
func SendMsgToMQ(mID int64, msg reply.ParamMsgInfoWithRly) {
	// 创建一个 RocketMQ 生产者，连接到指定的 RocketMQ 服务器
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.PrivateSetting.RocketMQ.Addr, global.PrivateSetting.RocketMQ.Port)}))
	if err != nil {
		// 如果创建生产者失败，立刻抛出异常，并停止程序
		panic(fmt.Sprintf("生成 Producer 失败: %s", err))
	}
	// 启动生产者，如果生产者启动失败，则抛出异常并停止程序
	if err := p.Start(); err != nil {
		panic(err)
	}
	// 构建消息的唯一标识符（UID），用于识别和定位该消息的接收者
	uID := fmt.Sprintf("accountID:%d", mID)
	// 消息内容
	sendMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化消息失败", err)
		return
	}
	// 使用 SendSync 方法发送消息，uID 是消息的主题（topic），SystemMsg 是消息的内容
	// primitive.NewMessage 将消息打包成 RocketMQ 消息
	res, err := p.SendSync(context.Background(), primitive.NewMessage(uID, sendMsg))
	if err != nil {
		fmt.Println("发送失败", err)
	} else {
		fmt.Println("发送成功， res：", res.String())
	}
	// 关闭生产者，如果关闭时发生错误，则抛出异常
	if err := p.Shutdown(); err != nil {
		panic(err)
	}
}

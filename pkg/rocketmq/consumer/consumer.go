package consumer

import (
	"Chat/global"
	"Chat/model/chat"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func StartConsumer(accountID int64) {
	// 使用给定的 accountID 创建一个唯一的主题或标签 uID
	uID := fmt.Sprintf("accountID:%d", accountID)
	// 创建一个新的 RocketMQ 推送消费者对象
	// consumer.WithNameServer 用于指定 RocketMQ NameServer 的地址
	// consumer.WithGroupName 指定消费者的消费组名称，这里使用 uID 作为组名称
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.PrivateSetting.RocketMQ.Addr, global.PrivateSetting.RocketMQ.Port)}),
		consumer.WithGroupName(uID),
	)
	if err != nil {
		fmt.Println("创建消费者失败：", err)
		return
	}

	// 订阅消息，uID 是订阅的主题(topic)，consumer.MessageSelector{}
	if err := c.Subscribe(uID, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			// 处理每条收到的消息，将消息发送到 global.ChatMap 中的指定客户端
			global.ChatMap.Send(accountID, chat.ClientSendMsg, msgs[i])
			// 打印接收到的消息
			fmt.Println("获取到值:", msgs[i])
		}
		// 返回处理成功的结果
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		fmt.Println("订阅消息失败:", err)
		return
	}
	// 启动消费者
	if err := c.Start(); err != nil {
		fmt.Println("启动消费者失败:", err)
		return
	}

	// 使用 defer 语句在函数退出时关闭消费者，确保资源得到释放
	defer func() {
		if err := c.Shutdown(); err != nil {
			fmt.Println("关闭消费者失败:", err)
		} else {
			fmt.Println("消费者已关闭！")
		}
	}()

	// 阻塞当前协程，保持消费者运行状态。通过 select{} 实现一个永远不会退出的循环
	select {}
}

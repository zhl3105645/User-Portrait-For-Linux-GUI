package main

import (
	"backend/biz/hadoop"
	"backend/biz/mq"
	"backend/cmd/dal"
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/gopkg/util/logger"
	"time"
)

func main() {
	// init
	dal.Init() // sql 初始化
	hadoop.Init(context.Background())
	InitConsumer() // consumer 初始化

	// 空转
	for {
		time.Sleep(time.Second * 60)
	}
}

func InitConsumer() {
	// 订阅主题、消费
	endPoint := []string{mq.EndPoint}
	// 创建一个consumer实例
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer(endPoint),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("ConsumerGroupName"),
	)

	// 订阅topic
	err = c.Subscribe(mq.TopicName, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			go handleMsg(msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		logger.Error("subscribe message error: ", err.Error())
	}

	// 启动consumer
	err = c.Start()

	if err != nil {
		logger.Error("consumer start error: ", err.Error())
		//panic(any("init mq consumer failed."))
	}
}

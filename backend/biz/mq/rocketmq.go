package mq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var pro rocketmq.Producer

const (
	EndPoint  = "127.0.0.1:9876"
	TopicName = "profile"
)

type Type int

const (
	RuleGene  Type = 3 // 生成规则数据
	LabelGene Type = 5 // 生成标签数据
)

type GeneMsg struct {
	AppId int64 // 应用ID
	Type  Type  // 消息类型
	Param int64 // 具体参数
}

func Init() {
	createTopicAndProducer()
}

func createTopicAndProducer() {
	endPoint := []string{EndPoint}
	// 创建主题
	ad, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(endPoint)))
	if err != nil {
		fmt.Printf("connection error: %s\n", err.Error())
	}
	err = ad.CreateTopic(context.Background(), admin.WithTopicCreate(TopicName))
	if err != nil {
		fmt.Printf("createTopic error: %s\n", err.Error())
	}

	// 创建一个producer实例
	pro, _ = rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
		producer.WithRetry(2),
		producer.WithGroupName("ProducerGroupName"),
	)
	// 启动
	err = pro.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		//panic(any("init mq producer failed."))
	}
}

func SendSyncMessage(message string) error {
	// 发送消息
	result, err := pro.SendSync(context.Background(), &primitive.Message{
		Topic: TopicName,
		Body:  []byte(message),
	})

	if err != nil {
		fmt.Printf("send message error: %s\n", err.Error())
	} else {
		fmt.Printf("send message seccess: result=%s\n", result.String())
	}

	return err
}

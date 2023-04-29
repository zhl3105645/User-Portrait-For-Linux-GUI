package main

import (
	"backend/biz/mq"
	"backend/cmd/dal"
	"backend/consumer/basic_behavior_gene"
	"backend/consumer/component_gene"
	"backend/consumer/config"
	"backend/consumer/label_gene"
	"backend/consumer/model_gene"
	"backend/consumer/rule_gene"
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/gopkg/util/logger"
	"time"
)

func main() {
	// init
	dal.Init()     // sql 初始化
	InitConsumer() // consumer 初始化
	// oldSolution() // 简单办法:使用文件通信

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
			logger.Info("subscribe callback : ", msgs[i])
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

func oldSolution() {
	config.StatusChan = make(chan *config.StatusChange, 10)

	// 开始
	go config.ReceiveStatusChange()

	for {
		configs, err := config.ReadConfig()
		if err != nil {
			logger.Error("read file err, err=%s", err.Error())
			continue
		}
		s, _ := json.Marshal(configs)
		logger.Info("Configs=", string(s))

		for appId, appConfig := range configs.Configs {
			if appConfig == nil {
				continue
			}

			for taskTyp, status := range appConfig.Config {
				if status != config.Begin {
					continue
				}

				// begin -> running
				config.StatusChan <- &config.StatusChange{
					AppId:    appId,
					TaskType: taskTyp,
					Status:   config.Running,
				}

				// 开始
				switch taskTyp {
				case config.ComponentGene:
					go component_gene.Gene(appId)
				case config.BasicBehaviorGene:
					go basic_behavior_gene.Gene(appId)
				case config.RuleGene:
					go rule_gene.Gene(appId)
				case config.ModelGene:
					modelId := appConfig.Param[config.ModelGene]
					go model_gene.Gene(appId, modelId)
				case config.LabelGene:
					labelId := appConfig.Param[config.LabelGene]
					go label_gene.Gene(appId, labelId)
				default:
				}
			}
		}

		// 10s 读取一次
		time.Sleep(time.Second * 10)
	}
}

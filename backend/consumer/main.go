package main

import (
	"backend/cmd/dal"
	"backend/consumer/basic_behavior_gene"
	"backend/consumer/component_gene"
	"backend/consumer/config"
	"backend/consumer/label_gene"
	"backend/consumer/model_gene"
	"backend/consumer/rule_gene"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"time"
)

func main() {
	// init
	dal.Init() // sql 初始化
	config.StatusChan = make(chan *config.StatusChange, 10)

	// 开始
	go config.ReceiveStatusChange()

	// 测试写入
	test()

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

func test() {
	label_gene.Gene(2, label_gene.GitFre)
}

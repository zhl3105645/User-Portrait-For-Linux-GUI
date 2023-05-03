package main

import (
	"backend/biz/mq"
	"backend/consumer/crowd_gene"
	"backend/consumer/label_gene"
	"backend/consumer/rule_gene"
	"backend/consumer/seq_mining"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/gopkg/util/logger"
)

func handleMsg(msg *primitive.MessageExt) {
	body := msg.Body
	logger.Info("body=", string(body))
	param := &mq.GeneMsg{}
	err := json.Unmarshal(body, param)
	if err != nil {
		logger.Error("json unmarshal failed. err=", err.Error())
		return
	}

	if param == nil {
		return
	}

	switch param.Type {
	case mq.RuleGene:
		go rule_gene.Gene(param.AppId)
	case mq.LabelGene:
		go label_gene.Gene(param.AppId, param.Param)
	case mq.CrowdGene:
		go crowd_gene.Gene(param.AppId, param.Param)
	case mq.SeqMining:
		go seq_mining.Gene(param.AppId, param.Param)
	default:
	}

	return
}

package rule

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/mq"
	"context"
	"encoding/json"
)

type GeneRule struct {
	accountId int64

	//
	appId int64
}

func NewGeneRule(accountId int64) *GeneRule {
	return &GeneRule{
		accountId: accountId,
	}
}

func (p *GeneRule) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	msg := &mq.GeneMsg{
		AppId: p.appId,
		Type:  mq.RuleGene,
		Param: 0,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return microtype.JsonMarshalFailed
	}

	err = mq.SendSyncMessage(string(msgJson))
	if err != nil {
		return microtype.MQSendFailed
	}

	//conf, err := config.ReadConfig()
	//if err != nil {
	//	return err
	//}
	//
	//// app_id typ=RuleGene status=begin
	//if len(conf.Configs) == 0 {
	//	conf.Configs = make(map[int64]*config.Config)
	//	conf.Configs[p.appId] = &config.Config{}
	//}
	//
	//if conf.Configs[p.appId] == nil || len(conf.Configs[p.appId].Config) == 0 {
	//	conf.Configs[p.appId] = &config.Config{
	//		Config: make(map[config.TaskType]config.Status),
	//	}
	//}
	//
	//status := conf.Configs[p.appId].Config[config.RuleGene]
	//if status != config.Stop {
	//	return microtype.RuleGene
	//}
	//
	//conf.Configs[p.appId].Config[config.RuleGene] = config.Begin
	//
	//err = config.WriteConfig(conf)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (p *GeneRule) GetResp() *backend.GeneResp {
	return &backend.GeneResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

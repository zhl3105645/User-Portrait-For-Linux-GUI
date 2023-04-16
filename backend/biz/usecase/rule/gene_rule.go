package rule

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/consumer/config"
	"context"
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

	conf, err := config.ReadConfig()
	if err != nil {
		return err
	}

	// app_id typ=RuleGene status=begin
	if len(conf.AppConfigs) == 0 {
		conf.AppConfigs = make(map[int64]*config.AppConfig)
		conf.AppConfigs[p.appId] = &config.AppConfig{}
	}

	if conf.AppConfigs[p.appId] == nil || len(conf.AppConfigs[p.appId].Config) == 0 {
		conf.AppConfigs[p.appId] = &config.AppConfig{
			Config: make(map[config.TaskType]config.Status),
		}
	}

	status := conf.AppConfigs[p.appId].Config[config.RuleGene]
	if status != config.Stop {
		return microtype.RuleGene
	}

	conf.AppConfigs[p.appId].Config[config.RuleGene] = config.Begin

	err = config.WriteConfig(conf)
	if err != nil {
		return err
	}

	return nil
}

func (p *GeneRule) GetResp() *backend.GeneResp {
	return &backend.GeneResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

package component

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/consumer/config"
	"context"
)

type GeneComponent struct {
	accountId int64

	//
	appId int64
}

func NewGeneComponent(accountId int64) *GeneComponent {
	return &GeneComponent{
		accountId: accountId,
	}
}

func (p *GeneComponent) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	conf, err := config.ReadConfig()
	if err != nil {
		return err
	}

	// app_id typ=gene_component status=begin
	if len(conf.AppConfigs) == 0 {
		conf.AppConfigs = make(map[int64]*config.AppConfig)
	}
	conf.AppConfigs[p.appId] = &config.AppConfig{}

	if len(conf.AppConfigs[p.appId].Config) == 0 {
		conf.AppConfigs[p.appId].Config = make(map[config.TaskType]config.Status)
	}
	conf.AppConfigs[p.appId].Config[config.ComponentGene] = config.Begin

	err = config.WriteConfig(conf)
	if err != nil {
		return err
	}

	return nil
}

func (p *GeneComponent) GetResp() *backend.GeneComponentResp {
	return &backend.GeneComponentResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

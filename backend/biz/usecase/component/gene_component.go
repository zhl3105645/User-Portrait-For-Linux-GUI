package component

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
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
	//ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	//if err := ac.Load(ctx); err != nil {
	//	return err
	//}
	//
	//p.appId = ac.GetQueryAccount().AppID
	//
	//conf, err := config.ReadConfig()
	//if err != nil {
	//	return err
	//}
	//
	//// app_id typ=gene_component status=begin
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
	//status := conf.Configs[p.appId].Config[config.ComponentGene]
	//if status != config.Stop {
	//	return microtype.ComponentInGene
	//}
	//
	//conf.Configs[p.appId].Config[config.ComponentGene] = config.Begin
	//
	//err = config.WriteConfig(conf)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (p *GeneComponent) GetResp() *backend.GeneResp {
	return &backend.GeneResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

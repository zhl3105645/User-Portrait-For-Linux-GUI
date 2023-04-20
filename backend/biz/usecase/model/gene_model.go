package model

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"backend/consumer/config"
	"context"
)

type GeneModel struct {
	modelId int64

	//
	appId int64
}

func NewGeneModel(modelId int64) *GeneModel {
	return &GeneModel{
		modelId: modelId,
	}
}

func (g *GeneModel) Load(ctx context.Context) error {
	if g.modelId <= 0 {
		return microtype.ParamCheckFailed
	}

	model, err := query.DataModel.WithContext(ctx).Where(query.DataModel.ModelID.Eq(g.modelId)).First()
	if err != nil {
		return microtype.DataModelQueryFailed
	}
	g.appId = model.AppID

	conf, err := config.ReadConfig()
	if err != nil {
		return err
	}

	// app_id typ=RuleGene status=begin
	if len(conf.Configs) == 0 {
		conf.Configs = make(map[int64]*config.Config)
		conf.Configs[g.appId] = &config.Config{}
	}

	if conf.Configs[g.appId] == nil {
		conf.Configs[g.appId] = &config.Config{}
	}

	if len(conf.Configs[g.appId].Config) == 0 {
		conf.Configs[g.appId].Config = make(map[config.TaskType]config.Status)
	}

	if len(conf.Configs[g.appId].Param) == 0 {
		conf.Configs[g.appId].Param = make(map[config.TaskType]int64)
	}

	status := conf.Configs[g.appId].Config[config.ModelGene]
	if status != config.Stop {
		return microtype.ModelDataGeneFailed
	}

	conf.Configs[g.appId].Param[config.ModelGene] = model.ModelID
	conf.Configs[g.appId].Config[config.ModelGene] = config.Begin

	err = config.WriteConfig(conf)
	if err != nil {
		return err
	}

	return nil
}

func (g *GeneModel) GetResp() *backend.GeneResp {
	return &backend.GeneResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

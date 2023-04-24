package label

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"context"
)

type GeneLabel struct {
	labelId int64

	//
	appId int64
}

func NewGeneLabel(labelId int64) *GeneLabel {
	return &GeneLabel{
		labelId: labelId,
	}
}

func (g *GeneLabel) Load(ctx context.Context) error {
	//if g.labelId <= 0 {
	//	return microtype.ParamCheckFailed
	//}
	//
	//label, err := query.Label.WithContext(ctx).Where(query.Label.LabelID.Eq(g.labelId)).First()
	//if err != nil {
	//	return microtype.LabelQueryFailed
	//}
	//model, err := query.DataModel.WithContext(ctx).Where(query.DataModel.ModelID.Eq(label.ModelID)).First()
	//if err != nil {
	//	return microtype.DataModelQueryFailed
	//}
	//g.appId = model.AppID
	//
	//conf, err := config.ReadConfig()
	//if err != nil {
	//	return err
	//}
	//
	//// app_id typ=LabelGene status=begin
	//if len(conf.Configs) == 0 {
	//	conf.Configs = make(map[int64]*config.Config)
	//	conf.Configs[g.appId] = &config.Config{}
	//}
	//
	//if conf.Configs[g.appId] == nil {
	//	conf.Configs[g.appId] = &config.Config{}
	//}
	//
	//if len(conf.Configs[g.appId].Config) == 0 {
	//	conf.Configs[g.appId].Config = make(map[config.TaskType]config.Status)
	//}
	//
	//if len(conf.Configs[g.appId].Param) == 0 {
	//	conf.Configs[g.appId].Param = make(map[config.TaskType]int64)
	//}
	//
	//status := conf.Configs[g.appId].Config[config.LabelGene]
	//if status != config.Stop {
	//	return microtype.ModelDataGeneFailed
	//}
	//
	//conf.Configs[g.appId].Param[config.LabelGene] = label.LabelID
	//conf.Configs[g.appId].Config[config.LabelGene] = config.Begin
	//
	//err = config.WriteConfig(conf)
	//if err != nil {
	//	return err
	//}
	//
	return nil
}

func (g *GeneLabel) GetResp() *backend.GeneResp {
	return &backend.GeneResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

package data_source

import (
	"backend/biz/entity/account"
	"backend/biz/entity/data_source"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type DataSources struct {
	accountId int64

	//
	appId int64
	res   []*backend.DataSource
}

func NewDataSources(accountId int64) *DataSources {
	return &DataSources{
		accountId: accountId,
	}
}

func (d *DataSources) Load(ctx context.Context) error {
	ac := account.NewAccount(d.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	d.appId = ac.GetQueryAccount().AppID

	// 数据源
	sourceDO := query.DataSource
	sourceMO := sourceDO.WithContext(ctx)
	sources, err := sourceMO.Where(sourceDO.AppID.Eq(d.appId)).Find()
	if err != nil {
		return microtype.DataSourceQueryFailed
	}

	// 规则
	ruleDO := query.Rule
	ruleMO := ruleDO.WithContext(ctx)
	rules, err := ruleMO.Where(ruleDO.AppID.Eq(d.appId)).Find()
	if err != nil {
		return microtype.RuleQueryFailed
	}

	ruleMap := make(map[int64]string, len(rules))
	for _, r := range rules {
		if r == nil {
			continue
		}
		ruleMap[r.RuleID] = r.RuleDesc
	}

	// 模型
	modelDO := query.DataModel
	modelMO := modelDO.WithContext(ctx)
	models, err := modelMO.Where(modelDO.AppID.Eq(d.appId)).Find()
	if err != nil {
		return microtype.DataModelQueryFailed
	}

	modelMap := make(map[int64]string, len(models))
	for _, m := range models {
		if m == nil {
			continue
		}
		modelMap[m.ModelID] = m.ModelName
	}

	// 汇总
	cate1 := &backend.DataSource{
		Value:    int64(data_source.Basic),
		Label:    data_source.SourceTypeDesc[data_source.Basic],
		Children: make([]*backend.DataSource, 0),
	}
	cate2 := &backend.DataSource{
		Value:    int64(data_source.EventRule),
		Label:    data_source.SourceTypeDesc[data_source.EventRule],
		Children: make([]*backend.DataSource, 0),
	}
	cate3 := &backend.DataSource{
		Value:    int64(data_source.BehaviorRule),
		Label:    data_source.SourceTypeDesc[data_source.BehaviorRule],
		Children: make([]*backend.DataSource, 0),
	}
	cate4 := &backend.DataSource{
		Value: int64(data_source.AllEventRule),
		Label: data_source.SourceTypeDesc[data_source.AllEventRule],
	}
	cate5 := &backend.DataSource{
		Value: int64(data_source.AllBehaviorRule),
		Label: data_source.SourceTypeDesc[data_source.AllBehaviorRule],
	}
	cate6 := &backend.DataSource{
		Value: int64(data_source.Model),
		Label: data_source.SourceTypeDesc[data_source.Model],
	}

	for _, source := range sources {
		if source == nil || source.SourceValue == nil {
			continue
		}
		typ := data_source.Type(source.SourceType)

		if typ == data_source.Basic {
			cate1.Children = append(cate1.Children, &backend.DataSource{
				Value: *source.SourceValue,
				Label: data_source.BasicSourceDesc[int(*source.SourceValue)],
			})
		} else if typ == data_source.EventRule {
			cate2.Children = append(cate2.Children, &backend.DataSource{
				Value: *source.SourceValue,
				Label: ruleMap[*source.SourceValue],
			})
		} else if typ == data_source.BehaviorRule {
			cate3.Children = append(cate3.Children, &backend.DataSource{
				Value: *source.SourceValue,
				Label: ruleMap[*source.SourceValue],
			})
		} else if typ == data_source.Model {
			cate6.Children = append(cate6.Children, &backend.DataSource{
				Value: *source.SourceValue,
				Label: modelMap[*source.SourceValue],
			})
		}
	}

	d.res = make([]*backend.DataSource, 0)
	if len(cate1.Children) > 0 {
		d.res = append(d.res, cate1)
	}
	if len(cate2.Children) > 0 {
		d.res = append(d.res, cate2)
	}
	if len(cate3.Children) > 0 {
		d.res = append(d.res, cate3)
	}
	if len(cate2.Children) > 0 {
		d.res = append(d.res, cate4)
	}
	if len(cate3.Children) > 0 {
		d.res = append(d.res, cate5)
	}
	if len(cate6.Children) > 0 {
		d.res = append(d.res, cate6)
	}

	return nil
}

func (d *DataSources) GetResp() *backend.DataSourceResp {
	return &backend.DataSourceResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		DataSource: d.res,
	}
}

package model

import (
	"backend/biz/entity/account"
	"backend/biz/entity/chart"
	"backend/biz/entity/data_model"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
	"strconv"
)

type PageModel struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string
	ruleType  int64

	//
	appId  int64
	models []*backend.Model
	total  int64
}

func NewPageModel(accountId int64, pageNum int64, pageSize int64, search string) *PageModel {
	return &PageModel{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

func (p *PageModel) Load(ctx context.Context) error {
	// appId
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	// 模型
	modelDO := query.DataModel
	modelMO := modelDO.WithContext(ctx)

	offset := (p.pageNum - 1) * p.pageSize
	models, count, err := modelMO.Where(modelDO.AppID.Eq(p.appId), modelDO.ModelName.Like("%"+p.search+"%")).
		FindByPage(int(offset), int(p.pageSize))
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return microtype.DataModelQueryFailed
	}

	if len(models) == 0 {
		return nil
	}
	p.total = count

	// 模型数据
	modelIds := make([]int64, 0, len(models))
	for _, m := range models {
		if m == nil {
			continue
		}
		modelIds = append(modelIds, m.ModelID)
	}
	dataDO := query.ModelDatum
	dataMO := dataDO.WithContext(ctx)
	modelData, err := dataMO.Where(dataDO.ModelID.In(modelIds...)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return microtype.ModelDataQueryFailed
	}
	// model_id -> data
	modelDataMap := make(map[int64][]*model.ModelDatum)
	for _, data := range modelData {
		if data == nil {
			continue
		}
		if v, ok := modelDataMap[data.ModelID]; !ok || len(v) == 0 {
			modelDataMap[data.ModelID] = make([]*model.ModelDatum, 0)
		}
		modelDataMap[data.ModelID] = append(modelDataMap[data.ModelID], data)
	}

	// 规则数据
	rules, err := query.Rule.WithContext(ctx).Where(query.Rule.AppID.Eq(p.appId)).Find()
	if err != nil && errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query rule failed. err=", err.Error())
	}
	ruleDescMap := make(map[int64]string)
	for _, r := range rules {
		if r == nil {
			continue
		}
		ruleDescMap[r.RuleID] = r.RuleDesc
	}

	// 整合数据
	p.models = make([]*backend.Model, 0, len(models))
	for _, m := range models {
		if m == nil {
			continue
		}
		p.models = append(p.models, &backend.Model{
			ModelName: m.ModelName,
			ModelID:   m.ModelID,
			ModelType: m.ModelType,
			Option:    getOption(m.DataType, modelDataMap[m.ModelID], ruleDescMap),
		})
	}

	return nil
}

func (p *PageModel) GetResp() *backend.ModelInPageResp {
	return &backend.ModelInPageResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Models:     p.models,
		Total:      p.total,
	}
}

func getOption(dataType int64, data []*model.ModelDatum, ruleDescMap map[int64]string) *backend.ChartOption {
	if len(data) == 0 {
		return nil
	}
	option := &backend.ChartOption{
		XAxis:   nil,
		YAxis:   nil,
		Tooltip: nil,
		Series:  nil,
	}

	if dataType == int64(data_model.Float) {
		return chart.GetModelOption(chart.CntFloat, data, nil)
	} else if dataType == int64(data_model.TimeDuration) {
		return chart.GetModelOption(chart.TimeFloat, data, nil)
	} else if dataType == int64(data_model.TimePeriod) {
		enumMap := make(map[int64]int64)
		for _, d := range data {
			if d == nil {
				continue
			}
			numF, _ := strconv.ParseFloat(d.Data, 64)
			num := int64(numF)
			if v, ok := enumMap[num]; ok {
				enumMap[num] = v + 1
			} else {
				enumMap[num] = 1
			}
		}
		yData := make([]string, 0, len(enumMap))
		for num, cnt := range enumMap {
			s := struct {
				Value int64  `json:"value"`
				Name  string `json:"name"`
			}{
				Value: cnt,
				Name:  fmt.Sprintf("%d~%d点", num, (num+2)%24),
			}

			str, _ := json.Marshal(s)
			yData = append(yData, string(str))
		}

		toolTip := &backend.ToolTip{
			Trigger: "item",
		}
		series := &backend.Series{
			Type: "pie",
			Data: yData,
		}
		option.Tooltip = toolTip
		option.Series = []*backend.Series{series}
	} else if dataType == int64(data_model.MultiTimeDuration) {
		return chart.GetModelOption(chart.All, data, ruleDescMap)
	} else if dataType == int64(data_model.Enum) {
		return chart.GetModelOption(chart.Enum, data, nil)
	}

	return option
}

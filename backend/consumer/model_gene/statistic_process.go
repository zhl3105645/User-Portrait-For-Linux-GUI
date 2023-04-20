package model_gene

import (
	"backend/biz/entity/data_model"
	"backend/biz/entity/data_source"
	"backend/biz/entity/rule"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"time"
)

func statisticsProcess(ctx context.Context, mo *model.DataModel) {
	// 数据源
	if mo.SourceID == nil {
		logger.Error("source id is empty")
		return
	}
	sourceDO := query.DataSource
	sourceMO := sourceDO.WithContext(ctx)
	source, err := sourceMO.Where(sourceDO.SourceID.Eq(*mo.SourceID)).First()
	if err != nil {
		logger.Error("source query failed. err=", err.Error())
		return
	}

	// 记录
	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User
	records, err := recordMO.LeftJoin(userDO, userDO.UserID.EqCol(recordDO.UserID)).
		Where(userDO.AppID.Eq(mo.AppID)).Find()
	if err != nil {
		logger.Error("record query failed. err=", err.Error())
		return
	}

	// u_id -> records
	recordMap := make(map[int64][]*model.Record)
	for _, record := range records {
		if record == nil {
			continue
		}

		if v, ok := recordMap[record.UserID]; !ok || len(v) == 0 {
			recordMap[record.UserID] = make([]*model.Record, 0)
		}
		recordMap[record.UserID] = append(recordMap[record.UserID], record)
	}

	// 新数据
	newMOs := make([]*model.ModelDatum, 0)

	dataRes := make(map[int64]string)

	switch data_source.Type(source.SourceType) {
	case data_source.Basic:
		if source.SourceValue == nil || mo.CalculateType == nil {
			logger.Error("basic value or calculate type is empty")
			return
		}

		dataRes = basicProcess(int(*source.SourceValue), data_model.CalculateType(*mo.CalculateType), recordMap)
	case data_source.EventRule:
		// TODO
	case data_source.BehaviorRule:
		// TODO
	case data_source.AllEventRule:
		//if mo.CalculateType == nil || *mo.CalculateType != int64(data_model.RuleCnt) {
		//	logger.Error("type is wrong")
		//	return
		//}
		// TODO

	case data_source.AllBehaviorRule:
		if mo.CalculateType == nil || *mo.CalculateType != int64(data_model.RuleDuration) {
			logger.Error("param is wrong")
			return
		}

		dataRes = allBehaviorProcess(ctx, mo.AppID, recordMap)
	}

	for uId, value := range dataRes {
		newMOs = append(newMOs, &model.ModelDatum{
			ModelDataID: 0,
			Data:        fmt.Sprintf("%s", value),
			ModelID:     mo.ModelID,
			UserID:      uId,
		})
	}

	// 已有数据
	dataDO := query.ModelDatum
	dataMO := dataDO.WithContext(ctx)

	dbDatas, err := dataMO.Where(dataDO.ModelID.Eq(mo.ModelID)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query model data failed. err=", err.Error())
		return
	}

	dbDataMap := make(map[int64]*model.ModelDatum)
	for _, data := range dbDatas {
		if data == nil {
			continue
		}
		dbDataMap[data.UserID] = data
	}

	// 划分: db存在 db不存在
	createMos := make([]*model.ModelDatum, 0)
	updateMos := make([]model.ModelDatum, 0)

	for _, r := range newMOs {
		if r == nil {
			continue
		}
		if v, ok := dbDataMap[r.UserID]; ok && v != nil {
			updateMos = append(updateMos, model.ModelDatum{
				ModelDataID: v.ModelDataID,
				Data:        r.Data,
				ModelID:     v.ModelID,
				UserID:      v.UserID,
			})
		} else {
			createMos = append(createMos, r)
		}
	}

	// 写入
	err = dataMO.Create(createMos...)
	if err != nil {
		logger.Error("create model data failed. err=", err.Error())
		return
	}

	for _, r := range updateMos {
		_, err = dataMO.Where(dataDO.ModelDataID.Eq(r.ModelDataID)).Updates(r)
		if err != nil {
			logger.Error("update model data failed. err=", err.Error())
			return
		}
	}
}

func allBehaviorProcess(ctx context.Context, appId int64, recordMap map[int64][]*model.Record) map[int64]string {
	res := make(map[int64]string)
	// 规则
	eventRules, behaviorRules, err := rule.GetRuleModels(ctx, appId)
	if err != nil {
		logger.Error("err=", err.Error())
		return res
	}
	ruleMap := make(map[int64]string)
	for _, r := range eventRules {
		if r == nil {
			continue
		}

		ruleMap[r.RuleID] = r.RuleDesc
	}
	for _, r := range behaviorRules {
		if r == nil {
			continue
		}

		ruleMap[r.RuleID] = r.RuleDesc
	}

	// 数据转换
	for userId, records := range recordMap {
		// 用户维度
		behaviorId2Duraion := make(map[int64]int64) // rule_id -> duration
		for _, r := range records {
			if r == nil {
				continue
			}

			if r.BehaviorRuleValue == nil {
				continue
			}
			eles := rule.ParseRuleElements(*r.BehaviorRuleValue, ruleMap)
			if len(eles) == 0 {
				continue
			}

			// 单次记录时间
			durationMap := rule.GetBehaviorDuration(eles)
			for id, duration := range durationMap {
				if v, ok := behaviorId2Duraion[id]; ok {
					behaviorId2Duraion[id] = v + duration
				} else {
					behaviorId2Duraion[id] = duration
				}
			}
		}
		ave := make(map[int64]int64)
		for id, duration := range behaviorId2Duraion {
			ave[id] = duration / int64(len(records))
		}
		if len(ave) == 0 {
			continue
		}
		// json化
		jsonStr, err := json.Marshal(ave)
		if err != nil {
			logger.Error("json marshal failed. err=", err.Error())
		}
		res[userId] = string(jsonStr)
	}

	return res
}

func basicProcess(basicTyp int, calculateTyp data_model.CalculateType, recordMap map[int64][]*model.Record) map[int64]string {
	// uid -> 数据列表
	dataMap := make(map[int64][]float64, len(recordMap))
	for uId, records := range recordMap {
		dataMap[uId] = make([]float64, 0, len(records))
		for _, record := range records {
			if record == nil {
				continue
			}

			switch basicTyp {
			case data_source.SourceMouseClickCnt:
				if record.MouseClickCnt == nil {
					continue
				}
				dataMap[uId] = append(dataMap[uId], float64(*record.MouseClickCnt))
			case data_source.SourceMouseMoveCnt:
				if record.MouseMoveCnt == nil {
					continue
				}
				dataMap[uId] = append(dataMap[uId], float64(*record.MouseMoveCnt))
			case data_source.SourceMoveDis:
				if record.MouseMoveDis == nil {
					continue
				}
				dataMap[uId] = append(dataMap[uId], *record.MouseMoveDis)
			case data_source.SourceMouseWheelCnt:
				if record.MouseWheelCnt == nil {
					continue
				}
				dataMap[uId] = append(dataMap[uId], float64(*record.MouseWheelCnt))
			case data_source.SourceKeyClickCnt:
				if record.KeyClickCnt == nil {
					continue
				}
				dataMap[uId] = append(dataMap[uId], float64(*record.KeyClickCnt))
			case data_source.SourceKeyClickSpeed:
				if record.KeyClickSpeed == nil {
					continue
				}
				dataMap[uId] = append(dataMap[uId], *record.KeyClickSpeed)
			case data_source.SourceShortCut:
				if record.ShortcutCnt == nil {
					continue
				}
				dataMap[uId] = append(dataMap[uId], float64(*record.ShortcutCnt))
			case data_source.SourceUsePeriod:
				hour := time.UnixMilli(record.BeginTime).Hour()
				dataMap[uId] = append(dataMap[uId], float64(hour))
			case data_source.SourceUseTime:
				if record.UseTime == 0 {
					continue
				}
				dataMap[uId] = append(dataMap[uId], float64(record.UseTime))
			default:
			}
		}
	}

	// 计算方式
	res := make(map[int64]string, len(recordMap))
	if calculateTyp == data_model.Average {
		for uId, data := range dataMap {
			if len(data) == 0 {
				continue
			}

			sum := funk.SumFloat64(data)
			res[uId] = fmt.Sprintf("%.1f", sum/float64(len(data)))
		}
	} else if calculateTyp == data_model.Mode {
		for uId, data := range dataMap {
			if len(data) == 0 {
				continue
			}

			cntMap := make(map[float64]int)
			for _, d := range data {
				if v, ok := cntMap[d]; ok {
					cntMap[d] = v + 1
				} else {
					cntMap[d] = 1
				}
			}
			maxValue := data[0]
			maxCnt := cntMap[maxValue]
			for value, cnt := range cntMap {
				if cnt > maxCnt {
					maxCnt = cnt
					maxValue = value
				}
			}
			res[uId] = fmt.Sprintf("%.1f", maxValue)
		}
	}

	return res
}

package label_gene

import (
	"backend/biz/usecase/label"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"backend/consumer/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
	"strconv"
)

func Gene(appId int64, labelId int64) {
	defer geneDone(appId)

	ctx := context.Background()

	// 标签信息
	labelDO := query.Label
	labelMO := labelDO.WithContext(ctx)
	lab, err := labelMO.Where(labelDO.LabelID.Eq(labelId)).First()
	if err != nil {
		logger.Error("query label failed. err=", err.Error())
		return
	}
	convertRules := make([]*label.ConvertRule, 0)
	err = json.Unmarshal([]byte(lab.LabelConvertRule), &convertRules)
	if err != nil {
		logger.Error("json unmarshal failed. err=", err.Error())
		return
	}

	// 模型数据
	modelDO := query.ModelDatum
	modelMO := modelDO.WithContext(ctx)
	modelData, err := modelMO.Where(modelDO.ModelID.Eq(lab.ModelID)).Find()
	if err != nil {
		logger.Error("query model data failed. err=", err.Error())
		return
	}

	// 数据转换
	labelDataMap := convert(modelData, convertRules)
	if len(labelDataMap) == 0 {
		logger.Info("new data is empty. err=", err.Error())
		return
	}

	// 写入标签数据
	newMOs := make([]*model.LabelDatum, 0, len(labelDataMap))
	for userId, data := range labelDataMap {
		newMOs = append(newMOs, &model.LabelDatum{
			LabelDataID: 0,
			Data:        fmt.Sprintf("%d", data),
			LabelID:     lab.LabelID,
			UserID:      userId,
		})
	}

	// 已有数据
	dataDO := query.LabelDatum
	dataMO := dataDO.WithContext(ctx)
	dbMOs, err := dataMO.
		Where(dataDO.LabelID.Eq(lab.LabelID)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query db label data failed. err=", err.Error())
		return
	}

	dbDataMap := make(map[int64]*model.LabelDatum)
	for _, data := range dbMOs {
		if data == nil {
			continue
		}
		dbDataMap[data.UserID] = data
	}

	// 分为 更新 和 新增
	createMos := make([]*model.LabelDatum, 0)
	updateMos := make([]model.LabelDatum, 0)

	for _, r := range newMOs {
		if r == nil {
			continue
		}
		if v, ok := dbDataMap[r.UserID]; ok && v != nil {
			updateMos = append(updateMos, model.LabelDatum{
				LabelDataID: v.LabelDataID,
				Data:        r.Data,
				LabelID:     v.LabelID,
				UserID:      v.UserID,
			})
		} else {
			createMos = append(createMos, r)
		}
	}

	// 写入
	err = dataMO.Create(createMos...)
	if err != nil {
		logger.Error("create label data failed. err=", err.Error())
		return
	}

	for _, r := range updateMos {
		_, err = dataMO.
			Where(dataDO.LabelDataID.Eq(r.LabelDataID)).Updates(r)
		if err != nil {
			logger.Error("update label data failed. err=", err.Error())
			return
		}
	}

	return
}

func convert(data []*model.ModelDatum, rules []*label.ConvertRule) map[int64]int64 {
	res := make(map[int64]int64)
	for _, d := range data {
		if d == nil {
			continue
		}
		modelValue, err := strconv.ParseFloat(d.Data, 64)
		if err != nil {
			logger.Error("parse float failed. err=", err.Error())
			continue
		}
		for _, rule := range rules {
			if rule.Match(modelValue) {
				res[d.UserID] = rule.YValue
				break
			}
		}
	}

	return res
}

func geneDone(appId int64) {
	// running -> stop
	config.StatusChan <- &config.StatusChange{
		AppId:    appId,
		TaskType: config.LabelGene,
		Status:   config.Stop,
	}
}

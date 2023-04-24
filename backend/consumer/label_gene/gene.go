package label_gene

import (
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"backend/consumer/config"
	"context"
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
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

	labelDataMap := process(ctx, appId, labelId)

	// 写入标签数据
	newMOs := make([]*model.LabelDatum, 0, len(labelDataMap))
	for userId, data := range labelDataMap {
		newMOs = append(newMOs, &model.LabelDatum{
			LabelDataID: 0,
			Data:        data,
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

func geneDone(appId int64) {
	// running -> stop
	config.StatusChan <- &config.StatusChange{
		AppId:    appId,
		TaskType: config.LabelGene,
		Status:   config.Stop,
	}
}

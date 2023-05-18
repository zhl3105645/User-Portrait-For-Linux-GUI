package model_gene

//
//import (
//	"backend/biz/entity/data_model"
//	"backend/cmd/dal/query"
//	"backend/consumer/config"
//	"context"
//	"github.com/bytedance/gopkg/util/logger"
//)
//
//func Gene(appId, modelId int64) {
//	defer geneDone(appId)
//	ctx := context.Background()
//
//	modelDO := query.DataModel
//	modelMO := modelDO.WithContext(ctx)
//	model, err := modelMO.Where(modelDO.ModelID.Eq(modelId)).First()
//	if err != nil {
//		logger.Error("query model failed. err=", err.Error())
//		return
//	}
//
//	if model.ModelType == int64(data_model.Statistics) {
//		statisticsProcess(ctx, model)
//	} else if model.ModelType == int64(data_model.Learning) {
//		learningProcess(ctx, model)
//	}
//
//	return
//}
//
//func geneDone(appId int64) {
//	// running -> stop
//	config.StatusChan <- &config.StatusChange{
//		AppId:    appId,
//		TaskType: config.ModelGene,
//		Status:   config.Stop,
//	}
//}

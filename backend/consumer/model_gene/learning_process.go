package model_gene

import (
	"backend/biz/entity/data_model"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"gorm.io/gorm"
)

func learningProcess(ctx context.Context, mo *model.DataModel) {
	// 参数校验
	if mo.MlParam == nil {
		logger.Error("ml param is empty")
		return
	}
	mlParam := &backend.LearningParam{}
	err := json.Unmarshal([]byte(*mo.MlParam), mlParam)
	if err != nil {
		logger.Error("unmarshal failed. err=", err.Error())
		return
	}

	if mlParam.HTTPType != int64(data_model.Post) || mlParam.HTTPAddr == "" {
		logger.Error("参数错误")
		return
	}

	// 用户数据
	users, err := query.User.WithContext(ctx).Where(query.User.AppID.Eq(mo.AppID)).Find()
	if err != nil {
		logger.Error("query user failed. err=", err.Error())
		return
	}

	// 模型数据
	modelDataMap := make(map[int64]map[int64]string) // source_id -> user_id -> data
	for _, param := range mlParam.BodyParams {
		if param == nil {
			continue
		}

		modelDatas, err := query.ModelDatum.WithContext(ctx).
			Where(query.ModelDatum.ModelID.Eq(param.ModelID)).Find()
		if err != nil {
			logger.Error("query model data failed. err=", err.Error())
			return
		}

		modelDataMap[param.ModelID] = make(map[int64]string)
		for _, data := range modelDatas {
			if data == nil {
				continue
			}
			modelDataMap[param.ModelID][data.UserID] = data.Data
		}
	}

	// httpClient
	c, err := client.NewClient()
	if err != nil {
		return
	}

	resMap := make(map[int64]string, len(users)) // user_id -> data
	for _, user := range users {
		params := make(map[string]string)
		for _, param := range mlParam.BodyParams {
			params[param.Name] = modelDataMap[param.ModelID][user.UserID]
		}
		res, err := post(ctx, c, mlParam.HTTPAddr, params, mlParam.GetHTTPRespName())
		if err != nil {
			logger.Error("http request wrong. err=", err.Error())
			continue
		}

		if res == "" {
			continue
		}

		resMap[user.UserID] = res
	}

	// 写入模型数据
	newMOs := make([]*model.ModelDatum, 0, len(resMap))
	for userId, data := range resMap {
		newMOs = append(newMOs, &model.ModelDatum{
			//ModelDataID: 0,
			Data:    data,
			ModelID: mo.ModelID,
			UserID:  userId,
		})
	}

	// 已有数据
	dbMOs, err := query.ModelDatum.WithContext(ctx).
		Where(query.ModelDatum.ModelID.Eq(mo.ModelID)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query model data failed. err=", err.Error())
		return
	}

	dbDataMap := make(map[int64]*model.ModelDatum)
	for _, data := range dbMOs {
		if data == nil {
			continue
		}
		dbDataMap[data.UserID] = data
	}

	// 分为 更新 和 新增
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
	err = query.ModelDatum.WithContext(ctx).Create(createMos...)
	if err != nil {
		logger.Error("create model data failed. err=", err.Error())
		return
	}

	for _, r := range updateMos {
		_, err = query.ModelDatum.WithContext(ctx).
			Where(query.ModelDatum.ModelDataID.Eq(r.ModelDataID)).Updates(r)
		if err != nil {
			logger.Error("update model data failed. err=", err.Error())
			return
		}
	}
}

func post(ctx context.Context, client *client.Client, addr string, params map[string]string, respName string) (string, error) {
	var postArgs protocol.Args
	for k, v := range params {
		if v == "" {
			return "", fmt.Errorf("param is nil")
		}
		postArgs.Set(k, v)
	}
	status, body, _ := client.Post(ctx, nil, addr, &postArgs)
	if status != 200 {
		return "", fmt.Errorf("status != 200")
	}

	m := make(map[string]string)
	err := json.Unmarshal(body, &m)
	if err != nil {
		return "", err
	}

	return m[respName], nil
}

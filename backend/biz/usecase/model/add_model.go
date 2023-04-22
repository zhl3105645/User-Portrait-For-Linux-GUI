package model

import (
	"backend/biz/entity/account"
	"backend/biz/entity/data_model"
	"backend/biz/entity/data_source"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"strconv"
)

type AddModel struct {
	accountId int64
	req       backend.AddModelReq

	//
	appId int64
}

func NewAddModel(accountId int64, req backend.AddModelReq) *AddModel {
	return &AddModel{
		accountId: accountId,
		req:       req,
	}
}

func (a *AddModel) Load(ctx context.Context) error {
	if !a.check() {
		return microtype.ParamCheckFailed
	}

	ac := account.NewAccount(a.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	a.appId = ac.GetQueryAccount().AppID

	typ, _ := strconv.ParseInt(a.req.ModelType, 10, 64)
	httpDataType := int64(0)
	if a.req.LearningParam != nil {
		httpDataType = a.req.LearningParam.HTTPRespDataType
	}

	mo := &model.DataModel{
		ModelID:       0,
		ModelType:     typ,
		AppID:         a.appId,
		DataType:      int64(data_model.GetDataType(typ, a.req.SourceType, a.req.SourceValue, a.req.CalculateType, httpDataType)),
		ModelName:     a.req.ModelName,
		SourceID:      nil,
		CalculateType: proto.Int64(a.req.CalculateType),
		MlParam:       nil,
	}

	if typ == int64(data_model.Statistics) {
		sourceId := int64(0)
		sourceDO := query.DataSource
		sourceMO := sourceDO.WithContext(ctx)
		if a.req.SourceType == int64(data_source.AllEventRule) || a.req.SourceType == int64(data_source.AllBehaviorRule) {
			queryMo, err := sourceMO.
				Where(sourceDO.AppID.Eq(a.appId),
					sourceDO.SourceType.Eq(a.req.SourceType)).
				First()
			if err != nil {
				return microtype.DataSourceQueryFailed
			}
			sourceId = queryMo.SourceID
		} else {
			queryMo, err := sourceMO.
				Where(sourceDO.AppID.Eq(a.appId),
					sourceDO.SourceType.Eq(a.req.SourceType),
					sourceDO.SourceValue.Eq(a.req.SourceValue)).
				First()
			if err != nil {
				return microtype.DataSourceQueryFailed
			}
			sourceId = queryMo.SourceID
		}
		mo.SourceID = proto.Int64(sourceId)
	} else if typ == int64(data_model.Learning) {
		s, err := json.Marshal(a.req.LearningParam)
		if err != nil {
			return microtype.JsonMarshalFailed
		}

		mo.MlParam = proto.String(string(s))
	}

	// 写入模型
	modelDO := query.DataModel
	modelMO := modelDO.WithContext(ctx)

	err := modelMO.Create(mo)
	if err != nil {
		return microtype.DataModelCreateFailed
	}

	// 添加数据源
	err = data_source.AddModelSource(ctx, mo.ModelID, mo.AppID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddModel) GetResp() *backend.AddModelResp {
	return &backend.AddModelResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

func (a *AddModel) check() bool {
	if a.req.ModelName == "" {
		return false
	}
	if a.req.ModelType != strconv.FormatInt(int64(data_model.Statistics), 10) && a.req.ModelType != strconv.FormatInt(int64(data_model.Learning), 10) {
		return false
	}

	if a.req.ModelType == strconv.FormatInt(int64(data_model.Learning), 10) {
		if a.req.LearningParam == nil {
			return false
		}
	}

	return true
}

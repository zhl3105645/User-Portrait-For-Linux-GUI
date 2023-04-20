package model

import (
	"backend/biz/entity/data_source"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type DeleteModel struct {
	modelId int64
}

func NewDeleteModel(modelId int64) *DeleteModel {
	return &DeleteModel{
		modelId: modelId,
	}
}

func (d *DeleteModel) Load(ctx context.Context) error {
	_, err := query.ModelDatum.WithContext(ctx).
		Where(query.ModelDatum.ModelID.Eq(d.modelId)).
		Delete()
	if err != nil {
		return microtype.ModelDataDeleteFailed
	}

	_, err = query.DataModel.WithContext(ctx).
		Where(query.DataModel.ModelID.Eq(d.modelId)).
		Delete()
	if err != nil {
		return microtype.DataModelDeleteFailed
	}

	err = data_source.DeleteModelSource(ctx, d.modelId)
	if err != nil {
		return err
	}

	return nil
}

func (d *DeleteModel) GetResp() *backend.DeleteModelResp {
	return &backend.DeleteModelResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

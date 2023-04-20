package label

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type DeleteLabel struct {
	labelId int64
}

func NewDeleteLabel(labelId int64) *DeleteLabel {
	return &DeleteLabel{
		labelId: labelId,
	}
}

func (d *DeleteLabel) Load(ctx context.Context) error {
	_, err := query.LabelDatum.WithContext(ctx).
		Where(query.LabelDatum.LabelID.Eq(d.labelId)).
		Delete()
	if err != nil {
		return microtype.LabelDataDeleteFailed
	}

	_, err = query.Label.WithContext(ctx).
		Where(query.Label.LabelID.Eq(d.labelId)).
		Delete()
	if err != nil {
		return microtype.LabelDeleteFailed
	}

	return nil
}

func (d *DeleteLabel) GetResp() *backend.DeleteLabelResp {
	return &backend.DeleteLabelResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

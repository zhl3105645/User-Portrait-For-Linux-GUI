package element

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type DeleteElement struct {
	elementId int64
}

func NewDeleteElement(elementId int64) *DeleteElement {
	return &DeleteElement{
		elementId: elementId,
	}
}

func (d *DeleteElement) Load(ctx context.Context) error {
	if d.elementId <= 0 {
		return microtype.ParamCheckFailed
	}

	_, err := query.RuleElement.WithContext(ctx).
		Where(query.RuleElement.RuleElementID.Eq(d.elementId)).
		Delete()
	if err != nil {
		return microtype.ElementDeleteFailed
	}

	return nil
}

func (d *DeleteElement) GetResp() *backend.DeleteElementResp {
	return &backend.DeleteElementResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

package crowd

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type DeleteCrowd struct {
	crowdId int64
}

func NewDeleteCrowd(crowdId int64) *DeleteCrowd {
	return &DeleteCrowd{
		crowdId: crowdId,
	}
}

func (d *DeleteCrowd) Load(ctx context.Context) error {
	// 删除人群数据
	_, err := query.CrowdRelation.WithContext(ctx).Where(query.CrowdRelation.CrowdID.Eq(d.crowdId)).Delete()
	if err != nil {
		return err
	}
	// 删除人群
	_, err = query.Crowd.WithContext(ctx).Where(query.Crowd.CrowdID.Eq(d.crowdId)).Delete()
	if err != nil {
		return err
	}

	return nil
}
func (d *DeleteCrowd) GetResp() *backend.DeleteCrowdResp {
	return &backend.DeleteCrowdResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

package crowd

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/mq"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
)

type GeneCrowd struct {
	crowdId int64

	//
	appId int64
}

func NewGeneCrowd(crowdId int64) *GeneCrowd {
	return &GeneCrowd{
		crowdId: crowdId,
	}
}

func (g *GeneCrowd) Load(ctx context.Context) error {
	crowd, err := query.Crowd.WithContext(ctx).Where(query.Crowd.CrowdID.Eq(g.crowdId)).First()
	if err != nil {
		return err
	}

	g.appId = crowd.AppID

	msg := &mq.GeneMsg{
		AppId: g.appId,
		Type:  mq.CrowdGene,
		Param: g.crowdId,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return microtype.JsonMarshalFailed
	}

	err = mq.SendSyncMessage(string(msgJson))
	if err != nil {
		return microtype.MQSendFailed
	}

	return nil
}

func (g *GeneCrowd) GetResp() *backend.GeneResp {
	return &backend.GeneResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

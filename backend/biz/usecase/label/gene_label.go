package label

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/mq"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
)

type GeneLabel struct {
	labelId int64

	//
	appId int64
}

func NewGeneLabel(labelId int64) *GeneLabel {
	return &GeneLabel{
		labelId: labelId,
	}
}

func (g *GeneLabel) Load(ctx context.Context) error {
	label, err := query.Label.WithContext(ctx).Where(query.Label.LabelID.Eq(g.labelId)).First()
	if err != nil {
		return err
	}

	g.appId = label.AppID

	msg := &mq.GeneMsg{
		AppId: g.appId,
		Type:  mq.LabelGene,
		Param: g.labelId,
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

func (g *GeneLabel) GetResp() *backend.GeneResp {
	return &backend.GeneResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

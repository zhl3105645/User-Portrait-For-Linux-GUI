package seq_mining

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/mq"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"time"
)

type SeqMining struct {
	accountId int64
	percent   float64
	taskName  string

	//
	appId int64
}

func NewSeqMining(accountId int64, percent float64, taskName string) *SeqMining {
	return &SeqMining{
		accountId: accountId,
		percent:   percent,
		taskName:  taskName,
	}
}

func (s *SeqMining) Load(ctx context.Context) error {
	ac := account.NewAccount(s.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	s.appId = ac.GetQueryAccount().AppID

	mo := &model.SeqMiningTask{
		TaskName:     s.taskName,
		CreateTime:   time.Now(),
		Status:       StatusBegin,
		Percent:      int64(s.percent),
		Event2number: nil,
		Result:       nil,
		AppID:        s.appId,
	}
	err := query.SeqMiningTask.WithContext(ctx).Create(mo)
	if err != nil {
		return err
	}

	msg := &mq.GeneMsg{
		AppId: s.appId,
		Type:  mq.SeqMining,
		Param: mo.TaskID,
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

func (s *SeqMining) GetResp() *backend.SeqMiningResp {
	return &backend.SeqMiningResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

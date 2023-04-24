package label

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Labels struct {
	accountId int64

	//
	appId int64
	res   []*backend.Label
}

func NewLabels(accountId int64) *Labels {
	return &Labels{
		accountId: accountId,
	}
}

func (l *Labels) Load(ctx context.Context) error {
	ac := account.NewAccount(l.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	l.appId = ac.GetQueryAccount().AppID

	labels, err := query.Label.WithContext(ctx).Where(query.Label.AppID.Eq(l.appId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return microtype.LabelQueryFailed
	}

	l.res = make([]*backend.Label, 0, len(labels))
	for _, lab := range labels {
		if lab == nil {
			continue
		}

		l.res = append(l.res, &backend.Label{
			LabelName: lab.LabelName,
			LabelID:   lab.LabelID,
		})
	}

	return nil
}

func (l *Labels) GetResp() *backend.LabelsResp {
	return &backend.LabelsResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Labels:     l.res,
	}
}

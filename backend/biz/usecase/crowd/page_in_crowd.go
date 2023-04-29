package crowd

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

type PageInCrowd struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string

	//
	appId  int64
	total  int64
	crowds []*backend.Crowd
}

func NewPageInCrowd(accountId int64, pageNum int64, pageSize int64, search string) *PageInCrowd {
	return &PageInCrowd{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

func (p *PageInCrowd) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	off := (p.pageNum - 1) * p.pageSize
	res, count, err := query.Crowd.WithContext(ctx).
		Where(query.Crowd.AppID.Eq(p.appId), query.Crowd.CrowdName.Like("%"+p.search+"%")).
		FindByPage(int(off), int(p.pageSize))
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return err
	}
	p.total = count

	crowdIds := make([]int64, 0, len(res))
	for _, c := range res {
		crowdIds = append(crowdIds, c.CrowdID)
	}

	crowdId2UserNum := make(map[int64]int64)
	crowdRelations, err := query.CrowdRelation.WithContext(ctx).Where(query.CrowdRelation.CrowdID.In(crowdIds...)).Find()
	if err == nil {
		for _, relation := range crowdRelations {
			if cnt, ok := crowdId2UserNum[relation.CrowdID]; ok {
				crowdId2UserNum[relation.CrowdID] = cnt + 1
			} else {
				crowdId2UserNum[relation.CrowdID] = 1
			}
		}
	} else {
		return err
	}

	p.crowds = make([]*backend.Crowd, 0, len(res))
	for _, r := range res {
		if r == nil {
			continue
		}

		rules := make([]*backend.DivideRule, 0)

		err := json.Unmarshal([]byte(r.CrowdDivideRule), &rules)
		if err != nil {
			continue
		}

		c := &backend.Crowd{
			CrowdID:     r.CrowdID,
			CrowdName:   r.CrowdName,
			CrowdDesc:   r.CrowdDesc,
			DivideRules: rules,
			UserNum:     crowdId2UserNum[r.CrowdID],
		}
		p.crowds = append(p.crowds, c)
	}

	return nil
}

func (p *PageInCrowd) GetResp() *backend.CrowdInPageResp {
	return &backend.CrowdInPageResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Total:      p.total,
		Crowds:     p.crowds,
	}
}

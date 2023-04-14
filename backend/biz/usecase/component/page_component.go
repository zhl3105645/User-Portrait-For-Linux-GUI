package component

import (
	"backend/biz/entity/account"
	"backend/biz/entity/component"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"context"
)

type PageComponent struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string

	//
	appId int64
	cp    *component.Component
}

func NewPageComponent(accountId int64, pageNum int64, pageSize int64, search string) *PageComponent {
	return &PageComponent{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

func (p *PageComponent) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	p.cp = component.NewComponent(component.QueryPage, p.appId, &component.QueryParam{
		PageNum:  p.pageNum,
		PageSize: p.pageSize,
		Search:   p.search,
	}, nil)
	if err := p.cp.Load(ctx); err != nil {
		return err
	}

	return nil
}

func (p *PageComponent) GetResp() *backend.ComponentInPageResp {
	resp := &backend.ComponentInPageResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Components: nil,
		Total:      p.cp.GetTotal(),
	}

	coms := make([]*backend.Component, 0)
	for _, v := range p.cp.GetQueryComponent() {
		if v == nil {
			continue
		}

		com := &backend.Component{
			ComponentID:   v.ComponentID,
			ComponentName: v.ComponentName,
			ComponentType: v.ComponentType,
		}

		if v.ComponentDesc != nil {
			com.ComponentDesc = *v.ComponentDesc
		}

		coms = append(coms, com)
	}

	resp.Components = coms

	return resp
}

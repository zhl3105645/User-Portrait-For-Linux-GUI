package user

import (
	"backend/biz/entity/account"
	"backend/biz/entity/event_data"
	"backend/biz/entity/user"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"context"
	"sync"
)

type PageUser struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string

	//
	appId     int64
	us        *user.User
	recordNum map[int64]int64
}

func NewPageUser(accountId int64, pageNum int64, pageSize int64, search string) *PageUser {
	return &PageUser{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
		recordNum: make(map[int64]int64),
	}
}

func (p *PageUser) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	var (
		wg        sync.WaitGroup
		userErr   error
		recordErr error
	)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			_ = recover()
		}()

		p.us = user.NewUser(user.QueryPage, nil, &user.QueryParam{
			AppId:    p.appId,
			PageNum:  p.pageNum,
			PageSize: p.pageSize,
			Search:   p.search,
		}, 0)
		if err := p.us.Load(ctx); err != nil {
			userErr = err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			_ = recover()
		}()

		ed := event_data.NewEvent(p.appId)
		if err := ed.Load(ctx); err != nil {
			recordErr = err
			return
		}

		p.recordNum = ed.GetEventsRecordNum()
	}()

	wg.Wait()

	if userErr != nil {
		return userErr
	}

	if recordErr != nil {
		return recordErr
	}

	return nil
}

func (p *PageUser) GetResp() *backend.UserInPageResp {
	resp := &backend.UserInPageResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Total:      p.us.GetTotal(),
	}

	users := make([]*backend.User, 0)
	for _, v := range p.us.GetPageQueryUser() {
		if v == nil {
			continue
		}

		users = append(users, &backend.User{
			UserID:    v.UserID,
			UserName:  v.UserName,
			RecordNum: p.recordNum[v.UserID],
		})
	}

	resp.Users = users

	return resp
}

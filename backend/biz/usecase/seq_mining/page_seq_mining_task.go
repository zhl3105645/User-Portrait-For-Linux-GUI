package seq_mining

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type PageTask struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string

	//
	appId int64
	total int64
	res   []*backend.SeqMiningTask
}

func NewPageTask(accountId int64, pageNum int64, pageSize int64, search string) *PageTask {
	return &PageTask{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

func (p *PageTask) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	offset := (p.pageNum - 1) * p.pageSize
	res, count, err := query.SeqMiningTask.WithContext(ctx).
		Where(query.SeqMiningTask.AppID.Eq(p.appId), query.SeqMiningTask.TaskName.Like("%"+p.search+"%")).
		Order(query.SeqMiningTask.CreateTime.Desc()).
		FindByPage(int(offset), int(p.pageSize))
	if err != nil {
		return err
	}
	p.total = count

	p.res = make([]*backend.SeqMiningTask, 0, len(res))
	for _, r := range res {
		if r == nil {
			continue
		}

		p.res = append(p.res, &backend.SeqMiningTask{
			TaskID:     r.TaskID,
			TaskName:   r.TaskName,
			CreateTime: r.CreateTime.Format("2006-01-02 15:04:05"),
			TaskStatus: r.Status,
			Percent:    r.Percent,
		})
	}

	return nil
}

func (p *PageTask) GetResp() *backend.SeqMiningTaskInPageResp {
	return &backend.SeqMiningTaskInPageResp{
		StatusCode:     microtype.SuccessErr.Code,
		StatusMsg:      microtype.SuccessErr.Msg,
		Total:          p.total,
		SeqMiningTasks: p.res,
	}
}

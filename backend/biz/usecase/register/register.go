package register

import (
	"backend/biz/entity/account"
	"backend/biz/entity/app"
	"backend/biz/entity/data_source"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"context"
)

type Register struct {
	req backend.RegisterReq

	app *model.App
}

func NewRegister(req backend.RegisterReq) *Register {
	return &Register{
		req: req,
	}
}

func (r *Register) Load(ctx context.Context) error {
	// 新增app
	ap := app.NewApp(0, r.req.AppName, false)
	if err := ap.Load(ctx); err != nil {
		return err
	}

	if mo, err := ap.Add(ctx); err != nil {
		return err
	} else {
		r.app = mo
	}

	// 新增账号
	ac := account.NewAccount(
		0,
		r.app.AppID,
		r.req.AccountName,
		r.req.AccountPwd,
		account.Creator,
		account.Create,
	)

	if err := ac.Load(ctx); err != nil {
		return err
	}

	// 初始化数据源
	err := data_source.InitDataSource(ctx, r.app.AppID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Register) GetResp() *backend.RegisterResp {
	return &backend.RegisterResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

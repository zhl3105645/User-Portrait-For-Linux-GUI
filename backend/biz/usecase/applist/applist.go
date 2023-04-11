package applist

import (
	"backend/biz/entity/app"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"context"
)

type AppList struct {
	apps []*model.App
}

func NewAppList() *AppList {
	return &AppList{}
}

func (a *AppList) Load(ctx context.Context) error {
	ap := app.NewApp(0, "", true)
	if err := ap.Load(ctx); err != nil {
		return err
	}

	a.apps = ap.FindAll()

	return nil
}

func (a *AppList) GetResp() *backend.AppListResp {
	resp := &backend.AppListResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}

	list := make([]*backend.App, 0, len(a.apps))
	for _, ap := range a.apps {
		if ap == nil || ap.AppName == nil {
			continue
		}
		list = append(list, &backend.App{
			AppID:   ap.AppID,
			AppName: *ap.AppName,
		})
	}

	resp.Apps = list

	return resp
}

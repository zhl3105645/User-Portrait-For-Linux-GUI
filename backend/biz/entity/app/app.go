package app

import (
	"backend/biz/microtype"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"github.com/golang/protobuf/proto"
)

type App struct {
	appId    int64
	appName  string
	queryAll bool

	ap  *model.App
	aps []*model.App
}

func NewApp(appId int64, appName string, queryAll bool) *App {
	return &App{
		appId:    appId,
		appName:  appName,
		queryAll: queryAll,
	}
}

func (a *App) Load(ctx context.Context) error {
	// 查询全部
	if a.queryAll {
		if aps, err := query.App.WithContext(ctx).Find(); err != nil {
			return microtype.AppFindAllFailed
		} else {
			a.aps = aps
		}

		return nil
	}

	// 查询某一应用
	ap, err := query.App.WithContext(ctx).
		Where(query.App.AppID.Eq(a.appId)).
		Or(query.App.AppName.Eq(a.appName)).First()

	if err != nil {
		return nil
	}

	a.ap = ap

	return nil
}

func (a *App) FindAll() []*model.App {
	if len(a.aps) == 0 {
		return nil
	}

	return a.aps
}

// Add 添加应用
func (a *App) Add(ctx context.Context) (*model.App, error) {
	if a.appName == "" {
		return nil, microtype.AppParamCheckFailed
	}

	if a.ap != nil {
		return nil, microtype.AppExist
	}

	mo := &model.App{
		AppName: proto.String(a.appName),
	}

	err := query.App.WithContext(ctx).Create(mo)
	if err != nil {
		return nil, microtype.AppAddFailed
	}

	return mo, nil
}

// Query 查询应用
func (a *App) Query() *model.App {
	if a.ap != nil {
		return a.ap
	}

	return nil
}

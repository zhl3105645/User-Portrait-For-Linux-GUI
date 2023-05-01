package register

import (
	"backend/biz/entity/account"
	"backend/biz/entity/app"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/usecase/label"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"github.com/golang/protobuf/proto"
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
	//err := data_source.InitDataSource(ctx, r.app.AppID)
	//if err != nil {
	//	return err
	//}
	// 初始化标签
	parentMos := []*model.Label{
		&model.Label{
			LabelName: "基础信息",
			AppID:     r.app.AppID,
			FixType:   label.BasicInfo,
		},
		&model.Label{
			LabelName: "使用频率",
			AppID:     r.app.AppID,
			FixType:   label.UseFre,
		},
	}
	err := query.Label.WithContext(ctx).Create(parentMos...)
	if err != nil {
		return err
	}
	basicId, useFreId := parentMos[0].LabelID, parentMos[1].LabelID
	childMos := []*model.Label{
		&model.Label{
			LabelName:         "性别",
			IsLeaf:            1,
			DataType:          1,
			ParentLabelID:     proto.Int64(basicId),
			LabelSemanticDesc: proto.String("{\"1\":\"男\",\"2\":\"女\"}"),
			AppID:             r.app.AppID,
			FixType:           label.Gender,
		},
		&model.Label{
			LabelName:     "年龄",
			IsLeaf:        1,
			DataType:      2,
			ParentLabelID: proto.Int64(basicId),
			AppID:         r.app.AppID,
			FixType:       label.Age,
		},
		&model.Label{
			LabelName:     "职业",
			IsLeaf:        1,
			DataType:      2,
			ParentLabelID: proto.Int64(basicId),
			AppID:         r.app.AppID,
			FixType:       label.Career,
		},
		&model.Label{
			LabelName:     "平均使用时长",
			IsLeaf:        1,
			DataType:      2,
			ParentLabelID: proto.Int64(useFreId),
			AppID:         r.app.AppID,
			FixType:       label.UseTime,
		},
		&model.Label{
			LabelName:     "使用时间段",
			IsLeaf:        1,
			DataType:      2,
			ParentLabelID: proto.Int64(useFreId),
			AppID:         r.app.AppID,
			FixType:       label.UsePeriod,
		},
	}
	err = query.Label.WithContext(ctx).Create(childMos...)
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

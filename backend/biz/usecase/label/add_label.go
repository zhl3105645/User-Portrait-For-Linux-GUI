package label

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"strconv"
)

type AddLabel struct {
	accountId int64
	req       backend.AddLabelReq

	//
	appId    int64
	dataType int64
}

func NewAddLabel(accountId int64, req backend.AddLabelReq) *AddLabel {
	return &AddLabel{
		accountId: accountId,
		req:       req,
	}
}

func (a *AddLabel) Load(ctx context.Context) error {
	if !a.check() {
		return microtype.ParamCheckFailed
	}

	ac := account.NewAccount(a.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	a.appId = ac.GetQueryAccount().AppID

	// 添加
	mo := &model.Label{
		LabelID:           0,
		LabelName:         a.req.LabelName,
		IsLeaf:            int64(1),
		DataType:          a.dataType,
		ParentLabelID:     nil,
		LabelSemanticDesc: nil,
		AppID:             a.appId,
	}

	if a.req.ParentLabelID > 0 {
		mo.ParentLabelID = proto.Int64(a.req.ParentLabelID)
	}
	if len(a.req.ConvertRules) > 0 {
		res := make(map[string]string)
		for _, r := range a.req.ConvertRules {
			if r == nil {
				continue
			}
			res[r.Data] = r.Desc
		}
		desc, err := json.Marshal(res)
		if err != nil {
			return microtype.JsonMarshalFailed
		}
		mo.LabelSemanticDesc = proto.String(string(desc))
	}

	err := query.Label.WithContext(ctx).Create(mo)
	if err != nil {
		return microtype.LabelCreateFailed
	}

	// 更新父标签
	if a.req.ParentLabelID > 0 {
		queryMO, err := query.Label.WithContext(ctx).Where(query.Label.LabelID.Eq(a.req.ParentLabelID)).First()
		if err != nil {
			return microtype.LabelQueryFailed
		}

		updateMO := model.Label{
			LabelID:           queryMO.LabelID,
			LabelName:         queryMO.LabelName,
			IsLeaf:            0,
			DataType:          queryMO.DataType,
			ParentLabelID:     queryMO.ParentLabelID,
			LabelSemanticDesc: queryMO.LabelSemanticDesc,
			AppID:             queryMO.AppID,
		}
		_, err = query.Label.WithContext(ctx).Where(query.Label.LabelID.Eq(queryMO.LabelID)).Updates(updateMO)
		if err != nil {
			return microtype.LabelUpdateFailed
		}
	}

	return nil
}

func (a *AddLabel) GetResp() *backend.AddLabelResp {
	return &backend.AddLabelResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

func (a *AddLabel) check() bool {
	if a.req.LabelName == "" {
		return false
	}

	if a.req.DataType != "" {
		d, err := strconv.ParseInt(a.req.DataType, 10, 64)
		if err != nil {
			return false
		}

		a.dataType = d
	}

	return true
}

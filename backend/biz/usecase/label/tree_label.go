package label

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
	"strconv"
)

type TreeLabels struct {
	accountId int64
	needEnum  bool // 是否需要枚举值

	//
	appId int64
	res   *backend.TreeLabel
	// id -> label
	labelId2Label map[int64]*model.Label
	// parent -> children
	labelId2Children map[int64][]*model.Label
}

func NewTreeLabels(accountId int64, appId int64, needEnum bool) *TreeLabels {
	return &TreeLabels{
		accountId: accountId,
		appId:     appId,
		needEnum:  needEnum,
	}
}

func (t *TreeLabels) Load(ctx context.Context) error {
	if t.appId == 0 {
		ac := account.NewAccount(t.accountId, 0, "", "", 0, account.IdQuery)
		if err := ac.Load(ctx); err != nil {
			return err
		}

		t.appId = ac.GetQueryAccount().AppID
	}

	// 查询全部标签
	labels, err := query.Label.WithContext(ctx).Where(query.Label.AppID.Eq(t.appId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		return microtype.LabelQueryFailed
	}

	// id -> label
	t.labelId2Label = make(map[int64]*model.Label)
	// parent -> children
	t.labelId2Children = make(map[int64][]*model.Label)
	// 一级标签
	firstLabels := make([]*model.Label, 0)
	for _, l := range labels {
		t.labelId2Label[l.LabelID] = l
		if l.ParentLabelID != nil {
			parentId := *l.ParentLabelID
			if _, ok := t.labelId2Children[parentId]; ok {
				t.labelId2Children[parentId] = append(t.labelId2Children[parentId], l)
			} else {
				t.labelId2Children[parentId] = []*model.Label{l}
			}
		} else {
			firstLabels = append(firstLabels, l)
		}
	}

	// 计算最终结果
	t.res = &backend.TreeLabel{
		Name:     "标签",
		Value:    0,
		Children: make([]*backend.TreeLabel, 0),
	}

	// 递归计算
	for _, child := range firstLabels {
		t.res.Children = append(t.res.Children, t.getRecLabel(child))
	}

	return nil
}

func (t *TreeLabels) getRecLabel(root *model.Label) *backend.TreeLabel {
	res := &backend.TreeLabel{
		Name:     root.LabelName,
		Value:    root.LabelID,
		Children: make([]*backend.TreeLabel, 0),
	}
	//
	rootId := root.LabelID
	children := t.labelId2Children[rootId]
	if len(children) == 0 {
		// 叶子标签
		// 可枚举的变量
		if root.DataType == CanEnum && root.LabelSemanticDesc != nil {
			desc := make(map[string]string)
			err := json.Unmarshal([]byte(*root.LabelSemanticDesc), &desc)
			if err != nil {
				logger.Error("json unmarshal failed. err=", err.Error())
			} else if t.needEnum {
				labels := make([]*backend.TreeLabel, 0, len(desc))
				for k, v := range desc {
					n, _ := strconv.ParseInt(k, 10, 64)
					labels = append(labels, &backend.TreeLabel{
						Name:  v,
						Value: n,
					})
				}
				res.Children = labels
			}
		}

		return res
	}

	// 非叶子标签
	for _, child := range children {
		res.Children = append(res.Children, t.getRecLabel(child))
	}

	return res
}

func (t *TreeLabels) GetResp() *backend.TreeLabelResp {
	return &backend.TreeLabelResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Data:       t.res,
	}
}

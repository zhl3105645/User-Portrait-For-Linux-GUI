package profile

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/usecase/label"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
)

type Profile struct {
	userId int64

	//
	appId         int64
	root          *backend.TreeLabel
	labelId2Data  map[int64]string
	labelId2Label map[int64]*model.Label
}

func NewProfile(userId int64) *Profile {
	return &Profile{
		userId: userId,
	}
}

func (p *Profile) Load(ctx context.Context) error {
	// 用户信息
	user, err := query.User.WithContext(ctx).Where(query.User.UserID.Eq(p.userId)).First()
	if err != nil {
		return microtype.UserQueryFailed
	}

	p.appId = user.AppID

	// 树状标签
	lb := label.NewTreeLabels(0, p.appId, false)
	if err := lb.Load(ctx); err != nil {
		return err
	}

	p.root = lb.GetResp().Data
	if p.root == nil {
		return nil
	}

	p.root.Name = "画像"

	// 叶子标签
	leafLabel, err := query.Label.WithContext(ctx).
		Where(query.Label.AppID.Eq(p.appId), query.Label.IsLeaf.Eq(1)).
		Find()
	if err != nil {
		return microtype.LabelQueryFailed
	}
	p.labelId2Label = make(map[int64]*model.Label)

	// 叶子标签数据
	leafIds := make([]int64, 0, len(leafLabel))
	for _, lab := range leafLabel {
		if lab == nil {
			continue
		}
		leafIds = append(leafIds, lab.LabelID)
		p.labelId2Label[lab.LabelID] = lab
	}

	labelData, err := query.LabelDatum.WithContext(ctx).
		Where(query.LabelDatum.UserID.Eq(p.userId), query.LabelDatum.LabelID.In(leafIds...)).
		Find()
	if err != nil {
		return microtype.LabelDataQueryFailed
	}

	p.labelId2Data = make(map[int64]string)
	for _, d := range labelData {
		if d == nil {
			continue
		}
		p.labelId2Data[d.LabelID] = d.Data
	}

	// 添加标签值
	p.rec(p.root)

	return nil
}

func (p *Profile) rec(root *backend.TreeLabel) bool {
	if len(root.Children) == 0 {
		lab := p.labelId2Label[root.Value]
		data := p.labelId2Data[root.Value]
		if lab == nil {
			return false
		}

		if lab.DataType == label.CanEnum && lab.LabelSemanticDesc != nil {
			descMap := make(map[string]string)
			err := json.Unmarshal([]byte(*lab.LabelSemanticDesc), &descMap)
			if err != nil {
				return false
			}
			data = descMap[data]
		}
		if data == "" {
			return false
		}

		root.Name = root.Name + ": " + data
		return true
	}

	newChild := make([]*backend.TreeLabel, 0)
	for _, child := range root.Children {
		if p.rec(child) {
			newChild = append(newChild, child)
		}
	}
	if len(newChild) == 0 {
		return false
	}

	root.Children = newChild

	return true
}

func (p *Profile) GetResp() *backend.ProfileResp {
	return &backend.ProfileResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Label:      p.root,
	}
}

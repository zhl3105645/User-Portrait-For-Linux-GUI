package profile

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/usecase/label"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
)

type Profile struct {
	userId int64

	//
	appId         int64
	user          *model.User
	root          *backend.TreeLabel
	labelId2Data  map[int64]string
	labelId2Label map[int64]*model.Label
	radars        []*backend.Radar
	groupLabels   []*backend.GroupLabel
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
	p.user = user

	//// 树状标签
	//lb := label.NewTreeLabels(0, p.appId, false)
	//if err := lb.Load(ctx); err != nil {
	//	return err
	//}
	//
	//p.root = lb.GetResp().Data
	//if p.root == nil {
	//	return nil
	//}
	//
	//p.root.Name = "画像"
	//
	//// 叶子标签
	//leafLabel, err := query.Label.WithContext(ctx).
	//	Where(query.Label.AppID.Eq(p.appId), query.Label.IsLeaf.Eq(1)).
	//	Find()
	//if err != nil {
	//	return microtype.LabelQueryFailed
	//}
	//p.labelId2Label = make(map[int64]*model.Label)
	//
	//// 叶子标签数据
	//leafIds := make([]int64, 0, len(leafLabel))
	//for _, lab := range leafLabel {
	//	if lab == nil {
	//		continue
	//	}
	//	leafIds = append(leafIds, lab.LabelID)
	//	p.labelId2Label[lab.LabelID] = lab
	//}
	//
	//labelData, err := query.LabelDatum.WithContext(ctx).
	//	Where(query.LabelDatum.UserID.Eq(p.userId), query.LabelDatum.LabelID.In(leafIds...)).
	//	Find()
	//if err != nil {
	//	return microtype.LabelDataQueryFailed
	//}
	//
	//p.labelId2Data = make(map[int64]string)
	//for _, d := range labelData {
	//	if d == nil {
	//		continue
	//	}
	//	p.labelId2Data[d.LabelID] = d.Data
	//}
	//
	//// 添加标签值
	//p.rec(p.root)

	// 行为时长雷达图
	p.radars = p.getRadar(ctx)

	// 标签数据
	p.groupLabels = p.getGroupLabel(ctx)

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

func (p *Profile) getGroupLabel(ctx context.Context) []*backend.GroupLabel {
	// 叶子标签
	labs, err := query.Label.WithContext(ctx).Where(query.Label.IsLeaf.Eq(1)).Find()
	if err != nil {
		return nil
	}

	labIds := make([]int64, 0, len(labs))
	parentIds := make([]int64, 0)
	parentId2children := make(map[int64][]*model.Label)
	for _, lab := range labs {
		if lab == nil {
			continue
		}
		labIds = append(labIds, lab.LabelID)
		if lab.ParentLabelID != nil {
			if _, ok := parentId2children[*lab.ParentLabelID]; ok {
				parentId2children[*lab.ParentLabelID] = append(parentId2children[*lab.ParentLabelID], lab)
			} else {
				parentIds = append(parentIds, *lab.ParentLabelID)
				parentId2children[*lab.ParentLabelID] = []*model.Label{lab}
			}
		}
	}

	// 标签数据
	labData, err := query.LabelDatum.WithContext(ctx).Where(query.LabelDatum.LabelID.In(labIds...), query.LabelDatum.UserID.Eq(p.userId)).Find()
	if err != nil {
		return nil
	}
	labelId2Data := make(map[int64]string) // 一个标签只存在一个数据
	for _, d := range labData {
		labelId2Data[d.LabelID] = d.Data
	}

	// 父标签
	parentLabs, err := query.Label.WithContext(ctx).Where(query.Label.LabelID.In(parentIds...)).Find()
	if err != nil {
		return nil
	}

	res := make([]*backend.GroupLabel, 0)
	for _, parentLab := range parentLabs {
		if parentLab == nil {
			continue
		}
		m := &backend.GroupLabel{
			ParentLabelName: parentLab.LabelName,
			Labels:          make([]*backend.LabelValue, 0),
		}
		for _, lab := range parentId2children[parentLab.LabelID] {
			value := labelId2Data[lab.LabelID]
			if lab.DataType == label.CanEnum && lab.LabelSemanticDesc != nil {
				descMap := make(map[string]string)
				err := json.Unmarshal([]byte(*lab.LabelSemanticDesc), &descMap)
				if err == nil {
					value = descMap[value]
				}
			}

			m.Labels = append(m.Labels, &backend.LabelValue{
				LabelID:    lab.LabelID,
				LabelName:  lab.LabelName,
				LabelValue: value,
			})
		}
		res = append(res, m)
	}

	return res
}

func (p *Profile) getRadar(ctx context.Context) []*backend.Radar {
	// 应用 最大 & 平均 & 用户当前
	ap, err := query.App.WithContext(ctx).Where(query.App.AppID.Eq(p.appId)).First()
	if err != nil {
		logger.Error("query app failed. err=", err.Error())
		return nil
	}

	if ap.MaxBehaviorDurationMap == nil || ap.AveBehaviorDurationMap == nil || p.user.BehaviorDurationMap == nil {
		logger.Error("value is empty")
		return nil
	}

	maxMap := make(map[int64]int64)
	aveMap := make(map[int64]int64)
	userMap := make(map[int64]int64)

	err = json.Unmarshal([]byte(*ap.MaxBehaviorDurationMap), &maxMap)
	if err != nil {
		logger.Error("json unmarshal failed.", err.Error())
		return nil
	}
	err = json.Unmarshal([]byte(*ap.AveBehaviorDurationMap), &aveMap)
	if err != nil {
		logger.Error("json unmarshal failed.", err.Error())
		return nil
	}
	err = json.Unmarshal([]byte(*p.user.BehaviorDurationMap), &userMap)
	if err != nil {
		logger.Error("json unmarshal failed.", err.Error())
		return nil
	}

	// 规则描述
	rules, err := query.Rule.WithContext(ctx).Where(query.Rule.AppID.Eq(p.appId)).Find()
	if err != nil {
		logger.Error("query rules failed.", err.Error())
		return nil
	}
	ruleDescMap := make(map[int64]string)
	for _, r := range rules {
		if r == nil {
			continue
		}
		ruleDescMap[r.RuleID] = r.RuleDesc
	}

	res := make([]*backend.Radar, 0)
	for ruleId, duration := range userMap {
		if ruleId <= 0 {
			continue
		}

		res = append(res, &backend.Radar{
			Name: ruleDescMap[ruleId],
			Max:  maxMap[ruleId],
			Cur:  duration,
			Ave:  aveMap[ruleId],
		})
	}

	return res
}

func (p *Profile) GetResp() *backend.ProfileResp {
	return &backend.ProfileResp{
		StatusCode:  microtype.SuccessErr.Code,
		StatusMsg:   microtype.SuccessErr.Msg,
		Label:       p.root,
		Radars:      p.radars,
		GroupLabels: p.groupLabels,
	}
}

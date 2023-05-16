package profile

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/usecase/label"
	"backend/biz/util"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"strconv"
)

type GroupProfile struct {
	crowdId   int64
	accountId int64

	//
	appId               int64
	userIds             []int64
	stackBarModelLabels []*model.Label
	pieModelLabels      []*model.Label
	barModelLabel       *model.Label

	labelId2UserId2Data map[int64]map[int64]string

	//
	radars        []*backend.Radar
	pieLabels     []*backend.PieLabel
	stackBarLabel *backend.StackBarLabel
	barLabel      *backend.BarLabel
}

func NewGroupProfile(accountId int64, crowdId int64) *GroupProfile {
	return &GroupProfile{
		crowdId:             crowdId,
		accountId:           accountId,
		userIds:             make([]int64, 0),
		stackBarModelLabels: make([]*model.Label, 0),
		pieModelLabels:      make([]*model.Label, 0),
		labelId2UserId2Data: make(map[int64]map[int64]string),
	}
}

func (g *GroupProfile) Load(ctx context.Context) error {
	ac := account.NewAccount(g.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	g.appId = ac.GetQueryAccount().AppID

	// 确定人群所属用户
	if g.crowdId == -1 {
		users, err := query.User.WithContext(ctx).Where(query.User.AppID.Eq(g.appId)).Find()
		if err != nil {
			return err
		}
		for _, u := range users {
			if u == nil {
				continue
			}
			g.userIds = append(g.userIds, u.UserID)
		}
	} else {
		relations, err := query.CrowdRelation.WithContext(ctx).Where(query.CrowdRelation.CrowdID.Eq(g.crowdId)).Find()
		if err != nil {
			return err
		}
		for _, u := range relations {
			if u == nil {
				continue
			}
			g.userIds = append(g.userIds, u.UserID)
		}
	}

	if len(g.userIds) < 0 {
		return nil
	}

	// 叶子标签
	labs, err := query.Label.WithContext(ctx).Where(query.Label.AppID.Eq(g.appId), query.Label.IsLeaf.Eq(1)).Find()
	if err != nil {
		return nil
	}

	labIds := make([]int64, 0, len(labs))
	for _, lab := range labs {
		if lab == nil {
			continue
		}
		labIds = append(labIds, lab.LabelID)
		if lab.FixType == label.UsePeriod {
			g.barModelLabel = lab
		} else if lab.FixType > 0 {
			g.pieModelLabels = append(g.pieModelLabels, lab)
		} else if lab.DataType == label.CanEnum {
			g.stackBarModelLabels = append(g.stackBarModelLabels, lab)
		}
	}

	// 标签数据
	labData, err := query.LabelDatum.WithContext(ctx).Where(query.LabelDatum.LabelID.In(labIds...), query.LabelDatum.UserID.In(g.userIds...)).Find()
	if err != nil {
		return nil
	}
	for _, d := range labData {
		if _, ok := g.labelId2UserId2Data[d.LabelID]; !ok {
			g.labelId2UserId2Data[d.LabelID] = make(map[int64]string)
		}
		g.labelId2UserId2Data[d.LabelID][d.UserID] = d.Data
	}

	// 确定数据
	g.radars = g.getRadar(ctx)
	g.stackBarLabel = g.getStackBarLabel()
	g.pieLabels = g.getPieLabels()
	g.barLabel = g.getBarLabel()

	return nil
}

func (g *GroupProfile) getBarLabel() *backend.BarLabel {
	if g.barModelLabel == nil || g.barModelLabel.FixType != label.UsePeriod {
		return nil
	}

	// 数据
	data, ok := g.labelId2UserId2Data[g.barModelLabel.LabelID]
	if !ok {
		return nil
	}

	// 各小时分布
	hourCntMap := make(map[int64]int64)
	min, max := int64(24), int64(0)
	for _, d := range data {
		beginHour, endHour, err := util.GetBeginHourAndEndHour(d)
		if err != nil {
			continue
		}
		for i := beginHour; i < endHour; i++ {
			if cnt, ok := hourCntMap[i]; ok {
				hourCntMap[i] = cnt + 1
			} else {
				hourCntMap[i] = 1
			}
			if beginHour < min {
				min = beginHour
			}
			if endHour > max {
				max = endHour
			}
		}
	}

	res := &backend.BarLabel{
		XNames: make([]string, 0),
		Data:   make([]int64, 0),
	}

	for i := min; i <= max; i++ {
		res.XNames = append(res.XNames, strconv.FormatInt(i, 10))
		res.Data = append(res.Data, hourCntMap[i])
	}

	return res
}

func (g *GroupProfile) getPieLabels() []*backend.PieLabel {
	res := make([]*backend.PieLabel, 0, len(g.pieModelLabels))

	for _, lab := range g.pieModelLabels {
		data, ok := g.labelId2UserId2Data[lab.LabelID]
		if !ok {
			continue
		}

		pieData := make([]*backend.PieData, 0)
		if lab.FixType == label.Gender {
			descMap := make(map[string]string)
			err := json.Unmarshal([]byte(*lab.LabelSemanticDesc), &descMap)
			if err != nil {
				continue
			}

			dataCntMap := make(map[string]int64)
			for _, d := range data {
				if cnt, ok := dataCntMap[d]; ok {
					dataCntMap[d] = cnt + 1
				} else {
					dataCntMap[d] = 1
				}
			}

			for d, cnt := range dataCntMap {
				if cnt <= 0 {
					continue
				}
				pieData = append(pieData, &backend.PieData{
					Name:  descMap[d],
					Value: cnt,
				})
			}
		} else if lab.FixType == label.Age {
			descs := []string{"20岁以下", "20~23岁", "24~27岁", "27~30岁", "30岁以上"}
			cnts := make([]int64, 5)
			for _, d := range data {
				age, err := strconv.ParseInt(d, 10, 64)
				if err != nil {
					continue
				}
				if age < 20 {
					cnts[0] = cnts[0] + 1
				} else if age <= 23 {
					cnts[1] = cnts[1] + 1
				} else if age <= 27 {
					cnts[2] = cnts[2] + 1
				} else if age <= 30 {
					cnts[3] = cnts[3] + 1
				} else {
					cnts[4] = cnts[4] + 1
				}
			}
			for idx, cnt := range cnts {
				if cnt <= 0 {
					continue
				}
				pieData = append(pieData, &backend.PieData{
					Name:  descs[idx],
					Value: cnt,
				})
			}
		} else if lab.FixType == label.Career {
			cntMap := make(map[string]int64)
			for _, d := range data {
				if cnt, ok := cntMap[d]; ok {
					cntMap[d] = cnt + 1
				} else {
					cntMap[d] = 1
				}
			}
			for desc, cnt := range cntMap {
				pieData = append(pieData, &backend.PieData{
					Name:  desc,
					Value: cnt,
				})
			}
		} else if lab.FixType == label.UseTime {
			descs := []string{"0.5h以下", "0.5h~1h", "1h~1.5h", "1.5h~2h", "2h以上"}
			cnts := make([]int64, 5)
			for _, d := range data {
				timeMS, err := util.ParseTimeDurationFromStr(d)
				if err != nil {
					continue
				}
				hour := float64(timeMS) / float64(1000*60*60)
				if hour < 0.5 {
					cnts[0] = cnts[0] + 1
				} else if hour <= 1 {
					cnts[1] = cnts[1] + 1
				} else if hour <= 1.5 {
					cnts[2] = cnts[2] + 1
				} else if hour <= 2 {
					cnts[3] = cnts[3] + 1
				} else {
					cnts[4] = cnts[4] + 1
				}
			}
			for idx, cnt := range cnts {
				if cnt <= 0 {
					continue
				}
				pieData = append(pieData, &backend.PieData{
					Name:  descs[idx],
					Value: cnt,
				})
			}
		}

		if len(pieData) == 0 {
			continue
		}

		res = append(res, &backend.PieLabel{
			LabelName: lab.LabelName,
			Data:      pieData,
		})
	}

	return res
}

func (g *GroupProfile) getStackBarLabel() *backend.StackBarLabel {
	res := &backend.StackBarLabel{
		LabelNames:     make([]string, 0),
		LabelCnt:       make([][]int64, 0),
		LabelValueDesc: make([]string, 0),
	}

	// label_id -> value_desc -> cnt
	labelId2Desc2Cnt := make(map[int64]map[string]int64)
	descCnt := 0
	for _, lab := range g.stackBarModelLabels {
		if lab == nil || lab.LabelSemanticDesc == nil {
			continue
		}

		data, ok := g.labelId2UserId2Data[lab.LabelID]
		if !ok {
			continue
		}

		descMap := make(map[string]string)
		err := json.Unmarshal([]byte(*lab.LabelSemanticDesc), &descMap)
		if err != nil {
			return nil
		}

		dataCntMap := make(map[string]int64)
		for _, d := range data {
			if cnt, ok := dataCntMap[d]; ok {
				dataCntMap[d] = cnt + 1
			} else {
				dataCntMap[d] = 1
			}
		}

		desc2Cnt := make(map[string]int64)
		for d, cnt := range dataCntMap {
			desc2Cnt[descMap[d]] = cnt
			descCnt++
		}
		labelId2Desc2Cnt[lab.LabelID] = desc2Cnt
	}

	curIndex := 0
	for _, lab := range g.stackBarModelLabels {
		desc2Cnt := labelId2Desc2Cnt[lab.LabelID]
		if len(desc2Cnt) == 0 {
			continue
		}
		res.LabelNames = append(res.LabelNames, lab.LabelName)
		for desc, cnt := range desc2Cnt {
			if cnt == 0 {
				continue
			}
			res.LabelCnt = append(res.LabelCnt, make([]int64, len(labelId2Desc2Cnt)))
			res.LabelValueDesc = append(res.LabelValueDesc, desc)
			res.LabelCnt[len(res.LabelCnt)-1][curIndex] = cnt
		}
		curIndex += 1
	}

	return res
}

func (g *GroupProfile) getRadar(ctx context.Context) []*backend.Radar {
	// 应用 最大 & 平均 & 人群当前(可能不存在)
	ap, err := query.App.WithContext(ctx).Where(query.App.AppID.Eq(g.appId)).First()
	if err != nil {
		logger.Error("query app failed. err=", err.Error())
		return nil
	}

	if ap.MaxBehaviorDurationMap == nil || ap.AveBehaviorDurationMap == nil {
		logger.Error("value is empty")
		return nil
	}

	maxMap := make(map[int64]int64)
	aveMap := make(map[int64]int64)
	crowdMap := make(map[int64]int64)

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

	if g.crowdId > 0 {
		crowd, err := query.Crowd.WithContext(ctx).Where(query.Crowd.CrowdID.Eq(g.crowdId)).First()
		if err != nil {
			logger.Error("query crowd failed.", err.Error())
			return nil
		}
		if crowd.BehaviorDurationMap != nil {
			err = json.Unmarshal([]byte(*crowd.BehaviorDurationMap), &crowdMap)
			if err != nil {
				logger.Error("json unmarshal failed.", err.Error())
				return nil
			}
		}
	}

	// 规则描述
	rules, err := query.Rule.WithContext(ctx).Where(query.Rule.AppID.Eq(g.appId)).Find()
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
	for ruleId, max := range maxMap {
		if ruleId <= 0 {
			continue
		}

		res = append(res, &backend.Radar{
			Name: ruleDescMap[ruleId],
			Max:  max,
			Cur:  crowdMap[ruleId],
			Ave:  aveMap[ruleId],
		})
	}

	return res
}

func (g *GroupProfile) GetResp() *backend.GroupProfileResp {
	return &backend.GroupProfileResp{
		StatusCode:    microtype.SuccessErr.Code,
		StatusMsg:     microtype.SuccessErr.Msg,
		Radars:        g.radars,
		PieLabel:      g.pieLabels,
		StackBarLabel: g.stackBarLabel,
		BarLabel:      g.barLabel,
	}
}

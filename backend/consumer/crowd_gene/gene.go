package crowd_gene

import (
	"backend/biz/model/backend"
	"backend/biz/usecase/crowd"
	"backend/cmd/dal"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
)

func Gene(crowdId int64) {
	ctx := context.Background()

	cr, err := query.Crowd.WithContext(ctx).Where(query.Crowd.CrowdID.Eq(crowdId)).First()
	if err != nil {
		logger.Error("query crowd failed. err=", err.Error())
		return
	}

	// 解析规则
	divideRules := make([]*backend.DivideRule, 0)
	err = json.Unmarshal([]byte(cr.CrowdDivideRule), &divideRules)
	if err != nil {
		logger.Error("json unmarshal failed. err=", err.Error())
		return
	}

	labelIds := make([]int64, 0, len(divideRules))
	for _, rule := range divideRules {
		labelIds = append(labelIds, rule.LabelID)
	}

	// 用户标签数据
	labelData, err := query.LabelDatum.WithContext(ctx).Where(query.LabelDatum.LabelID.In(labelIds...)).Find()
	if err != nil {
		logger.Error("label data query failed. err=", err.Error())
		return
	}

	label2UserId2Data := make(map[int64]map[int64]string) // label_id -> user_id -> data
	for _, data := range labelData {
		labelId := data.LabelID
		if _, ok := label2UserId2Data[labelId]; !ok {
			label2UserId2Data[labelId] = make(map[int64]string)
		}
		label2UserId2Data[labelId][data.UserID] = data.Data
	}

	// 每个条件满足的用户 && 对用户做并集 or 交集
	userIds := make([]int64, 0)
	for idx, rule := range divideRules {
		us := getMatchUsers(rule, label2UserId2Data[rule.LabelID])
		if idx == 0 {
			userIds = us
		} else {
			userIds = crowd.Union(rule.UnionOperate, userIds, us)
		}
	}

	q := query.Use(dal.DB)

	err = q.Transaction(func(tx *query.Query) error {
		// 删除已有关系
		_, err = tx.CrowdRelation.WithContext(ctx).Where(tx.CrowdRelation.CrowdID.Eq(crowdId)).Delete()
		if err != nil {
			logger.Error("crowd relation delete failed. err=", err.Error())
			return err
		}

		// 写入新关系
		mos := make([]*model.CrowdRelation, 0, len(userIds))
		for _, userId := range userIds {
			mos = append(mos, &model.CrowdRelation{
				CrowdRelationID: 0,
				UserID:          userId,
				CrowdID:         crowdId,
			})
		}

		err = tx.CrowdRelation.WithContext(ctx).Create(mos...)
		if err != nil {
			logger.Error("crowd relation create failed. err=", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	// 行为数据更新
	users, err := query.User.WithContext(ctx).Where(query.User.UserID.In(userIds...)).Find()
	if err != nil {
		logger.Error("query user data failed. err=", err.Error())
		return
	}
	aveDurationMap := make(map[int64]int64)
	total := int64(0)
	for _, u := range users {
		if u == nil || u.BehaviorDurationMap == nil {
			continue
		}

		durationMap := make(map[int64]int64)
		err = json.Unmarshal([]byte(*u.BehaviorDurationMap), &durationMap)
		if err != nil {
			logger.Error("json unmarshal failed. err=", err.Error())
			continue
		}

		total++
		for id, duration := range durationMap {
			if cnt, ok := aveDurationMap[id]; ok {
				aveDurationMap[id] = cnt + duration
			} else {
				aveDurationMap[id] = duration
			}
		}
	}
	for id, duration := range aveDurationMap {
		aveDurationMap[id] = duration / total
	}

	bs, err := json.Marshal(aveDurationMap)
	if err != nil {
		logger.Error("json marshal failed. err=", err.Error())
		return
	}

	mo := model.Crowd{
		BehaviorDurationMap: proto.String(string(bs)),
	}
	_, err = query.Crowd.WithContext(ctx).Where(query.Crowd.CrowdID.Eq(crowdId)).Updates(mo)
	if err != nil {
		logger.Error("update crowd failed. err=", err.Error())
		return
	}

	return
}

func getMatchUsers(rule *backend.DivideRule, userId2Data map[int64]string) []int64 {
	res := make([]int64, 0)
	for userId, data := range userId2Data {
		if crowd.MatchDivide(rule.DivideOperate, data, rule.LabelData) {
			res = append(res, userId)
		}
	}

	return res
}

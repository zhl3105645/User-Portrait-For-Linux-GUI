package rule_gene

import (
	"backend/biz/entity/rule"
	"backend/cmd/dal/query"
	"context"
)

type EventRule struct {
	RuleID   int64
	RuleDesc string
	Elements []*rule.EventRuleElement // 事件
}

type BehaviorRule struct {
	RuleID   int64
	RuleDesc string
	Elements []*rule.BehaviorRuleElement // 事件规则ID
}

type result struct {
	RuleElementID    int64  `gorm:"column:rule_element_id;type:bigint;primaryKey;autoIncrement:true" json:"rule_element_id"` // 规则元素ID
	RuleElementValue string `gorm:"column:rule_element_value;type:text;not null" json:"rule_element_value"`                  // 规则元素值
	RuleID           int64  `gorm:"column:rule_id;type:bigint;not null" json:"rule_id"`
	RuleType         int64  `gorm:"column:rule_type;type:int;not null" json:"rule_type"`  // 规则类型
	RuleDesc         string `gorm:"column:rule_desc;type:text;not null" json:"rule_desc"` // 规则描述
}

func getRules(ctx context.Context, appId int64) ([]*EventRule, []*BehaviorRule, error) {
	ruleDO := query.Rule
	//ruleMO := ruleDO.WithContext(ctx)
	elementDO := query.RuleElement
	elementMO := elementDO.WithContext(ctx)

	res := make([]*result, 0)
	err := elementMO.Select(elementDO.ALL, ruleDO.RuleType, ruleDO.RuleDesc).
		Join(ruleDO, ruleDO.RuleID.EqCol(elementDO.RuleID)).
		Where(ruleDO.AppID.Eq(appId)).Scan(&res)
	if err != nil {
		return nil, nil, err
	}

	if len(res) == 0 {
		return nil, nil, err
	}

	resultMap := make(map[int64][]*result)
	for _, r := range res {
		if r == nil {
			continue
		}

		if v, ok := resultMap[r.RuleID]; !ok || v == nil {
			resultMap[r.RuleID] = make([]*result, 0)

		}
		resultMap[r.RuleID] = append(resultMap[r.RuleID], r)
	}

	eventRules := make([]*EventRule, 0)
	behaviorRules := make([]*BehaviorRule, 0)
	for rId, results := range resultMap {
		if len(results) == 0 {
			continue
		}

		ruleType := results[0].RuleType
		if ruleType == int64(rule.EventRule) {
			event := &EventRule{
				RuleID:   rId,
				RuleDesc: results[0].RuleDesc,
				Elements: nil,
			}

			eles := make([]*rule.EventRuleElement, 0, len(results))
			for _, result := range results {
				if result == nil {
					continue
				}

				ele := rule.ParseEventElement(result.RuleElementValue)
				if ele == nil {
					continue
				}

				eles = append(eles, ele)
			}

			event.Elements = eles
			eventRules = append(eventRules, event)
		} else if ruleType == int64(rule.BehaviorRule) {
			behavior := &BehaviorRule{
				RuleID:   rId,
				RuleDesc: results[0].RuleDesc,
				Elements: nil,
			}

			eles := make([]*rule.BehaviorRuleElement, 0, len(results))
			for _, result := range results {
				if result == nil {
					continue
				}

				ele := rule.ParseBehaviorElement(result.RuleElementValue)
				if ele == nil {
					continue
				}

				eles = append(eles, ele)
			}

			behavior.Elements = eles
			behaviorRules = append(behaviorRules, behavior)
		}
	}

	// 行为规则 未操作
	if len(behaviorRules) > 0 {
		behaviorRules = append(behaviorRules, &BehaviorRule{
			RuleID:   rule.BehaviorRuleNoOperate,
			RuleDesc: "未操作",
			Elements: []*rule.BehaviorRuleElement{
				&rule.BehaviorRuleElement{
					EventRuleIds: []int64{-1, -2},
				},
			},
		})
	}

	return eventRules, behaviorRules, nil
}

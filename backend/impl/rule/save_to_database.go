package rule

import (
	"backend/cmd/dal"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"strings"
)

func LoadRuleToDatabase(ctx context.Context) {
	dal.Init()
	// 需要提前清空 rule rule_element 表
	loadEventRuleToDatabase(ctx)
	loadBehaviorRuleDatabase(ctx)
}

func loadEventRuleToDatabase(ctx context.Context) {
	appId := int64(2)

	ruleDO := query.Rule
	ruleMO := ruleDO.WithContext(ctx)
	elementDO := query.RuleElement
	elementMO := elementDO.WithContext(ctx)

	eventRules, _ := GetRules()

	rules := make([]*model.Rule, 0, len(eventRules))
	ruleName2Rule := make(map[string]*model.Rule)

	for _, eventRule := range eventRules {
		if eventRule == nil {
			continue
		}
		r := &model.Rule{
			RuleType: 1,
			RuleDesc: eventRule.Name,
			AppID:    appId,
		}
		rules = append(rules, r)
		ruleName2Rule[r.RuleDesc] = r
	}

	err := ruleMO.Create(rules...)
	if err != nil {
		logger.Error("rule err=", err.Error())
		return
	}

	elements := make([]*model.RuleElement, 0)
	for _, eventRule := range eventRules {
		if eventRule == nil {
			continue
		}

		r := ruleName2Rule[eventRule.Name]
		if ruleMO == nil {
			continue
		}

		for _, s := range eventRule.Events {
			elements = append(elements, &model.RuleElement{
				RuleElementID:    0,
				RuleElementValue: s,
				RuleID:           r.RuleID,
			})
		}
	}

	err = elementMO.Create(elements...)
	if err != nil {
		logger.Error("element err=", err.Error())
		return
	}
}

func loadBehaviorRuleDatabase(ctx context.Context) {
	appId := int64(2)

	ruleDO := query.Rule
	ruleMO := ruleDO.WithContext(ctx)
	elementDO := query.RuleElement
	elementMO := elementDO.WithContext(ctx)

	_, behaviorRules := GetRules()

	rules := make([]*model.Rule, 0, len(behaviorRules))
	ruleName2Rule := make(map[string]*model.Rule)

	for _, behaviorRule := range behaviorRules {
		if behaviorRule == nil {
			continue
		}
		r := &model.Rule{
			RuleType: 2,
			RuleDesc: behaviorRule.Name,
			AppID:    appId,
		}
		rules = append(rules, r)
		ruleName2Rule[r.RuleDesc] = r
	}

	err := ruleMO.Create(rules...)
	if err != nil {
		logger.Error("rule err=", err.Error())
		return
	}

	elements := make([]*model.RuleElement, 0)
	for _, behaviorRule := range behaviorRules {
		if behaviorRule == nil {
			continue
		}

		r := ruleName2Rule[behaviorRule.Name]
		if ruleMO == nil {
			continue
		}

		for _, s := range behaviorRule.Behaviors {
			if strings.Contains(s, "-") {
				continue
			}
			elements = append(elements, &model.RuleElement{
				RuleElementID:    0,
				RuleElementValue: s,
				RuleID:           r.RuleID,
			})
		}
	}

	err = elementMO.Create(elements...)
	if err != nil {
		logger.Error("element err=", err.Error())
		return
	}
}

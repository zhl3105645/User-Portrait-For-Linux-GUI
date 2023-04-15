package element

import (
	"backend/biz/entity/account"
	"backend/biz/entity/rule"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"gorm.io/gorm"
)

type PageElement struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string
	ruleType  int64

	//
	appId int64
	resMo []*result
	total int64
	//
	events    []*backend.EventRuleElement
	behaviors []*backend.BehaviorRuleElement
}

func NewPageElement(ruleType int64, accountId int64, pageNum int64, pageSize int64, search string) *PageElement {
	return &PageElement{
		ruleType:  ruleType,
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

type result struct {
	RuleID           int64  `gorm:"column:rule_id;type:bigint;primaryKey;autoIncrement:true" json:"rule_id"` // 规则ID
	RuleType         int64  `gorm:"column:rule_type;type:int;not null" json:"rule_type"`                     // 规则类型
	RuleDesc         string `gorm:"column:rule_desc;type:text" json:"rule_desc"`                             // 规则描述
	AppID            int64  `gorm:"column:app_id;type:bigint;not null" json:"app_id"`
	RuleElementID    int64  `gorm:"column:rule_element_id;type:bigint;primaryKey;autoIncrement:true" json:"rule_element_id"` // 规则元素ID
	RuleElementValue string `gorm:"column:rule_element_value;type:text;not null" json:"rule_element_value"`                  // 规则元素值
}

func (p *PageElement) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	ruleDO := query.Rule
	ruleMO := ruleDO.WithContext(context.Background())
	elementDO := query.RuleElement

	res := make([]*result, 0)

	offset := (p.pageNum - 1) * p.pageSize
	count, err := ruleMO.Select(ruleDO.ALL, elementDO.RuleElementID, elementDO.RuleElementValue).LeftJoin(elementDO, ruleDO.RuleID.EqCol(elementDO.RuleID)).
		Where(ruleDO.AppID.Eq(p.appId), ruleDO.RuleType.Eq(p.ruleType), ruleDO.RuleDesc.Like("%"+p.search+"%")).
		ScanByPage(&res, int(offset), int(p.pageSize))
	if err != nil {
		return err
	}
	p.resMo = res
	p.total = count

	//
	p.events = make([]*backend.EventRuleElement, 0)
	p.behaviors = make([]*backend.BehaviorRuleElement, 0)

	for _, v := range p.resMo {
		if v == nil || v.RuleID <= 0 {
			continue
		}

		if v.RuleType == int64(rule.EventRule) {
			re := &backend.EventRuleElement{
				RuleID:    v.RuleID,
				RuleType:  v.RuleType,
				RuleDesc:  v.RuleDesc,
				ElementID: v.RuleElementID,
			}

			event := rule.ParseEventElement(v.RuleElementValue)
			if event != nil {
				re.EventType = event.EventType
				re.MouseClickType = event.MouseClickType
				re.MouseClickButton = event.MouseClickButton
				re.KeyClickType = event.KeyClickType
				re.KeyValue = event.KeyValue
				re.ComponentNamePrefix = event.ComponentNamePrefix
			}

			p.events = append(p.events, re)
		} else if v.RuleType == int64(rule.BehaviorRule) {
			re := &backend.BehaviorRuleElement{
				RuleID:     v.RuleID,
				RuleType:   v.RuleType,
				RuleDesc:   v.RuleDesc,
				ElementID:  v.RuleElementID,
				EventRules: make([]*backend.EventRule, 0),
			}

			behavior := rule.ParseBehaviorElement(v.RuleElementValue)

			if behavior != nil && len(behavior.EventRuleIds) > 0 {
				rules, err := ruleMO.Where(ruleDO.RuleID.In(behavior.EventRuleIds...)).Find()
				if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
					return microtype.RuleQueryFailed
				}

				mp := make(map[int64]*model.Rule)
				for _, r := range rules {
					if r == nil {
						continue
					}

					mp[r.RuleID] = r
				}

				for _, id := range behavior.EventRuleIds {
					if r, ok := mp[id]; ok && r != nil {
						re.EventRules = append(re.EventRules, &backend.EventRule{
							RuleID:   r.RuleID,
							RuleDesc: r.RuleDesc,
						})
					}
				}
			}

			p.behaviors = append(p.behaviors, re)
		}
	}

	return nil
}

func (p *PageElement) GetResp() *backend.ElementInPageResp {
	return &backend.ElementInPageResp{
		StatusCode:       microtype.SuccessErr.Code,
		StatusMsg:        microtype.SuccessErr.Msg,
		EventElements:    p.events,
		BehaviorElements: p.behaviors,
		Total:            p.total,
	}
}

package element

import (
	"backend/biz/entity/account"
	"backend/biz/entity/rule"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
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
	RuleID           int64   `gorm:"column:rule_id;type:bigint;primaryKey;autoIncrement:true" json:"rule_id"` // 规则ID
	RuleType         int64   `gorm:"column:rule_type;type:int;not null" json:"rule_type"`                     // 规则类型
	RuleDesc         *string `gorm:"column:rule_desc;type:text" json:"rule_desc"`                             // 规则描述
	AppID            int64   `gorm:"column:app_id;type:bigint;not null" json:"app_id"`
	RuleElementID    int64   `gorm:"column:rule_element_id;type:bigint;primaryKey;autoIncrement:true" json:"rule_element_id"` // 规则元素ID
	RuleElementValue string  `gorm:"column:rule_element_value;type:text;not null" json:"rule_element_value"`                  // 规则元素值
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

	return nil
}

func (p *PageElement) GetResp() *backend.ElementInPageResp {
	events := make([]*backend.EventRuleElement, 0)
	behaviors := make([]*backend.BehaviorRuleElement, 0)

	for _, v := range p.resMo {
		if v == nil || v.RuleID <= 0 {
			continue
		}

		if v.RuleType == int64(rule.EventRule) {
			res := &backend.EventRuleElement{
				RuleID:    v.RuleID,
				RuleType:  v.RuleType,
				RuleDesc:  "",
				ElementID: v.RuleElementID,
			}
			if v.RuleDesc != nil {
				res.RuleDesc = *v.RuleDesc
			}

			event := rule.ParseEventElement(v.RuleElementValue)
			if event != nil {
				res.EventType = event.EventType
				res.MouseClickType = event.MouseClickType
				res.MouseClickButton = event.MouseClickButton
				res.KeyClickType = event.KeyClickType
				res.KeyValue = event.KeyValue
				res.ComponentNamePrefix = event.ComponentNamePrefix
			}

			events = append(events, res)
		} else if v.RuleType == int64(rule.BehaviorRule) {
			res := &backend.BehaviorRuleElement{
				RuleID:    v.RuleID,
				RuleType:  v.RuleType,
				RuleDesc:  "",
				ElementID: v.RuleElementID,
			}
			if v.RuleDesc != nil {
				res.RuleDesc = *v.RuleDesc
			}

			behavior := rule.ParseBehaviorElement(v.RuleElementValue)
			if behavior != nil && len(behavior.EventRuleIds) > 0 {
				res.EventRuleIds = behavior.EventRuleIds
			}

			behaviors = append(behaviors, res)
		}
	}

	return &backend.ElementInPageResp{
		StatusCode:       microtype.SuccessErr.Code,
		StatusMsg:        microtype.SuccessErr.Msg,
		EventElements:    events,
		BehaviorElements: behaviors,
		Total:            p.total,
	}
}

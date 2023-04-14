package element

import (
	"backend/biz/entity/rule"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"gorm.io/gorm"
)

type AddElement struct {
	req backend.AddElementReq
}

func NewAddElement(req backend.AddElementReq) *AddElement {
	return &AddElement{
		req: req,
	}
}

func (a *AddElement) Load(ctx context.Context) error {
	if !a.check() {
		return microtype.ParamCheckFailed
	}

	queryMo, err := query.Rule.WithContext(ctx).Where(query.Rule.RuleID.Eq(a.req.RuleID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return microtype.ElementQueryFailed
	}

	switch queryMo.RuleType {
	case int64(rule.EventRule):
		eventRuleEle := &rule.EventRuleElement{
			EventType:           a.req.EventType,
			MouseClickType:      a.req.MouseClickType,
			MouseClickButton:    a.req.MouseClickButton,
			KeyClickType:        a.req.KeyClickType,
			KeyValue:            a.req.KeyValue,
			ComponentNamePrefix: a.req.ComponentNamePrefix,
		}

		createMo := &model.RuleElement{
			RuleElementID:    0,
			RuleElementValue: rule.GeneEventElement(eventRuleEle),
			RuleID:           queryMo.RuleID,
		}

		err := query.RuleElement.WithContext(ctx).Create(createMo)
		if err != nil {
			return microtype.ElementCreateFailed
		}
	case int64(rule.BehaviorRule):
		behaviorRuleEle := &rule.BehaviorRuleElement{
			EventRuleIds: a.req.EventRuleIds,
		}

		createMo := &model.RuleElement{
			RuleElementID:    0,
			RuleElementValue: rule.GeneBehaviorElement(behaviorRuleEle),
			RuleID:           queryMo.RuleID,
		}

		err := query.RuleElement.WithContext(ctx).Create(createMo)
		if err != nil {
			return microtype.ElementCreateFailed
		}
	default:

	}

	return nil
}

func (a *AddElement) GetResp() *backend.AddElementResp {
	return &backend.AddElementResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

func (a *AddElement) check() bool {
	if a.req.RuleID <= 0 {
		return false
	}

	return true
}

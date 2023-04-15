package main

import (
	"backend/cmd/dal"
	"backend/cmd/dal/query"
	"context"
	"testing"
)

func Test1(t *testing.T) {
	dal.Init()

	appId := 2

	ruleDO := query.Rule
	ruleMO := ruleDO.WithContext(context.Background())
	elementDO := query.RuleElement
	type Result struct {
		RuleID           int64  `gorm:"column:rule_id;type:bigint;primaryKey;autoIncrement:true" json:"rule_id"` // 规则ID
		RuleType         int64  `gorm:"column:rule_type;type:int;not null" json:"rule_type"`                     // 规则类型
		RuleDesc         string `gorm:"column:rule_desc;type:text" json:"rule_desc"`                             // 规则描述
		AppID            int64  `gorm:"column:app_id;type:bigint;not null" json:"app_id"`
		RuleElementID    int64  `gorm:"column:rule_element_id;type:bigint;primaryKey;autoIncrement:true" json:"rule_element_id"` // 规则元素ID
		RuleElementValue string `gorm:"column:rule_element_value;type:text;not null" json:"rule_element_value"`                  // 规则元素值
	}
	res := make([]*Result, 0)

	err := ruleMO.Select(ruleDO.ALL, elementDO.RuleElementID, elementDO.RuleElementValue).LeftJoin(elementDO, ruleDO.RuleID.EqCol(elementDO.RuleID)).Where(ruleDO.AppID.Eq(int64(appId))).Scan(&res)
	if err != nil {
		println(err.Error())
	}

}

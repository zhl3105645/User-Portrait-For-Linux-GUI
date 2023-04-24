package rule

import (
	"backend/biz/entity/account"
	"backend/biz/entity/rule"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/util"
	"backend/cmd/dal/query"
	"context"
)

type PageRuleData struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string

	//
	appId int64
	total int64
	data  []*backend.RuleData
}

func NewPageRuleData(accountId int64, pageNum int64, pageSize int64, search string) *PageRuleData {
	return &PageRuleData{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

type result struct {
	RecordID          int64   `gorm:"column:record_id;type:bigint;primaryKey;autoIncrement:true" json:"record_id"` // 使用记录ID
	UserID            int64   `gorm:"column:user_id;type:bigint;not null" json:"user_id"`                          // 用户ID
	BeginTime         int64   `gorm:"column:begin_time;type:bigint;not null" json:"begin_time"`                    // 开始时间
	UseTime           int64   `gorm:"column:use_time;type:bigint;not null" json:"use_time"`                        // 使用时长
	EventRuleValue    *string `gorm:"column:event_rule_value;type:text" json:"event_rule_value"`                   // 事件规则数据
	BehaviorRuleValue *string `gorm:"column:behavior_rule_value;type:text" json:"behavior_rule_value"`             // 行为规则数据
	UserName          string  `gorm:"column:user_name;type:varchar(256);not null" json:"user_name"`                // 用户名
}

func (p *PageRuleData) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User
	//userMO := userDO.WithContext(ctx)

	// 使用记录
	res := make([]*result, 0)
	offset := (p.pageNum - 1) * p.pageSize
	count, err := recordMO.Select(recordDO.RecordID, recordDO.UserID, recordDO.BeginTime, recordDO.UseTime, recordDO.EventRuleValue, recordDO.BehaviorRuleValue, userDO.UserName).
		Join(userDO, recordDO.UserID.EqCol(userDO.UserID)).
		Where(userDO.AppID.Eq(p.appId), userDO.UserName.Like("%"+p.search+"%")).
		ScanByPage(&res, int(offset), int(p.pageSize))
	if err != nil {
		return microtype.RuleDataQueryFailed
	}

	// 规则
	eventRules, behaviorRules, err := rule.GetRuleModels(ctx, p.appId)
	if err != nil {
		return microtype.RuleQueryFailed
	}
	ruleMap := make(map[int64]string)
	for _, r := range eventRules {
		if r == nil {
			continue
		}

		ruleMap[r.RuleID] = r.RuleDesc
	}
	for _, r := range behaviorRules {
		if r == nil {
			continue
		}

		ruleMap[r.RuleID] = r.RuleDesc
	}

	p.total = count
	p.data = make([]*backend.RuleData, 0, len(res))
	for _, r := range res {
		if r == nil {
			continue
		}

		d := &backend.RuleData{
			RecordID:         r.RecordID,
			UserID:           r.UserID,
			UserName:         r.UserName,
			BeginTime:        util.GeneTimeFromTimestampMs(r.BeginTime),
			UseTime:          util.GeneTimeDurationFromMs(r.UseTime),
			BehaviorRuleData: nil,
		}

		//if r.EventRuleValue != nil {
		//	eles := make([]*backend.RuleElement, 0)
		//	for _, ele := range rule.ParseRuleElements(*r.EventRuleValue, ruleMap) {
		//		if ele == nil || ele.RuleID <= 0 {
		//			continue
		//		}
		//		eles = append(eles, ele)
		//	}
		//	d.EventRuleData = &backend.EventRuleData{
		//		RuleElements: eles,
		//	}
		//}
		if r.BehaviorRuleValue != nil {
			eles := rule.ParseRuleElements(*r.BehaviorRuleValue, ruleMap)
			d.BehaviorRuleData = &backend.BehaviorRuleData{
				RuleElements: eles,
			}
			if len(eles) > 0 {
				durationMap := rule.GetBehaviorDuration(eles)
				d.BehaviorRuleData.BehaviorDuration = make(map[string]int64, 0)
				for id, duration := range durationMap {
					desc := ""
					if v, ok := ruleMap[id]; ok && v != "" {
						d.BehaviorRuleData.BehaviorDuration[desc] = duration
					}
				}
			}
		}

		p.data = append(p.data, d)
	}

	return nil
}

func (p *PageRuleData) GetResp() *backend.RuleDataInPageResp {
	return &backend.RuleDataInPageResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		RuleData:   p.data,
		Total:      p.total,
	}
}

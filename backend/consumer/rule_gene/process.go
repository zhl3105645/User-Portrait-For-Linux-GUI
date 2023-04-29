package rule_gene

import (
	"backend/biz/entity/event_data"
	"backend/biz/entity/rule"
	"backend/cmd/dal/model"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
	"strconv"
	"strings"
)

type RuleData struct {
	ID   int64 // 规则ID
	Time int64 // 发生时间
}

func process(events [][]string, eventRules []*rule.EventRuleModel, behaviorRules []*rule.BehaviorRuleModel, ruleDescMap map[int64]string) (*model.Record, map[int64]int64) {
	beginTimeMs := int64(0)
	useTimeMs := int64(0)
	beginTime, endTime, err := getAppUseTime(events)
	if err != nil {
		logger.Error("get app use time failed. err=", err.Error())
		return nil, nil
	}
	beginTimeMs = beginTime
	useTimeMs = endTime - beginTime

	eventRuleData := make([]*RuleData, 0, len(events))
	lastEventTimeMs, _ := strconv.ParseInt(events[1][1], 10, 64)
	for i := 2; i < len(events)-1; i++ {
		e := events[i]
		if len(e) < 11 {
			continue
		}

		eventTimeMs, err := strconv.ParseInt(e[event_data.EventTimeIndex], 10, 64)
		if err != nil {
			continue
		}

		// 事件规则数据
		if eventTimeMs-lastEventTimeMs > event_data.MaxNoOperateTimeS*1000 {
			eventRuleData = append(eventRuleData, &RuleData{
				ID:   rule.EventRuleStopOperate,
				Time: lastEventTimeMs,
			})
			eventRuleData = append(eventRuleData, &RuleData{
				ID:   rule.EventRuleBeginOperate,
				Time: eventTimeMs,
			})
		} else {
			if eventRuleId := getEventRuleID(e, eventRules); eventRuleId != 0 {
				eventRuleData = append(eventRuleData, &RuleData{
					ID:   eventRuleId,
					Time: eventTimeMs,
				})
			}
		}
		lastEventTimeMs = eventTimeMs
	}

	// 事件规则数据
	eventData := ""
	for _, ruleData := range eventRuleData {
		eventData = eventData + fmt.Sprintf("(%d,%d)", ruleData.ID, ruleData.Time)
	}

	// 行为规则数据
	behaviorRuleData := getBehaviorRuleIDs(eventRuleData, behaviorRules)
	behaviorData := ""
	for i := 0; i < len(behaviorRuleData); i++ {
		ruleData := behaviorRuleData[i]
		behaviorData = behaviorData + fmt.Sprintf("(%d,%d)", ruleData.ID, ruleData.Time)
	}

	eles := rule.ParseRuleElements(behaviorData, ruleDescMap)
	// 行为时长map
	durationMap := rule.GetBehaviorDuration(eles)

	return &model.Record{
		UserID:            0,
		BeginTime:         beginTimeMs,
		UseTime:           useTimeMs,
		EventRuleValue:    proto.String(eventData),
		BehaviorRuleValue: proto.String(behaviorData),
	}, durationMap
}

// 返回行为ID序列
func getBehaviorRuleIDs(eventRuleData []*RuleData, behaviorRules []*rule.BehaviorRuleModel) []*RuleData {
	if len(eventRuleData) == 0 || len(behaviorRules) == 0 {
		return nil
	}

	// 行为规则ID序列
	behaviorRuleIDs := make([][]int64, 0)
	eventRuleIdsIndex2BehaviorID := make(map[int]int64) // ID 索引 -> 规则 索引
	for _, behaviorRule := range behaviorRules {
		if behaviorRule == nil {
			continue
		}
		for _, ele := range behaviorRule.Elements {
			if ele == nil || len(ele.EventRuleIds) == 0 {
				continue
			}

			ids := ele.EventRuleIds

			behaviorRuleIDs = append(behaviorRuleIDs, ids)
			eventRuleIdsIndex2BehaviorID[len(behaviorRuleIDs)-1] = behaviorRule.RuleID
		}
	}

	// 思路：使用双指针思想，同时保存大序列和序列数组的当前指针，当序列数组指针到头即表示该序列命中
	// 指针
	eventIndex := 0
	// 指针数组
	ruleIndex := make([]int, 0, len(behaviorRuleIDs))
	for i := 0; i < len(behaviorRuleIDs); i++ {
		ruleIndex = append(ruleIndex, 0)
	}
	// 序列数组开头在大序列的位置
	ruleBeginIndex := make([]int, 0, len(behaviorRuleIDs))
	for i := 0; i < len(behaviorRuleIDs); i++ {
		ruleBeginIndex = append(ruleBeginIndex, 0)
	}

	// 命中的结果序列
	res := make([]*RuleData, 0)
	// 开始计算
	for eventIndex < len(eventRuleData) {
		for idx, ruleIDs := range behaviorRuleIDs {
			if eventRuleData[eventIndex].ID == ruleIDs[ruleIndex[idx]] {
				if ruleIndex[idx] == 0 {
					ruleBeginIndex[idx] = eventIndex
				}
				ruleIndex[idx]++
				// 存在一序列满足条件
				if ruleIndex[idx] == len(ruleIDs) {
					res = append(res, &RuleData{
						ID:   int64(eventRuleIdsIndex2BehaviorID[idx]),
						Time: eventRuleData[ruleBeginIndex[idx]].Time,
					})
					// 指针归零
					for i := 0; i < len(ruleIndex); i++ {
						ruleIndex[i] = 0
					}
					break
				}
			}
		}

		eventIndex++
	}

	if len(res) <= 2 {
		return res
	}

	// 合并中间相同的行为，保留头尾
	newRes := make([]*RuleData, 0, len(res))

	preData := res[0]
	newRes = append(newRes, preData)
	for i := 1; i < len(res)-1; i++ {
		curData := res[i]
		if curData.ID != preData.ID {
			preData = curData
			newRes = append(newRes, preData)
		}
	}

	newRes = append(newRes, res[len(res)-1])

	return newRes
}

func getEventRuleID(event []string, eventRules []*rule.EventRuleModel) int64 {
	if eventRules == nil {
		return 0
	}

	eventTyp := event[event_data.EventTypeIndex]
	mouseClickTyp := event[event_data.MouseClickTypeIndex]
	mouseClickButton := event[event_data.MouseClickButtonIndex]
	keyClickTyp := event[event_data.KeyClickTypeIndex]
	keyValue := event[event_data.KeyCodeIndex]
	comName := event[event_data.ComponentNameIndex]

	eventTypV := int64(0)
	mouseClickTypV := int64(0)
	mouseClickButtonV := int64(0)
	keyClickTypV := int64(0)

	if eventTyp != "" {
		v, err := strconv.ParseFloat(eventTyp, 64)
		if err == nil {
			eventTypV = int64(v)
		}
	}
	if mouseClickTyp != "" {
		v, err := strconv.ParseFloat(mouseClickTyp, 64)
		if err == nil {
			mouseClickTypV = int64(v)
		}
	}
	if mouseClickButton != "" {
		v, err := strconv.ParseFloat(mouseClickButton, 64)
		if err == nil {
			mouseClickButtonV = int64(v)
		}
	}
	if keyClickTyp != "" {
		v, err := strconv.ParseFloat(keyClickTyp, 64)
		if err == nil {
			keyClickTypV = int64(v)
		}
	}

	for _, eventRule := range eventRules {
		if eventRule == nil {
			continue
		}

		for _, ele := range eventRule.Elements {
			if ele == nil {
				continue
			}

			if ele.EventType == eventTypV &&
				ele.MouseClickType == mouseClickTypV &&
				ele.MouseClickButton == mouseClickButtonV &&
				ele.KeyClickType == keyClickTypV &&
				strings.HasPrefix(keyValue, ele.KeyValue) &&
				strings.HasPrefix(comName, ele.ComponentNamePrefix) {
				return eventRule.RuleID
			}
		}
	}

	return 0
}

func getAppUseTime(events [][]string) (beginTime int64, endTime int64, err error) {
	if len(events) < 3 {
		return 0, 0, fmt.Errorf("length < 3")
	}

	beginEvent := events[1]
	endEvent := events[len(events)-1]
	if len(beginEvent) < 2 || len(endEvent) < 2 || beginEvent[0] != string(event_data.AppStart) || endEvent[0] != string(event_data.AppQuit) {
		return 0, 0, fmt.Errorf("app start or stop data error")
	}

	beginTime, err1 := strconv.ParseInt(beginEvent[1], 10, 64)
	endTime, err2 := strconv.ParseInt(endEvent[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("time is not int")
	}

	return beginTime, endTime, nil
}

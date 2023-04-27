package rule

import (
	"backend/biz/entity/event_data"
	"backend/biz/model/backend"
	"fmt"
	"strconv"
	"strings"
)

type Type int

const (
	EventRule    Type = 1
	BehaviorRule Type = 2
)

// 特殊事件 停止操作 开始操作

const EventRuleStopOperate = -1
const EventRuleBeginOperate = -2

// 特殊行为 未操作

const BehaviorRuleNoOperate = -3
const BehaviorRuleNoOperateDesc = "未操作"

type EventRuleElement struct {
	EventType           int64
	MouseClickType      int64
	MouseClickButton    int64
	KeyClickType        int64
	KeyValue            string
	ComponentNamePrefix string
}

type BehaviorRuleElement struct {
	EventRuleIds []int64
}

func GeneEventElement(element *EventRuleElement) string {
	s := make([]string, 6)
	s[0] = strconv.FormatInt(element.EventType, 10)
	s[1] = strconv.FormatInt(element.MouseClickType, 10)
	s[2] = strconv.FormatInt(element.MouseClickButton, 10)
	s[3] = strconv.FormatInt(element.KeyClickType, 10)
	s[4] = element.KeyValue
	s[5] = element.ComponentNamePrefix

	return fmt.Sprintf("(%s)", strings.Join(s, "|"))
}

func ParseEventElement(value string) *EventRuleElement {
	if len(value) < 2 {
		return nil
	}
	eles := strings.Split(value[1:len(value)-1], "|")
	if len(eles) < 6 {
		return nil
	}

	res := make([]int64, 4)

	for i := 0; i < 4; i++ {
		if eles[i] == "" {
			continue
		}
		a, err := strconv.ParseFloat(eles[i], 64)
		if err != nil {
			return nil
		}

		res[i] = int64(a)
	}

	return &EventRuleElement{
		EventType:           res[0],
		MouseClickType:      res[1],
		MouseClickButton:    res[2],
		KeyClickType:        res[3],
		KeyValue:            eles[4],
		ComponentNamePrefix: eles[5],
	}
}

func GeneBehaviorElement(element *BehaviorRuleElement) string {
	if element == nil || len(element.EventRuleIds) == 0 {
		return ""
	}

	s := make([]string, 0, len(element.EventRuleIds))
	for _, id := range element.EventRuleIds {
		s = append(s, strconv.FormatInt(id, 10))
	}

	return fmt.Sprintf("(%s)", strings.Join(s, ","))
}

func ParseBehaviorElement(value string) *BehaviorRuleElement {
	if len(value) < 2 {
		return nil
	}
	eles := strings.Split(value[1:len(value)-1], ",")
	if len(eles) <= 0 {
		return nil
	}
	ids := make([]int64, 0, len(eles))
	for _, ele := range eles {
		id, err := strconv.ParseInt(ele, 10, 64)
		if err != nil {
			return nil
		}

		ids = append(ids, id)
	}

	return &BehaviorRuleElement{
		EventRuleIds: ids,
	}
}

func ParseRuleElements(value string, ruleMap map[int64]string) []*backend.RuleElement {
	if len(value) < 2 || len(ruleMap) == 0 {
		return nil
	}

	eles := strings.Split(value[1:len(value)-1], ")(")
	if len(eles) == 0 {
		return nil
	}

	res := make([]*backend.RuleElement, 0, len(eles))
	for _, ele := range eles {
		if ele == "" {
			continue
		}

		strs := strings.Split(ele, ",")
		if len(strs) != 2 {
			continue
		}

		ruleIdStr, timeStr := strs[0], strs[1]
		ruleId, err1 := strconv.ParseInt(ruleIdStr, 10, 64)
		time, err2 := strconv.ParseInt(timeStr, 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		ruleDesc := ""
		// 未操作
		if ruleId == BehaviorRuleNoOperate {
			ruleDesc = BehaviorRuleNoOperateDesc
		} else if v, ok := ruleMap[ruleId]; ok {
			ruleDesc = v
		}

		res = append(res, &backend.RuleElement{
			RuleID:    ruleId,
			RuleDesc:  ruleDesc,
			Timestamp: time,
		})
	}
	return res
}

func GetBehaviorDuration(elements []*backend.RuleElement) map[int64]int64 {
	res := make(map[int64]int64)
	for i := 1; i < len(elements); i++ {
		last := elements[i-1]
		cur := elements[i]

		time := cur.Timestamp - last.Timestamp

		if v, ok := res[last.RuleID]; ok {
			res[last.RuleID] = v + time
		} else {
			res[last.RuleID] = time
		}
	}

	return res
}

func GetEventCnt(elements []*backend.RuleElement) map[int64]int64 {
	res := make(map[int64]int64)
	for _, ele := range elements {
		if ele == nil {
			continue
		}
		if cnt, ok := res[ele.RuleID]; ok {
			res[ele.RuleID] = cnt + 1
		} else {
			res[ele.RuleID] = 1
		}
	}

	return res
}

// 判断是否符合规则
func MatchEvent(event []string, value []string) bool {
	if event_data.ComponentNameIndex >= len(event) {
		return false
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

	for _, v := range value {
		ele := ParseEventElement(v)
		if ele == nil {
			return false
		}

		if ele.EventType == eventTypV &&
			ele.MouseClickType == mouseClickTypV &&
			ele.MouseClickButton == mouseClickButtonV &&
			ele.KeyClickType == keyClickTypV &&
			strings.HasPrefix(keyValue, ele.KeyValue) &&
			strings.HasPrefix(comName, ele.ComponentNamePrefix) {
			return true
		}
	}

	return false
}

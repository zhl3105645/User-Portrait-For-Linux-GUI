package rule

import (
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
const BehaviorRuleNoOperate = -1

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

	for i := 0; i < 3; i++ {
		if eles[i] == "" {
			continue
		}

		a, err := strconv.ParseInt(eles[i], 10, 64)
		if err != nil {
			return nil
		}

		res[i] = a
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

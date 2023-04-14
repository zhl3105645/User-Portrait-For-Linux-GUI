package singleUse

import (
	"backend/biz/entity/event_data"
	"backend/impl/rule"
	"encoding/csv"
	"fmt"
	"gopkg.in/yaml.v2"
	"math"
	"os"
	"strconv"
	"strings"
)

type BasicBehavior struct {
	UseTimeMS int64 `yaml:"UseTimeMS"`

	MouseClickCnt int64 `yaml:"MouseClickCnt"`
	MouseMoveCnt  int64 `yaml:"MouseMoveCnt"`
	MouseWheelCnt int64 `yaml:"MouseWheelCnt"`
	KeyClickCnt   int64 `yaml:"KeyClickCnt"`
	ShortcutCnt   int64 `yaml:"ShortcutCnt"`

	MouseMoveDis  float64 `yaml:"MouseMoveDis"`
	KeyClickSpeed float64 `yaml:"KeyClickSpeed"` // 按键速度

	EventRuleData    string `yaml:"EventRuleData"`    // 事件规则数据
	BehaviorRuleData string `yaml:"BehaviorRuleData"` // 行为规则数据

	BehaviorTime map[int64]int64 // 各类行为时长
}

type AppComponent struct {
	Components []*QTComponent `yaml:"Components"`
}

type QTComponent struct {
	Name        string `yaml:"Name"`
	Description string `yaml:"Description"`
}

type RuleData struct {
	ID   int64 // 规则ID
	Time int64 // 发生时间
}

const StopOperate = -1
const BeginOperate = -2

func Process(filePath string, componentMap map[string]*QTComponent, eventRules []*rule.EventRule, behaviorRules []*rule.BehaviorRule) *BasicBehavior {
	// 1 获取数据
	file, err := os.Open(filePath)
	if err != nil {
		panic(any(err.Error()))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	events, err := reader.ReadAll()
	if err != nil {
		panic(any(err.Error()))
	}

	// 2 处理数据
	appUseTimeMs := int64(0)

	mouseClickCnt := int64(0)
	mouseMoveCnt := int64(0)
	mouseWheelCnt := int64(0)
	keyClickCnt := int64(0)
	shortcutCnt := int64(0)

	mouseMoveDis := 0.0
	keyContinueClickCnt := int64(0)  // 持续按键次数
	keyContinueClickTime := int64(0) // 持续按键时间

	eventRuleData := make([]*RuleData, 0, len(events))

	// 2.1 使用时长
	appUseTimeMs, err = getAppUseTime(events)
	if err != nil {
		println(err.Error())
		return nil
	}

	// 2.2 事件数据
	moving := false
	beginPos := ""
	lastEventTimeMs, _ := strconv.ParseInt(events[1][1], 10, 64)
	for i := 2; i < len(events)-1; i++ {
		e := events[i]
		if len(e) < 11 {
			continue
		}
		eventTimeMs, err := strconv.ParseInt(e[1], 10, 64)
		if err != nil {
			continue
		}
		switch event_data.EventType(e[0]) {
		case event_data.MouseClick:
			// 鼠标双击: 单击 双击 单击
			if event_data.MouseClickType(e[3]) == event_data.Two && i+1 < len(events)-1 {
				if len(events[i-1]) == 11 && event_data.EventType(events[i-1][0]) == event_data.MouseClick && event_data.MouseClickType(events[i-1][3]) == "1" &&
					len(events[i+1]) == 11 && event_data.EventType(events[i+1][0]) == event_data.MouseClick && event_data.MouseClickType(events[i+1][3]) == "1" {
					i = i + 1 // 跳过下一次
					continue
				}
			}
			mouseClickCnt++
		case event_data.MouseMove:
			typ := event_data.MouseMoveType(e[5])
			if typ == event_data.MoveBegin {
				moving = true
				beginPos = e[2]
			} else if typ == event_data.MoveEnd && moving {
				moving = false
				dis := getDistance(beginPos, e[2])
				if dis > 0 {
					mouseMoveDis += dis
					mouseMoveCnt++
				}
				beginPos = ""
			}
		case event_data.KeyClick:
			keyClickCnt++
			if event_data.EventType(events[i-1][0]) == event_data.KeyClick {
				nowTime, err1 := strconv.ParseInt(e[1], 10, 64)
				lastTime, err2 := strconv.ParseInt(events[i-1][1], 10, 64)
				if err1 == nil && err2 == nil {
					keyContinueClickCnt++
					keyContinueClickTime += nowTime - lastTime
				}
			}

		case event_data.MouseWheel:
			mouseWheelCnt++
		case event_data.Shortcut:
			shortcutCnt++
		default:
		}

		// 组件信息
		if event_data.EventType(e[0]) == event_data.MouseClick || event_data.EventType(e[0]) == event_data.MouseMove ||
			event_data.EventType(e[0]) == event_data.MouseWheel || event_data.EventType(e[0]) == event_data.KeyClick || event_data.EventType(e[0]) == event_data.Shortcut {
			if _, ok := componentMap[e[8]]; !ok {
				componentMap[e[8]] = &QTComponent{
					Name:        e[8],
					Description: e[10],
				}
			}
		}

		// 事件数据
		if eventTimeMs-lastEventTimeMs > event_data.MaxNoOperateTimeS*1000 {
			eventRuleData = append(eventRuleData, &RuleData{
				ID:   StopOperate,
				Time: lastEventTimeMs,
			})
			eventRuleData = append(eventRuleData, &RuleData{
				ID:   BeginOperate,
				Time: eventTimeMs,
			})
		} else {
			if eventRuleId := getEventRuleID(e, eventRules); eventRuleId != 0 {
				eventRuleData = append(eventRuleData, &RuleData{
					ID:   int64(eventRuleId),
					Time: eventTimeMs,
				})
			}
		}
		lastEventTimeMs = eventTimeMs
	}

	// 事件数据
	eventData := ""
	for _, ruleData := range eventRuleData {
		eventData = eventData + fmt.Sprintf("(%d,%d)", ruleData.ID, ruleData.Time)
	}

	// 行为数据
	behaviorRuleData := getBehaviorRuleIDs(eventRuleData, behaviorRules)
	behaviorData := ""
	behaviorTimeMap := make(map[int64]int64)
	for i := 0; i < len(behaviorRuleData); i++ {
		ruleData := behaviorRuleData[i]
		behaviorData = behaviorData + fmt.Sprintf("(%d,%d)", ruleData.ID, ruleData.Time)
		if i > 0 {
			last := behaviorRuleData[i-1]
			time := ruleData.Time - last.Time
			if v, ok := behaviorTimeMap[last.ID]; ok {
				behaviorTimeMap[last.ID] = v + time
			} else {
				behaviorTimeMap[last.ID] = time
			}
		}
	}

	// 写入文件
	basicBehaviorFile, err := os.Create("./impl/singleUse/basicBehavior.yaml")
	if err != nil {
		println(err.Error())
		return nil
	}
	defer basicBehaviorFile.Close()

	keySpeed := 0.0
	if keyContinueClickTime > 0 {
		keySpeed = float64(keyContinueClickCnt) / float64(keyContinueClickTime)
	}

	data := &BasicBehavior{
		UseTimeMS:        appUseTimeMs,
		MouseClickCnt:    mouseClickCnt,
		MouseMoveCnt:     mouseMoveCnt,
		MouseWheelCnt:    mouseWheelCnt,
		KeyClickCnt:      keyClickCnt,
		ShortcutCnt:      shortcutCnt,
		MouseMoveDis:     mouseMoveDis,
		KeyClickSpeed:    keySpeed,
		EventRuleData:    eventData,
		BehaviorRuleData: behaviorData,
		BehaviorTime:     behaviorTimeMap,
	}

	encoder := yaml.NewEncoder(basicBehaviorFile)

	if err := encoder.Encode(data); err != nil {
		println(err.Error())
	}

	return data
}

// 返回行为ID序列
func getBehaviorRuleIDs(eventRuleData []*RuleData, behaviorRules []*rule.BehaviorRule) []*RuleData {
	if len(eventRuleData) == 0 || len(behaviorRules) == 0 {
		return nil
	}

	// 行为规则ID序列
	behaviorRuleIDs := make([][]int64, 0)
	ruleID2BehaviorID := make(map[int]int) // ID 索引 -> 规则 索引
	for _, r := range behaviorRules {
		if r == nil {
			continue
		}
		for _, str := range r.Behaviors {
			ids := make([]int64, 0, len(r.Behaviors))
			idStrs := strings.Split(str[1:len(str)-1], ",")
			for _, idStr := range idStrs {
				if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					ids = append(ids, id)
				}
			}
			behaviorRuleIDs = append(behaviorRuleIDs, ids)
			ruleID2BehaviorID[len(behaviorRuleIDs)-1] = r.Id
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
						ID:   int64(ruleID2BehaviorID[idx]),
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

	return res
}

func getEventRuleID(event []string, eventRules []*rule.EventRule) int {
	if eventRules == nil {
		return 0
	}

	eventTyp := event[0]
	mouseClickTyp := event[3]
	mouseClickButton := event[4]
	keyClickTyp := event[6]
	keyValue := event[7]
	comName := event[8]

	for _, eventRule := range eventRules {
		if eventRule == nil {
			continue
		}

		for _, eventStr := range eventRule.Events {
			eventValues := strings.Split(eventStr[1:len(eventStr)-1], "|")
			if len(eventValues) < 6 {
				continue
			}

			if eventTyp == eventValues[0] && mouseClickTyp == eventValues[1] &&
				mouseClickButton == eventValues[2] && keyClickTyp == eventValues[3] &&
				keyValue == eventValues[4] && strings.HasPrefix(comName, eventValues[5]) {
				return eventRule.Id
			}
		}
	}

	return 0
}

func getDistance(beginPos, endPos string) float64 {
	beginXY := strings.Split(beginPos[1:len(beginPos)-1], ",")
	endXY := strings.Split(endPos[1:len(endPos)-1], ",")
	if len(beginXY) < 2 || len(endXY) < 2 {
		return 0
	}

	beginX, err1 := strconv.ParseInt(beginXY[0], 10, 64)
	beginY, err2 := strconv.ParseInt(beginXY[1], 10, 64)
	endX, err3 := strconv.ParseInt(endXY[0], 10, 64)
	endY, err4 := strconv.ParseInt(endXY[1], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return 0
	}

	return math.Sqrt(math.Pow(float64(beginX-endX), 2) + math.Pow(float64(beginY-endY), 2))
}

func getAppUseTime(events [][]string) (int64, error) {
	if len(events) < 3 {
		return 0, fmt.Errorf("length < 3")
	}

	beginEvent := events[1]
	endEvent := events[len(events)-1]
	if len(beginEvent) < 2 || len(endEvent) < 2 || beginEvent[0] != string(event_data.AppStart) || endEvent[0] != string(event_data.AppQuit) {
		return 0, fmt.Errorf("app start or stop data error")
	}

	beginTime, err1 := strconv.ParseInt(beginEvent[1], 10, 64)
	endTime, err2 := strconv.ParseInt(endEvent[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("time is not int")
	}

	return endTime - beginTime, nil
}

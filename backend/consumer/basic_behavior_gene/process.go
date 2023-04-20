package basic_behavior_gene

import (
	"backend/biz/entity/event_data"
	"backend/cmd/dal/model"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
	"math"
	"strconv"
	"strings"
)

func process(events [][]string) *model.Record {
	// 基础数据
	beginTimeMs := int64(0)
	appUseTimeMs := int64(0)

	mouseClickCnt := int64(0)
	mouseMoveCnt := int64(0)
	mouseWheelCnt := int64(0)
	keyClickCnt := int64(0)
	shortcutCnt := int64(0)

	mouseMoveDis := 0.0
	keyContinueClickCnt := int64(0)    // 持续按键次数
	keyContinueClickTimeMs := int64(0) // 持续按键时间

	// 1. 使用时长
	beginTime, endTime, err := getAppUseTime(events)
	if err != nil {
		logger.Error("get app use time failed. err=", err.Error())
		return nil
	}
	beginTimeMs = beginTime
	appUseTimeMs = endTime - beginTime

	// 2. 事件数据
	moving := false
	beginPos := ""
	for i := 2; i < len(events)-1; i++ {
		e := events[i]

		switch event_data.EventType(e[event_data.EventTypeIndex]) {
		case event_data.MouseClick:
			if len(e) <= event_data.MouseClickTypeIndex {
				continue
			}
			// 鼠标双击: 单击 双击 单击
			if event_data.MouseClickType(e[event_data.MouseClickTypeIndex]) == event_data.Two && i+1 < len(events)-1 {
				if len(events[i-1]) > event_data.MouseClickTypeIndex && event_data.EventType(events[i-1][event_data.EventTypeIndex]) == event_data.MouseClick && event_data.MouseClickType(events[i-1][event_data.MouseClickTypeIndex]) == event_data.One &&
					len(events[i+1]) > event_data.MouseClickTypeIndex && event_data.EventType(events[i+1][event_data.EventTypeIndex]) == event_data.MouseClick && event_data.MouseClickType(events[i+1][event_data.MouseClickTypeIndex]) == event_data.One {
					i = i + 1 // 跳过下一次
					continue
				}
			}
			mouseClickCnt++
		case event_data.MouseMove:
			if len(e) <= event_data.MouseMoveTypeIndex {
				continue
			}
			typ := event_data.MouseMoveType(e[event_data.MouseMoveTypeIndex])
			if typ == event_data.MoveBegin {
				moving = true
				beginPos = e[event_data.MousePosIndex]
			} else if typ == event_data.MoveEnd && moving {
				moving = false
				dis := getDistance(beginPos, e[event_data.MousePosIndex])
				if dis > 0 {
					mouseMoveDis += dis
					mouseMoveCnt++
				}
				beginPos = ""
			}
		case event_data.KeyClick:
			if len(e) <= event_data.EventTimeIndex {
				continue
			}
			keyClickCnt++
			if event_data.EventType(events[i-1][event_data.EventTypeIndex]) == event_data.KeyClick {
				nowTime, err1 := strconv.ParseInt(e[event_data.EventTimeIndex], 10, 64)
				lastTime, err2 := strconv.ParseInt(events[i-1][event_data.EventTimeIndex], 10, 64)
				if err1 == nil && err2 == nil {
					keyContinueClickCnt++
					keyContinueClickTimeMs += nowTime - lastTime
				}
			}

		case event_data.MouseWheel:
			mouseWheelCnt++
		case event_data.Shortcut:
			shortcutCnt++
		default:
		}
	}

	speed := 0.0
	if keyContinueClickTimeMs != 0 {
		speed = float64(keyContinueClickCnt) / float64(keyContinueClickTimeMs) * 1000 * 60
	}

	return &model.Record{
		UserID:        0,
		BeginTime:     beginTimeMs,
		UseTime:       appUseTimeMs,
		MouseClickCnt: proto.Int64(mouseClickCnt),
		MouseMoveCnt:  proto.Int64(mouseMoveCnt),
		MouseMoveDis:  proto.Float64(mouseMoveDis),
		MouseWheelCnt: proto.Int64(mouseWheelCnt),
		KeyClickCnt:   proto.Int64(keyClickCnt),
		KeyClickSpeed: proto.Float64(speed),
		ShortcutCnt:   proto.Int64(shortcutCnt),
	}
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

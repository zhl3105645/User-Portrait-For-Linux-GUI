package label_gene

import (
	"backend/biz/entity/event_data"
	"backend/biz/entity/rule"
	"backend/biz/util"
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"strconv"
	"strings"
)

func processEventCntLabel(ctx context.Context, appId int64, labelId int64) map[int64]string {
	res := make(map[int64]string)
	// 数据文件路径
	userEventPath := getUserEventPath(ctx, appId)
	if len(userEventPath) == 0 {
		return res
	}

	// 处理数据
	userIds := make([]int64, 0, len(userEventPath))
	plMap := make(map[int64]string)         // 编程语言 u_id -> program language
	codeSpeedMap := make(map[int64]float64) // 打字速度 u_id -> code speed
	shortcutCntMap := make(map[int64]int64) // 快捷键次数 u_id -> cnt
	gitCntMap := make(map[int64]int64)      // git 操作次数 u_id -> cnt
	for userId, paths := range userEventPath {
		userIds = append(userIds, userId)
		cCnt, cppCnt := int64(0), int64(0)
		keyClickCnt, keyClickDuration := int64(0), int64(0)
		shortcutCnt := int64(0)
		gitCnt := int64(0)
		for _, path := range paths {
			events, err := openFile(path)
			if err != nil {
				logger.Error("open file failed. err=", err.Error())
				continue
			}
			switch labelId {
			case ProgramLanguage:
				c, cpp := processProgramLanguage(events)
				cCnt = cppCnt + c
				cppCnt = cppCnt + cpp
			case CodeSpeed:
				cnt, duration := processCodeSpeed(events)
				keyClickCnt = keyClickCnt + cnt
				keyClickDuration = keyClickDuration + duration
			case ShortcutFre:
				shortcutCnt = shortcutCnt + processShortcutCnt(events)
			case GitFre:
				gitCnt = gitCnt + processGitCnt(events)
			default:
			}
		}
		switch labelId {
		case ProgramLanguage:
			if cCnt > cppCnt {
				plMap[userId] = "1"
			} else {
				plMap[userId] = "2"
			}
		case CodeSpeed:
			codeSpeedMap[userId] = float64(keyClickCnt) / float64(keyClickDuration)
		case ShortcutFre:
			shortcutCntMap[userId] = shortcutCnt
		case GitFre:
			if gitCnt > 0 {
				gitCntMap[userId] = gitCnt
			}
		default:

		}
	}

	// 结果
	switch labelId {
	case ProgramLanguage:
		res = plMap
	case CodeSpeed:
		gradeMap := util.GradeByPercent(codeSpeedMap, []float64{0.3, 0.7})
		for userId, grade := range gradeMap {
			res[userId] = fmt.Sprintf("%d", grade)
		}
	case ShortcutFre:
		gradeMap := util.GradeByPercent(util.ConvertIntMap2Float(shortcutCntMap), []float64{0.3, 0.7})
		for userId, grade := range gradeMap {
			res[userId] = fmt.Sprintf("%d", grade)
		}
	case GitFre:
		gradeMap := util.GradeByPercent(util.ConvertIntMap2Float(gitCntMap), []float64{0.3, 0.7})
		for _, userId := range userIds {
			if v, ok := gitCntMap[userId]; !ok || v == 0 {
				res[userId] = "1"
			} else if grade, ok := gradeMap[userId]; ok {
				res[userId] = fmt.Sprintf("%d", grade+1)
			}
		}
	default:
	}

	return res
}

func processGitCnt(events [][]string) int64 {
	cnt := int64(0)
	// git操作
	value := []string{
		"(3|1|1|||MainWindow.menubar.menuGit)",
		"(3|1|1|||MainWindow.<class_name=QMenu>.actionGit_Create_Repository)",
		"(3|1|1|||MainWindow.<class_name=QMenu,1>.actionGit)",
	}
	for _, event := range events {
		if rule.MatchEvent(event, value) {
			cnt++
		}
	}
	return cnt
}

func processShortcutCnt(events [][]string) int64 {
	cnt := int64(0)
	for _, event := range events {
		if event_data.EventTypeIndex >= len(events) {
			continue
		}
		if event[event_data.EventTypeIndex] == string(event_data.Shortcut) {
			cnt++
		}
	}

	return cnt
}

func processCodeSpeed(events [][]string) (keyClickCnt int64, duration int64) {
	keyClickCnt = 0
	duration = 0
	// 代码区键盘输入
	value := []string{
		"(5|0|0|1||MainWindow.centralwidget.EditorPanel.splitterEditorPanel.EditorTabs)",
		"(5|0|0|2||MainWindow.centralwidget.EditorPanel.splitterEditorPanel.EditorTabs)",
	}

	match := false // 上次事件是否为键盘输入
	lastTime := int64(0)
	for _, event := range events {
		timeStr := event[event_data.EventTimeIndex]
		timeStamp, err := strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			logger.Error("event time parse failed. err=", err.Error())
			continue
		}
		if rule.MatchEvent(event, value) {
			if match == true { // 计算次数 & 时间
				keyClickCnt++
				duration = duration + timeStamp - lastTime
			}
			match = true
		} else {
			match = false
		}

		lastTime = timeStamp
	}

	return keyClickCnt, duration
}

func processProgramLanguage(events [][]string) (cCnt int64, cppCnt int64) {
	cMap := make(map[string]bool)
	cppMap := make(map[string]bool)
	for _, event := range events {
		if event_data.ComponentExtraIndex > len(event)-1 {
			continue
		}

		extra := event[event_data.ComponentExtraIndex]
		if strings.HasSuffix(extra, ".c") {
			cMap[extra] = true
		} else if strings.HasSuffix(extra, ".cpp") {
			cppMap[extra] = true
		}
	}

	return int64(len(cMap)), int64(len(cppMap))
}

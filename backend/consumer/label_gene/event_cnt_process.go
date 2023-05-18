package label_gene

import (
	"backend/biz/entity/rule"
	"backend/biz/hadoop"
	"backend/biz/util"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
	"strings"
)

func processEventCntLabel(ctx context.Context, appId int64, labelId int64) map[int64]string {
	res := make(map[int64]string)

	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User
	// 查询 mysql 已有记录
	records, err := recordMO.Join(userDO, recordDO.UserID.EqCol(userDO.UserID)).
		Where(userDO.AppID.Eq(appId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query mysql failed. Err=", err.Error())
		return nil
	}
	userId2Records := make(map[int64][]*model.Record)
	for _, rec := range records {
		if len(userId2Records[rec.UserID]) == 0 {
			userId2Records[rec.UserID] = make([]*model.Record, 0)
		}
		userId2Records[rec.UserID] = append(userId2Records[rec.UserID], rec)
	}

	// 处理数据
	userIds := make([]int64, 0)
	plMap := make(map[int64]string)         // 编程语言 u_id -> program language
	codeSpeedMap := make(map[int64]float64) // 打字速度 u_id -> code speed
	shortcutCntMap := make(map[int64]int64) // 快捷键次数 u_id -> cnt
	gitCntMap := make(map[int64]int64)      // git 操作次数 u_id -> cnt
	for userId, recs := range userId2Records {
		userIds = append(userIds, userId)
		cCnt, cppCnt := int64(0), int64(0)
		keyClickCnt, keyClickDuration := int64(0), int64(0)
		shortcutCnt := int64(0)
		gitCnt := int64(0)
		for _, rec := range recs {
			events, err := hadoop.QueryEventsByRecordId(ctx, rec.RecordID)
			if err != nil {
				logger.Error(fmt.Sprintf("query hadoop failed. recordId=%d, err=%s", rec.RecordID, err.Error()))
				continue
			}
			logger.Info(fmt.Sprintf("reordId=%d, 记录长度=%d", rec.RecordID, len(events)))

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
			gitCntMap[userId] = gitCnt
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
		for userId, grade := range gradeMap {
			res[userId] = fmt.Sprintf("%d", grade)
		}
	default:
	}

	return res
}

func processGitCnt(events []*hadoop.Event) int64 {
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

func processShortcutCnt(events []*hadoop.Event) int64 {
	cnt := int64(0)
	for _, event := range events {
		if event.EventType == hadoop.Shortcut {
			cnt++
		}
	}

	return cnt
}

func processCodeSpeed(events []*hadoop.Event) (keyClickCnt int64, duration int64) {
	keyClickCnt = 0
	duration = 0
	// 代码区键盘输入
	value := []string{
		"(5|0|0|1||MainWindow.centralwidget.EditorPanel.splitterEditorPanel.EditorTabs)",
		"(5|0|0|2||MainWindow.centralwidget.EditorPanel.splitterEditorPanel.EditorTabs)",
	}

	match := false // 上次事件是否为键盘输入
	lastTime := int64(0)
	for i := 1; i < len(events); i++ {
		event := events[i]
		timeStamp := event.EventTime
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

func processProgramLanguage(events []*hadoop.Event) (cCnt int64, cppCnt int64) {
	cMap := make(map[string]bool)
	cppMap := make(map[string]bool)
	for _, event := range events {
		extra := event.ComponentExtra
		if strings.HasSuffix(extra, ".c") {
			cMap[extra] = true
		} else if strings.HasSuffix(extra, ".cpp") {
			cppMap[extra] = true
		}
	}

	return int64(len(cMap)), int64(len(cppMap))
}

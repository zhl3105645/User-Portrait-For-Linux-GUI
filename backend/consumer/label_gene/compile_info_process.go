package label_gene

import (
	"backend/biz/entity/event_data"
	"backend/biz/entity/rule"
	"backend/biz/util"
	"backend/consumer/common"
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"strings"
)

func ProcessCompileInfo(ctx context.Context, appId int64) map[int64]string {
	res := make(map[int64]string)
	// 数据文件路径
	userEventPath := common.GetUserEventPath(ctx, appId)
	if len(userEventPath) == 0 {
		return res
	}

	// 采集各个用户的报警和错误信息 - 每次编译或运行只记录一次相同的报警
	warningCntMap := make(map[int64]map[string]int64) // user_id -> warning -> cnt
	errorCntMap := make(map[int64]map[string]int64)   // user_id -> error -> cnt
	for userId, paths := range userEventPath {
		warningCnt := make(map[string]int64)
		errorCnt := make(map[string]int64)
		for _, path := range paths {
			events, err := common.OpenFile(path)
			if err != nil {
				logger.Error("open file failed. err=", err.Error())
				continue
			}

			warnMap, errMap := getCompileInfo(events)
			for msg, cnt := range warnMap {
				if v, ok := warningCnt[msg]; ok {
					warningCnt[msg] = v + cnt
				} else {
					warningCnt[msg] = cnt
				}
			}
			for msg, cnt := range errMap {
				if v, ok := errorCnt[msg]; ok {
					errorCnt[msg] = v + cnt
				} else {
					errorCnt[msg] = cnt
				}
			}
		}
		if len(warningCnt) > 0 {
			warningCntMap[userId] = warningCnt
		}
		if len(errorCnt) > 0 {
			errorCntMap[userId] = errorCnt
		}
	}

	// 写入文件，作为记录
	// common.WriteToDataToPath(warningCntMap, "D:\\graudation2\\code\\backend\\consumer\\label_gene\\data\\compile_warning.csv", []string{"user_id", "msg", "cnt"})
	// common.WriteToDataToPath(errorCntMap, "D:\\graudation2\\code\\backend\\consumer\\label_gene\\data\\compile_error.csv", []string{"user_id", "msg", "cnt"})

	// 错误权重累计
	warningWeight := int64(1)
	errorWeight := int64(2)
	scoreMap := make(map[int64]int64)
	for userId, desc2Cnt := range warningCntMap {
		cnt := int64(0)
		for _, c := range desc2Cnt {
			cnt += c
		}

		if score, ok := scoreMap[userId]; ok {
			scoreMap[userId] = score + cnt*warningWeight
		} else {
			scoreMap[userId] = cnt * warningWeight
		}
	}
	for userId, desc2Cnt := range errorCntMap {
		cnt := int64(0)
		for _, c := range desc2Cnt {
			cnt += c
		}

		if score, ok := scoreMap[userId]; ok {
			scoreMap[userId] = score + cnt*errorWeight
		} else {
			scoreMap[userId] = cnt * errorWeight
		}
	}

	// 权重分级
	gradeMap := util.GradeByPercent(util.ConvertIntMap2Float(scoreMap), []float64{0.3, 0.7})
	for userId, grade := range gradeMap {
		res[userId] = fmt.Sprintf("%d", 4-grade)
	}

	return res
}

func getCompileInfo(events [][]string) (warnMap map[string]int64, errMap map[string]int64) {
	warnMap = make(map[string]int64)
	errMap = make(map[string]int64)

	compileBtn := []string{
		"(3|1|1|||MainWindow.menubar.menuExecute.actionCompile)",
		"(3|1|1|||MainWindow.toolbarCompile.<class_name=QToolButton>)",
		"(7|||1|F9|)",
		"(3|1|1|||MainWindow.menubar.menuExecute.actionRebuild)",
		"(3|1|1||MainWindow.toolbarCompile.<class_name=QToolButton,2>)",
		"(7|||1|F12|)",
	}

	runBtn := []string{
		"(3|1|1|||MainWindow.menubar.menuExecute.actionRun)",
		"(3|1|1|||MainWindow.toolbarCompile.<class_name=QToolButton,1>)",
		"(3|1|1|||MainWindow.<class_name=QMenu>.actionRun)",
		"(7|||1|F10|)",
	}

	compileIssueTab := []string{
		"(3|1|1|||MainWindow.dockMessages.tabMessages.qt_tabwidget_stackedwidget.tabIssues.tableIssues)",
		"(3|1|2|||MainWindow.dockMessages.tabMessages.qt_tabwidget_stackedwidget.tabIssues.tableIssues)",
		"(3|2|1|||MainWindow.dockMessages.tabMessages.qt_tabwidget_stackedwidget.tabIssues.tableIssues)",
		"(3|2|2|||MainWindow.dockMessages.tabMessages.qt_tabwidget_stackedwidget.tabIssues.tableIssues)",
		"(4|||||MainWindow.dockMessages.tabMessages.qt_tabwidget_stackedwidget.tabIssues.tableIssues)",
	}

	existMap := make(map[string]bool)
	for _, event := range events {
		// 一次运行 or 编译 清空记录
		if rule.MatchEvent(event, compileBtn) || rule.MatchEvent(event, runBtn) {
			existMap = make(map[string]bool)
		}

		// 位于编译信息框
		if !rule.MatchEvent(event, compileIssueTab) {
			continue
		}

		// 报警内容
		if event_data.ComponentExtraIndex < len(event) {
			msg := event[event_data.ComponentExtraIndex]
			if _, ok := existMap[msg]; ok {
				continue
			}
			if strings.HasPrefix(msg, "[警告]") {
				if cnt, ok := warnMap[msg]; ok {
					warnMap[msg] = cnt + 1
				} else {
					warnMap[msg] = 1
				}
				existMap[msg] = true
			} else if strings.HasPrefix(msg, "[错误]") {
				if cnt, ok := errMap[msg]; ok {
					errMap[msg] = cnt + 1
				} else {
					errMap[msg] = 1
				}
				existMap[msg] = true
			}
		}
	}

	return warnMap, errMap
}

package label_gene

import (
	"backend/biz/entity/event_data"
	"backend/biz/util"
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"math"
	"strconv"
	"time"
)

func processUsePreLabel(ctx context.Context, appId int64, labelId int64) map[int64]string {
	res := make(map[int64]string)
	// 数据文件路径
	userEventPath := getUserEventPath(ctx, appId)
	if len(userEventPath) == 0 {
		return res
	}
	// 开始时间，结束时间
	beginTimeMap := make(map[int64][]int64) // u_id -> []time
	endTimeMap := make(map[int64][]int64)
	for userId, paths := range userEventPath {
		beginTimeMap[userId] = make([]int64, 0, len(paths))
		endTimeMap[userId] = make([]int64, 0, len(paths))
		for _, path := range paths {
			events, err := openFile(path)
			if err != nil {
				logger.Error("open file failed. err=", err.Error())
				continue
			}
			beginTime, endTime, err := getBeginAndEndTime(events)
			if err != nil {
				logger.Error("get app use time failed. err=", err.Error())
				continue
			}

			beginTimeMap[userId] = append(beginTimeMap[userId], beginTime)
			endTimeMap[userId] = append(endTimeMap[userId], endTime)
		}
	}

	// 处理得到使用时长、使用时间段、活跃度
	useTimeMap, usePeriodMap, useActivityMap := getUseFreLabels(beginTimeMap, endTimeMap)
	switch labelId {
	case UseTime:
		return useTimeMap
	case UsePeriod:
		return usePeriodMap
	case UseActivity:
		return useActivityMap
	}

	return nil
}

func getUseFreLabels(beginTimeMap map[int64][]int64, endTimeMap map[int64][]int64) (useTimeMap map[int64]string, usePeriodMap map[int64]string, useActivityMap map[int64]string) {
	useTimeMap = make(map[int64]string)
	usePeriodMap = make(map[int64]string)
	useActivityMap = make(map[int64]string)

	//
	activityMap := make(map[int64]float64)
	for userId, beginTimes := range beginTimeMap {
		endTimes := endTimeMap[userId]
		if len(endTimes) == 0 || len(beginTimes) == 0 || len(endTimes) != len(beginTimes) {
			logger.Warn("begin time and end time data is wrong.")
			continue
		}
		//
		dayUseTime := make(map[string]int64) // 天 -> 时间
		totalTime := int64(0)                // 总使用时间
		periodCnt := make(map[string]int64)  // 时间段 -> 次数
		activity := 0.0
		for idx := range beginTimes {
			beginTime := beginTimes[idx]
			endTime := endTimes[idx]
			// 使用时间
			useTime := endTime - beginTime
			date := time.UnixMilli(beginTime).Format("2006-01-02")
			if v, ok := dayUseTime[date]; ok {
				dayUseTime[date] = v + useTime
			} else {
				dayUseTime[date] = useTime
			}
			totalTime += useTime
			// 时间段
			beginHour := time.UnixMilli(beginTime).Hour()
			endHour := time.UnixMilli(endTime).Hour()
			period := fmt.Sprintf("%d~%d时", beginHour, endHour+1)
			if v, ok := periodCnt[period]; ok {
				periodCnt[period] = v + 1
			} else {
				periodCnt[period] = 1
			}
			// 活跃度
			activity = activity + activityExp(time.UnixMilli(beginTime), useTime)
		}
		// 天平均时间
		useTimeMap[userId] = util.GeneTimeDurationFromMs(totalTime / int64(len(dayUseTime)))
		// 最大时间段
		maxPeriod := ""
		maxCnt := int64(0)
		for period, cnt := range periodCnt {
			if cnt > maxCnt {
				maxCnt = cnt
				maxPeriod = period
			}
		}
		usePeriodMap[userId] = maxPeriod
		// 活跃度
		activityMap[userId] = activity
	}

	// 活跃度分级
	activityGradeMap := util.GradeByPercent(activityMap, []float64{0.3, 0.7})
	for userId, grade := range activityGradeMap {
		useActivityMap[userId] = fmt.Sprintf("%d", grade)
	}

	return
}

func getBeginAndEndTime(events [][]string) (beginTime int64, endTime int64, err error) {
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

// 活跃度计算公式
func activityExp(ti time.Time, useTime int64) float64 {
	now := time.Now()
	day := now.Sub(ti).Abs().Hours() / 24
	return math.Pow(math.E, -day) * float64(useTime)
}

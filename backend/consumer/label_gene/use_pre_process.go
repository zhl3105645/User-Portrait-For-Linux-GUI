package label_gene

import (
	"backend/biz/util"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
	"math"
	"time"
)

func processUsePreLabel(ctx context.Context, appId int64, labelId int64) map[int64]string {
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
	// 开始时间，结束时间
	beginTimeMap := make(map[int64][]int64) // u_id -> []time
	endTimeMap := make(map[int64][]int64)
	for _, rec := range records {
		uId := rec.UserID

		if len(beginTimeMap[uId]) == 0 {
			beginTimeMap[uId] = make([]int64, 0)
		}
		if len(endTimeMap[uId]) == 0 {
			endTimeMap[uId] = make([]int64, 0)
		}
		beginTimeMap[uId] = append(beginTimeMap[uId], rec.BeginTime)
		endTimeMap[uId] = append(endTimeMap[uId], rec.BeginTime+rec.UseTime)
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

// 活跃度计算公式
func activityExp(ti time.Time, useTime int64) float64 {
	now := time.Now()
	day := math.Abs(now.Sub(ti).Hours() / 24)
	return math.Pow(math.E, -day) * float64(useTime)
}

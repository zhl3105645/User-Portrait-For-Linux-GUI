package util

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

func GeneTimeFromTimestampMs(timestampMs int64) string {
	t := time.UnixMilli(timestampMs).Format("2006-01-02 15:04:05")
	return t
}

func GeneTimeDurationFromMs(timeMs int64) string {
	day := timeMs / 1000 / 60 / 60 / 24  // 天
	hour := timeMs / 1000 / 60 / 60 % 24 // 小时
	min := timeMs / 1000 / 60 % 60       // 分
	sec := timeMs / 1000 % 60            // 秒

	dayStr, hourStr, minStr, secStr := "", "", "", ""

	if day > 0 {
		dayStr = fmt.Sprintf("%d天", day)
	}
	if hour > 0 {
		hourStr = fmt.Sprintf("%d时", hour)
	}
	if min > 0 {
		minStr = fmt.Sprintf("%d分", min)
	}
	if sec > 0 {
		secStr = fmt.Sprintf("%d秒", sec)
	}

	return dayStr + hourStr + minStr + secStr
}

func Decimal(value float64, p int) float64 {
	value, _ = strconv.ParseFloat(strconv.FormatFloat(value, 'f', p, 64), 64)
	return value
}

// GradeByPercent 按百分比分级, 返回等级从1开始
func GradeByPercent(data map[int64]float64, percents []float64) map[int64]int64 {
	floatData := make([]float64, 0, len(data))
	for _, d := range data {
		floatData = append(floatData, d)
	}
	sort.Float64s(floatData)

	indexs := make([]int, 0, len(percents))
	for _, p := range percents {
		idx := int(math.Round(float64(len(data)) * p))
		if idx >= len(floatData) {
			idx = len(floatData) - 1
		}
		indexs = append(indexs, idx)
	}

	res := make(map[int64]int64, len(data))
	for userId, d := range data {
		grade := 1 //默认最低等级
		if d < floatData[indexs[0]] {
			grade = 1
		} else if d >= floatData[indexs[len(percents)-1]] {
			grade = len(percents) + 1
		} else {
			for g, idx := range indexs {
				if d >= floatData[idx] {
					grade = g + 2
				} else {
					break
				}
			}
		}

		res[userId] = int64(grade)
	}

	return res
}

func ConvertIntMap2Float(value map[int64]int64) map[int64]float64 {
	if value == nil {
		return nil
	}
	res := make(map[int64]float64, len(value))
	for k, v := range value {
		res[k] = float64(v)
	}

	return res
}

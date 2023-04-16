package util

import (
	"fmt"
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
		hourStr = fmt.Sprintf("%d小时", hour)
	}
	if min > 0 {
		minStr = fmt.Sprintf("%d分钟", min)
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

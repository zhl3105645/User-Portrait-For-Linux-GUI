package util

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
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

// GradeByPercent 按值域百分比分级, 返回等级从1开始
func GradeByPercent(data map[int64]float64, percents []float64) map[int64]int64 {
	if len(data) == 0 {
		return nil
	}

	floatData := make([]float64, 0, len(data))
	for _, d := range data {
		floatData = append(floatData, d)
	}
	sort.Float64s(floatData)

	min := floatData[0]
	max := floatData[len(floatData)-1]
	// n个等级，n+1个分界线
	grades := make([]float64, 0, len(percents)+1)
	grades = append(grades, min)
	for _, per := range percents {
		grades = append(grades, min+per*(max-min))
	}

	res := make(map[int64]int64, len(data))
	for userId, d := range data {
		grade := 0
		for grade < len(grades) {
			if d >= grades[grade] {
				grade++
			} else {
				break
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

// GetCharNumberOfString 获取字符串的权重
func GetCharNumberOfString(s string) int {
	chCnt := 0
	enStr := ""
	for _, v := range s {
		if unicode.Is(unicode.Han, v) {
			chCnt++
		} else {
			enStr = enStr + string(v)
		}
	}

	enCnt := 0
	ens := strings.Split(enStr, " ")
	for _, en := range ens {
		if en != "" {
			enCnt++
		}
	}

	return chCnt/2 + enCnt // 汉字权重为英文一半
}

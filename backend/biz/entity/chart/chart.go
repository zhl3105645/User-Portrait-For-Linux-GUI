package chart

import (
	"backend/biz/model/backend"
	"backend/biz/util"
	"backend/cmd/dal/model"
	"encoding/json"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/thoas/go-funk"
	"strconv"
)

type Biz int

const (
	CntFloat  Biz = 1 // eg. 次数平均数
	TimeFloat Biz = 2 // eg. 时长平均数
	Enum      Biz = 3 // eg. 聚类
	All       Biz = 4 // eg. 行为时长 求总的比例

	Label Biz = 5 // eg. 标签
)

func GetLabelOption(biz Biz, data []*model.LabelDatum, labelDescMap map[int64]string) *backend.ChartOption {
	if biz != Label {
		return nil
	}

	option := &backend.ChartOption{
		XAxis:   nil,
		YAxis:   nil,
		Tooltip: nil,
		Series:  nil,
	}

	cntMap := make(map[int64]int64) // data -> cnt
	for _, d := range data {
		if d == nil {
			continue
		}

		i, _ := strconv.ParseInt(d.Data, 10, 64)
		if cnt, ok := cntMap[i]; ok {
			cntMap[i] = cnt + 1
		} else {
			cntMap[i] = 1
		}
	}

	descMap := make(map[string]int64) // desc -> cnt
	for value, cnt := range cntMap {
		desc, ok := labelDescMap[value]
		if !ok || desc == "" {
			continue
		}

		descMap[desc] = cnt
	}

	option.Tooltip = &backend.ToolTip{
		Trigger:   "item",
		Formatter: "{c} ({d}%)",
	}
	option.Series = []*backend.Series{pie(descMap)}

	return option
}

func GetModelOption(biz Biz, data []*model.ModelDatum, ruleDescMap map[int64]string) *backend.ChartOption {
	option := &backend.ChartOption{
		XAxis:   nil,
		YAxis:   nil,
		Tooltip: nil,
		Series:  nil,
	}

	switch biz {
	case CntFloat:
		fallthrough
	case TimeFloat:
		ds := getFloat(data)
		if len(ds) == 0 {
			return nil
		}

		// series 柱状图 + 曲线图 x轴
		cnt, prop, xLabel := splitGroup(biz, ds, 20)
		series := []*backend.Series{
			&backend.Series{
				Type: "bar",
				Data: convertIntToString(cnt),
			},
			&backend.Series{
				Type:       "line",
				Smooth:     true,
				YAxisIndex: 1,
				Data:       convertFloatToString(prop),
			},
		}

		// y轴
		yAxis0 := &backend.Axis{
			Type:     "value",
			Data:     nil,
			Name:     "频次",
			Position: "left",
			AxisLabel: &backend.AxisLabel{
				Show: true,
			},
		}
		yAxis1 := &backend.Axis{
			Type:     "value",
			Data:     nil,
			Name:     "概率",
			Position: "right",
			AxisLabel: &backend.AxisLabel{
				Show: true,
			},
		}

		// x 轴
		xAxis := &backend.Axis{
			Type: "category",
			Data: xLabel,
			AxisLabel: &backend.AxisLabel{
				Rotate: -90,
				Show:   true,
			},
		}

		// 赋值
		option.Series = series
		option.YAxis = []*backend.Axis{yAxis0, yAxis1}
		option.XAxis = []*backend.Axis{xAxis}
		option.Tooltip = &backend.ToolTip{
			Trigger: "axis",
			AxisPointer: &backend.AxisPointer{
				Type: "cross",
			},
		}
		option.Toolbox = &backend.ToolBox{
			Feature: &backend.Feature{
				DataView:    &backend.View{Show: true},
				SaveAsImage: &backend.View{Show: true},
			},
		}
	case Enum:
		ds := getString(data)
		if len(ds) == 0 {
			return nil
		}
		m := make(map[string]int64)
		for _, d := range ds {
			if v, ok := m[d]; ok {
				m[d] = v + 1
			} else {
				m[d] = 1
			}
		}

		option.Tooltip = &backend.ToolTip{
			Trigger:   "item",
			Formatter: "{c} ({d}%)",
		}
		option.Series = []*backend.Series{pie(m)}
	case All:
		maps := getMap(data)
		allMap := make(map[int64]int64)
		for _, m := range maps {
			for k, v := range m {
				if x, ok := allMap[k]; ok {
					allMap[k] = x + v
				} else {
					allMap[k] = v
				}
			}
		}

		m := make(map[string]int64)
		for ruleId, cnt := range allMap {
			if v, ok := ruleDescMap[ruleId]; ok && v != "" {
				m[v] = cnt
			} else {
				m[fmt.Sprintf("规则%d", ruleId)] = cnt
			}
		}

		option.Tooltip = &backend.ToolTip{
			Trigger:   "item",
			Formatter: "{c} ({d}%)",
		}
		option.Series = []*backend.Series{pie(m)}
	}

	return option
}

func getMap(data []*model.ModelDatum) []map[int64]int64 {
	ds := make([]map[int64]int64, 0, len(data))
	for _, d := range data {
		if d == nil || d.Data == "" {
			continue
		}
		m := make(map[int64]int64)

		err := json.Unmarshal([]byte(d.Data), &m)
		if err != nil {
			logger.Error("json unmarshal failed. err=", err.Error())
			continue
		}
		for k, _ := range m {
			if k <= 0 {
				delete(m, k)
			}
		}
		if len(m) == 0 {
			continue
		}

		ds = append(ds, m)
	}

	return ds
}

func getFloat(data []*model.ModelDatum) []float64 {
	ds := make([]float64, 0, len(data))
	for _, d := range data {
		if d == nil {
			continue
		}
		v, _ := strconv.ParseFloat(d.Data, 64)
		if v == 0 {
			continue
		}
		ds = append(ds, v)
	}

	return ds
}

func getString(data []*model.ModelDatum) []string {
	ds := make([]string, 0, len(data))
	for _, d := range data {
		if d == nil {
			continue
		}
		if d.Data == "" {
			continue
		}
		ds = append(ds, d.Data)
	}

	return ds
}

// splitGroup 分为 group 组，返回 每组的数量，每组的占比
func splitGroup(biz Biz, data []float64, group int) ([]int, []float64, []string) {
	min := funk.MinFloat64(data)
	max := funk.MaxFloat64(data)
	if min == max {
		return []int{len(data)}, []float64{1}, []string{fmt.Sprintf("%.2f", min)}
	}

	// 每组数量
	groupCnt := make([]int, group)
	// 每组数值大小
	gap := (max - min) / float64(group)

	for _, d := range data {
		index := int((d - min) / gap)
		if index == group {
			index = group - 1
		}
		groupCnt[index] = groupCnt[index] + 1
	}
	// 占比
	prop := make([]float64, group)
	for idx, cnt := range groupCnt {
		prop[idx] = float64(cnt) / float64(len(data))
	}

	// x 轴，每三个转换一次
	xLabel := make([]string, group)
	for i := 0; 3*i < group; i++ {
		idx := 3 * i
		value := min + float64(idx)*gap
		xLabel[idx] = fmt.Sprintf("%.2f", value)
		if biz == TimeFloat {
			xLabel[idx] = util.GeneTimeDurationFromMs(int64(value))
		}
	}

	return groupCnt, prop, xLabel
}

func convertFloatToString(fs []float64) []string {
	res := make([]string, 0, len(fs))
	for _, f := range fs {
		res = append(res, fmt.Sprintf("%.2f", f))
	}

	return res
}

func convertIntToString(fs []int) []string {
	res := make([]string, 0, len(fs))
	for _, f := range fs {
		res = append(res, fmt.Sprintf("%d", f))
	}

	return res
}

func pie(m map[string]int64) *backend.Series {
	yData := make([]string, 0, len(m))
	for name, cnt := range m {
		s := struct {
			Value int64  `json:"value"`
			Name  string `json:"name"`
		}{
			Value: cnt,
			Name:  name,
		}

		str, _ := json.Marshal(s)
		yData = append(yData, string(str))
	}

	return &backend.Series{
		Type: "pie",
		Data: yData,
	}
}

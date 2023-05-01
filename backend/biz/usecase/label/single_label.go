package label

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/util"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"strconv"
)

type SingleLabel struct {
	labelId int64

	//
	labData []*model.LabelDatum

	//
	chartType int64
	pieLabel  *backend.PieLabel
	barLabel  *backend.BarLabel
}

func NewSingleLabel(labelId int64) *SingleLabel {
	return &SingleLabel{
		labelId: labelId,
	}
}

func (s *SingleLabel) Load(ctx context.Context) error {
	lab, err := query.Label.WithContext(ctx).Where(query.Label.LabelID.Eq(s.labelId)).First()
	if err != nil {
		return err
	}

	labData, err := query.LabelDatum.WithContext(ctx).Where(query.LabelDatum.LabelID.Eq(s.labelId)).Find()
	if err != nil {
		return err
	}
	s.labData = labData
	s.pieLabel = &backend.PieLabel{
		LabelName: lab.LabelName,
		Data:      make([]*backend.PieData, 0),
	}
	s.barLabel = &backend.BarLabel{
		XNames: make([]string, 0),
		Data:   make([]int64, 0),
	}

	if lab.DataType == CanEnum {
		s.chartType = Pie
		descMap := make(map[string]string)
		err := json.Unmarshal([]byte(*lab.LabelSemanticDesc), &descMap)
		if err != nil {
			return err
		}

		dataCntMap := make(map[string]int64)
		for _, d := range labData {
			if cnt, ok := dataCntMap[d.Data]; ok {
				dataCntMap[d.Data] = cnt + 1
			} else {
				dataCntMap[d.Data] = 1
			}
		}

		for d, cnt := range dataCntMap {
			if cnt <= 0 {
				continue
			}
			s.pieLabel.Data = append(s.pieLabel.Data, &backend.PieData{
				Name:  descMap[d],
				Value: cnt,
			})
		}
	} else if lab.FixType == Age {
		s.chartType = Pie
		descs := []string{"20岁以下", "20~23岁", "24~27岁", "27~30岁", "30岁以上"}
		cnts := make([]int64, 5)
		for _, d := range labData {
			age, err := strconv.ParseInt(d.Data, 10, 64)
			if err != nil {
				continue
			}
			if age < 20 {
				cnts[0] = cnts[0] + 1
			} else if age <= 23 {
				cnts[1] = cnts[1] + 1
			} else if age <= 27 {
				cnts[2] = cnts[2] + 1
			} else if age <= 30 {
				cnts[3] = cnts[3] + 1
			} else {
				cnts[4] = cnts[4] + 1
			}
		}
		for idx, cnt := range cnts {
			if cnt <= 0 {
				continue
			}
			s.pieLabel.Data = append(s.pieLabel.Data, &backend.PieData{
				Name:  descs[idx],
				Value: cnt,
			})
		}
	} else if lab.FixType == Career {
		s.chartType = Pie
		cntMap := make(map[string]int64)
		for _, d := range labData {
			if cnt, ok := cntMap[d.Data]; ok {
				cntMap[d.Data] = cnt + 1
			} else {
				cntMap[d.Data] = 1
			}
		}
		for desc, cnt := range cntMap {
			s.pieLabel.Data = append(s.pieLabel.Data, &backend.PieData{
				Name:  desc,
				Value: cnt,
			})
		}
	} else if lab.FixType == UseTime {
		s.chartType = Pie
		descs := []string{"0.5h以下", "0.5h~1h", "1h~1.5h", "1.5h~2h", "2h以上"}
		cnts := make([]int64, 5)
		for _, d := range labData {
			timeMS, err := util.ParseTimeDurationFromStr(d.Data)
			if err != nil {
				continue
			}
			hour := float64(timeMS) / float64(1000*60*60)
			if hour < 0.5 {
				cnts[0] = cnts[0] + 1
			} else if hour <= 1 {
				cnts[1] = cnts[1] + 1
			} else if hour <= 1.5 {
				cnts[2] = cnts[2] + 1
			} else if hour <= 2 {
				cnts[3] = cnts[3] + 1
			} else {
				cnts[4] = cnts[4] + 1
			}
		}
		for idx, cnt := range cnts {
			if cnt <= 0 {
				continue
			}
			s.pieLabel.Data = append(s.pieLabel.Data, &backend.PieData{
				Name:  descs[idx],
				Value: cnt,
			})
		}
	} else if lab.FixType == UsePeriod {
		s.chartType = Bar
		// 各小时分布
		hourCntMap := make(map[int64]int64)
		min, max := int64(24), int64(0)
		for _, d := range labData {
			beginHour, endHour, err := util.GetBeginHourAndEndHour(d.Data)
			if err != nil {
				continue
			}
			for i := beginHour; i < endHour; i++ {
				if cnt, ok := hourCntMap[i]; ok {
					hourCntMap[i] = cnt + 1
				} else {
					hourCntMap[i] = 1
				}
				if beginHour < min {
					min = beginHour
				}
				if endHour > max {
					max = endHour
				}
			}
		}
		for i := min; i <= max; i++ {
			s.barLabel.XNames = append(s.barLabel.XNames, strconv.FormatInt(i, 10))
			s.barLabel.Data = append(s.barLabel.Data, hourCntMap[i])
		}
	}

	return nil
}

func (s *SingleLabel) GetResp() *backend.SingleLabelResp {
	return &backend.SingleLabelResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		ChartType:  s.chartType,
		BarLabel:   s.barLabel,
		PieLabel:   s.pieLabel,
	}
}

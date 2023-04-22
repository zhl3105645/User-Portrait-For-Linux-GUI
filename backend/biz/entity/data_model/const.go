package data_model

import "backend/biz/entity/data_source"

type Type int

const (
	Statistics Type = 1
	Learning   Type = 2
)

type CalculateType int

// 数据计算规则 Statistics
const (
	Average      CalculateType = 1 // 平均数
	Mode         CalculateType = 2 // 众数
	RuleCnt      CalculateType = 3 // 各类事件规则次数平均数
	RuleDuration CalculateType = 4 // 各类行为规则时长平均数
	TopRule      CalculateType = 5 // Top 规则
)

type HttpType int

const (
	Post = 1 // 默认
)

type FrontDataType int

const (
	Continues FrontDataType = 1 // 连续
	FrontEnum FrontDataType = 2 // 枚举
)

type DataType int

// 模型存储数据类型
const (
	Float         DataType = 1  // 连续浮点数, eg 速度
	TimePeriod    DataType = 2  // 时间段
	TimeDuration  DataType = 3  // 时长
	Rule2Duration DataType = 4  // map, eg. rule_id -> duration
	Enum          DataType = 5  // 枚举值 eg 用户聚类
	Rule2Int      DataType = 6  // map, eg. rule_id -> cnt
	Int           DataType = 7  // 连续整数，eg 次数
	RuleId        DataType = 10 // 规则Id
)

func GetDataType(modelType int64, sourceType int64, sourceValue int64, calculateType int64, httpDataType int64) DataType {
	if modelType == int64(Statistics) {
		if sourceType == int64(data_source.Basic) {
			if sourceValue == int64(data_source.SourceUsePeriod) {
				return TimePeriod
			} else if sourceValue == int64(data_source.SourceUseTime) {
				return TimeDuration
			} else {
				return Float
			}
		} else if sourceType == int64(data_source.EventRule) {
			if calculateType == int64(Average) {
				return Float
			}
		} else if sourceType == int64(data_source.BehaviorRule) {
			if calculateType == int64(Average) {
				return TimeDuration
			}
		} else if sourceType == int64(data_source.AllBehaviorRule) {
			if calculateType == int64(RuleDuration) {
				return Rule2Duration
			} else if calculateType == int64(TopRule) {
				return RuleId
			}
		} else if sourceType == int64(data_source.AllEventRule) {
			if calculateType == int64(RuleCnt) {
				return Rule2Int
			} else if calculateType == int64(TopRule) {
				return RuleId
			}
		}
	} else if modelType == int64(Learning) {
		if httpDataType == int64(Continues) {
			return Float
		} else if httpDataType == int64(FrontEnum) {
			return Enum
		}
	}

	return Float
}

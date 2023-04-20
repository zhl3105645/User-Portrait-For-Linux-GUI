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
)

type Feature int

const (
	Label   Feature = 1 // 标签
	Predict Feature = 2 // 预测
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
	Float             DataType = 1  // 浮点数
	TimePeriod        DataType = 2  // 时间段
	TimeDuration      DataType = 3  // 时长
	MultiTimeDuration DataType = 4  // 行为时长 map
	Enum              DataType = 5  // 整数枚举值 eg 用户聚类
	TopRule           DataType = 10 // Top 规则
)

func GetDataType(modelType int64, sourceType int64, sourceValue int64, httpDataType int64) DataType {
	if modelType == int64(Statistics) {
		if sourceType == int64(data_source.Basic) && sourceValue == int64(data_source.SourceUsePeriod) {
			return TimePeriod
		} else if sourceType == int64(data_source.Basic) && sourceValue == int64(data_source.SourceUseTime) {
			return TimeDuration
		} else if sourceType == int64(data_source.AllBehaviorRule) {
			return MultiTimeDuration
		}

		return Float
	} else if modelType == int64(Learning) {
		if httpDataType == int64(Continues) {
			return Float
		} else if httpDataType == int64(FrontEnum) {
			return Enum
		}

		return Float
	}

	return Float
}

package data_source

import (
	"backend/biz/entity/rule"
	"backend/biz/microtype"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
)

type Type int

const (
	Basic           Type = 1 // 基础数据
	EventRule       Type = 2 // 具体事件规则
	BehaviorRule    Type = 3 // 具体行为规则
	AllEventRule    Type = 4 // 全部的事件规则
	AllBehaviorRule Type = 5 // 全部的行为规则
	Model           Type = 6 // 模型数据
)

var SourceTypeDesc = map[Type]string{
	Basic:           "基础数据",
	EventRule:       "事件规则",
	BehaviorRule:    "行为规则",
	AllEventRule:    "全部事件规则",
	AllBehaviorRule: "全部行为规则",
	Model:           "数据模型",
}

// Basic 的具体值
const (
	SourceMouseClickCnt = 1
	SourceMouseMoveCnt  = 2
	SourceMoveDis       = 3
	SourceMouseWheelCnt = 4
	SourceKeyClickCnt   = 5
	SourceKeyClickSpeed = 6
	SourceShortCut      = 7
	SourceUsePeriod     = 8 // 使用时间段 小时为单位
	SourceUseTime       = 9 // 使用时长
)

var BasicSourceDesc = map[int]string{
	SourceMouseClickCnt: "鼠标点击次数",
	SourceMouseMoveCnt:  "鼠标移动次数",
	SourceMoveDis:       "鼠标移动距离",
	SourceMouseWheelCnt: "鼠标滚轮次数",
	SourceKeyClickCnt:   "键盘点击次数",
	SourceKeyClickSpeed: "键盘点击速度",
	SourceShortCut:      "快捷键次数",
	SourceUsePeriod:     "应用使用时间段",
	SourceUseTime:       "应用使用时长",
}

var BasicSlice = []int{
	SourceMouseClickCnt,
	SourceMouseMoveCnt,
	SourceMoveDis,
	SourceMouseWheelCnt,
	SourceKeyClickCnt,
	SourceKeyClickSpeed,
	SourceShortCut,
	SourceUsePeriod,
	SourceUseTime,
}

func InitDataSource(ctx context.Context, appId int64) error {
	mos := make([]*model.DataSource, 0, len(BasicSlice))

	// 基础行为
	for _, id := range BasicSlice {
		mos = append(mos, &model.DataSource{
			SourceType:  int64(Basic),
			SourceValue: proto.Int64(int64(id)),
			AppID:       appId,
		})
	}

	// 全部事件规则
	mos = append(mos, &model.DataSource{
		SourceType: int64(AllEventRule),
		AppID:      appId,
	})

	// 全部行为规则
	mos = append(mos, &model.DataSource{
		SourceType: int64(AllBehaviorRule),
		AppID:      appId,
	})

	// 插入
	sourceDO := query.DataSource
	sourceMO := sourceDO.WithContext(ctx)
	err := sourceMO.Create(mos...)
	if err != nil {
		logger.Error("source create failed. err=", err.Error())
		return microtype.DataSourceCreateFailed
	}

	return nil
}

func AddRuleSource(ctx context.Context, ruleType int64, ruleID int64) error {
	mo := &model.DataSource{
		SourceType:  0,
		SourceValue: proto.Int64(ruleID),
	}

	if ruleType == int64(rule.EventRule) {
		mo.SourceType = int64(EventRule)
	} else if ruleType == int64(rule.BehaviorRule) {
		mo.SourceType = int64(BehaviorRule)
	}

	sourceDO := query.DataSource
	sourceMO := sourceDO.WithContext(ctx)
	err := sourceMO.Create(mo)
	if err != nil {
		logger.Error("source create failed. err=", err.Error())
		return microtype.DataSourceCreateFailed
	}

	return nil
}

func DeleteRuleSource(ctx context.Context, ruleID int64) error {
	sourceDO := query.DataSource
	sourceMO := sourceDO.WithContext(ctx)
	_, err := sourceMO.
		Where(sourceDO.SourceType.In(int64(EventRule), int64(BehaviorRule)),
			sourceDO.SourceValue.Eq(ruleID)).
		Delete()
	if err != nil {
		logger.Error("source delete failed. err=", err.Error())
		return microtype.DataSourceDeleteFailed
	}

	return nil
}

func AddModelSource(ctx context.Context, modelId, appId int64) error {
	mo := &model.DataSource{
		SourceType:  int64(Model),
		SourceValue: proto.Int64(modelId),
		AppID:       appId,
	}

	sourceDO := query.DataSource
	sourceMO := sourceDO.WithContext(ctx)
	err := sourceMO.Create(mo)
	if err != nil {
		logger.Error("source create failed. err=", err.Error())
		return microtype.DataSourceCreateFailed
	}

	return nil
}

func DeleteModelSource(ctx context.Context, modelId int64) error {
	sourceDO := query.DataSource
	sourceMO := sourceDO.WithContext(ctx)
	_, err := sourceMO.Where(sourceDO.SourceType.Eq(int64(Model)),
		sourceDO.SourceValue.Eq(modelId)).Delete()

	if err != nil {
		logger.Error("source delete failed. err=", err.Error())
		return microtype.DataSourceDeleteFailed
	}

	return nil
}

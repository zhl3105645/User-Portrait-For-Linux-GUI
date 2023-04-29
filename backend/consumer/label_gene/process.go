package label_gene

import (
	"context"
)

// 标签ID配置
const (
	// 基础信息
	Gender = 2 // 男 女
	Age    = 3
	Career = 4
	// 使用频率
	UseTime     = 6
	UsePeriod   = 7
	UseActivity = 18 // 不活跃 一般 活跃
	// 技能
	GitNorm     = 17 // 不规范 一般规范 规范
	CodeSpeed   = 13 // 慢 中 快
	CodeAbility = 19 // 弱 中 强
	// 使用偏好
	ProgramLanguage = 12 // C C++
	ShortcutFre     = 15 // 偶尔 适中 经常
	GitFre          = 16 // 偶尔 适中 经常
	BehaviorPrefer  = 20 // 行为偏好
)

func process(ctx context.Context, appId int64, labelId int64) map[int64]string {
	if labelId == Gender || labelId == Age || labelId == Career {
		return processUserLabel(ctx, appId, labelId)
	}

	if labelId == UseTime || labelId == UsePeriod || labelId == UseActivity {
		return processUsePreLabel(ctx, appId, labelId)
	}

	if labelId == ProgramLanguage || labelId == CodeSpeed || labelId == ShortcutFre || labelId == GitFre {
		return processEventCntLabel(ctx, appId, labelId)
	}

	if labelId == CodeAbility {
		return ProcessCompileInfo(ctx, appId)
	}

	if labelId == GitNorm {
		return ProcessGitDesc(ctx, appId)
	}

	if labelId == BehaviorPrefer {
		return processBehaviorPrefer(ctx)
	}

	return nil
}

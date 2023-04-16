package basic_behavior

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/util"
	"backend/cmd/dal/query"
	"context"
)

type PageBasicBehavior struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string

	//
	appId int64
	res   []*result
	total int64
}

func NewPageBasicBehavior(accountId int64, pageNum int64, pageSize int64, search string) *PageBasicBehavior {
	return &PageBasicBehavior{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

type result struct {
	RecordID          int64    `gorm:"column:record_id;type:bigint;primaryKey;autoIncrement:true" json:"record_id"` // 使用记录ID
	UserID            int64    `gorm:"column:user_id;type:bigint;not null" json:"user_id"`                          // 用户ID
	BeginTime         int64    `gorm:"column:begin_time;type:bigint;not null" json:"begin_time"`                    // 开始时间
	UseTime           *int64   `gorm:"column:use_time;type:bigint" json:"use_time"`                                 // 使用时长
	MouseClickCnt     *int64   `gorm:"column:mouse_click_cnt;type:bigint" json:"mouse_click_cnt"`                   // 鼠标点击次数
	MouseMoveCnt      *int64   `gorm:"column:mouse_move_cnt;type:bigint" json:"mouse_move_cnt"`                     // 鼠标移动次数
	MouseMoveDis      *float64 `gorm:"column:mouse_move_dis;type:double" json:"mouse_move_dis"`                     // 鼠标移动距离
	MouseWheelCnt     *int64   `gorm:"column:mouse_wheel_cnt;type:bigint" json:"mouse_wheel_cnt"`                   // 鼠标滚轮次数
	KeyClickCnt       *int64   `gorm:"column:key_click_cnt;type:bigint" json:"key_click_cnt"`                       // 键盘点击次数
	KeyClickSpeed     *float64 `gorm:"column:key_click_speed;type:double" json:"key_click_speed"`                   // 键盘点击速度 字符/min
	ShortcutCnt       *int64   `gorm:"column:shortcut_cnt;type:bigint" json:"shortcut_cnt"`                         // 快捷键次数
	EventRuleValue    *string  `gorm:"column:event_rule_value;type:text" json:"event_rule_value"`                   // 事件规则数据
	BehaviorRuleValue *string  `gorm:"column:behavior_rule_value;type:text" json:"behavior_rule_value"`             // 行为规则数据
	UserName          string   `gorm:"column:user_name;type:varchar(256);not null" json:"user_name"`                // 用户名
}

func (p *PageBasicBehavior) Load(ctx context.Context) error {
	ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}

	p.appId = ac.GetQueryAccount().AppID

	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User
	//userMO := userDO.WithContext(ctx)

	res := make([]*result, 0)
	offset := (p.pageNum - 1) * p.pageSize
	count, err := recordMO.Select(recordDO.ALL, userDO.UserName).
		Join(userDO, recordDO.UserID.EqCol(userDO.UserID)).
		Where(userDO.AppID.Eq(p.appId), userDO.UserName.Like("%"+p.search+"%")).
		ScanByPage(&res, int(offset), int(p.pageSize))
	if err != nil {
		return microtype.BasicBehaviorQueryFailed
	}

	p.total = count
	p.res = res

	return nil
}

func (p *PageBasicBehavior) GetResp() *backend.BasicBehaviorInPageResp {
	resp := &backend.BasicBehaviorInPageResp{
		StatusCode:     microtype.SuccessErr.Code,
		StatusMsg:      microtype.SuccessErr.Msg,
		BasicBehaviors: nil,
		Total:          p.total,
	}

	behaviors := make([]*backend.BasicBehavior, 0, len(p.res))
	for _, v := range p.res {
		if v == nil {
			continue
		}

		r := &backend.BasicBehavior{
			RecordID:      v.RecordID,
			UserID:        v.UserID,
			UserName:      v.UserName,
			BeginTime:     util.GeneTimeFromTimestampMs(v.BeginTime),
			UseTime:       "",
			MouseClickCnt: 0,
			MouseMoveCnt:  0,
			MouseMoveDis:  0,
			MouseWheelCnt: 0,
			KeyClickCnt:   0,
			KeyClickSpeed: 0,
			ShortcutCnt:   0,
		}

		if v.UseTime != nil {
			r.UseTime = util.GeneTimeDurationFromMs(*v.UseTime)
		}
		if v.MouseClickCnt != nil {
			r.MouseClickCnt = *v.MouseClickCnt
		}
		if v.MouseMoveCnt != nil {
			r.MouseMoveCnt = *v.MouseMoveCnt
		}
		if v.MouseMoveDis != nil {
			r.MouseMoveDis = util.Decimal(*v.MouseMoveDis, 2)
		}
		if v.MouseWheelCnt != nil {
			r.MouseWheelCnt = *v.MouseWheelCnt
		}
		if v.KeyClickCnt != nil {
			r.KeyClickCnt = *v.KeyClickCnt
		}
		if v.KeyClickSpeed != nil {
			r.KeyClickSpeed = util.Decimal(*v.KeyClickSpeed, 2)
		}
		if v.ShortcutCnt != nil {
			r.ShortcutCnt = *v.ShortcutCnt
		}

		behaviors = append(behaviors, r)
	}

	resp.BasicBehaviors = behaviors

	return resp
}

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameRecord = "record"

// Record mapped from table <record>
type Record struct {
	RecordID          int64    `gorm:"column:record_id;type:bigint;primaryKey;autoIncrement:true" json:"record_id"` // 使用记录ID
	UserID            int64    `gorm:"column:user_id;type:bigint;not null" json:"user_id"`                          // 用户ID
	BeginTime         int64    `gorm:"column:begin_time;type:bigint;not null" json:"begin_time"`                    // 开始时间
	UseTime           int64    `gorm:"column:use_time;type:bigint;not null" json:"use_time"`                        // 使用时长
	MouseClickCnt     *int64   `gorm:"column:mouse_click_cnt;type:bigint" json:"mouse_click_cnt"`                   // 鼠标点击次数
	MouseMoveCnt      *int64   `gorm:"column:mouse_move_cnt;type:bigint" json:"mouse_move_cnt"`                     // 鼠标移动次数
	MouseMoveDis      *float64 `gorm:"column:mouse_move_dis;type:double" json:"mouse_move_dis"`                     // 鼠标移动距离
	MouseWheelCnt     *int64   `gorm:"column:mouse_wheel_cnt;type:bigint" json:"mouse_wheel_cnt"`                   // 鼠标滚轮次数
	KeyClickCnt       *int64   `gorm:"column:key_click_cnt;type:bigint" json:"key_click_cnt"`                       // 键盘点击次数
	KeyClickSpeed     *float64 `gorm:"column:key_click_speed;type:double" json:"key_click_speed"`                   // 键盘点击速度 字符/min
	ShortcutCnt       *int64   `gorm:"column:shortcut_cnt;type:bigint" json:"shortcut_cnt"`                         // 快捷键次数
	EventRuleValue    *string  `gorm:"column:event_rule_value;type:text" json:"event_rule_value"`                   // 事件规则数据
	BehaviorRuleValue *string  `gorm:"column:behavior_rule_value;type:text" json:"behavior_rule_value"`             // 行为规则数据
}

// TableName Record's table name
func (*Record) TableName() string {
	return TableNameRecord
}

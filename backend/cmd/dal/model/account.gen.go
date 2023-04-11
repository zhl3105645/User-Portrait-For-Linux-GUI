// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameAccount = "account"

// Account mapped from table <account>
type Account struct {
	AccountID         int64  `gorm:"column:account_id;type:bigint;primaryKey;autoIncrement:true" json:"account_id"`   // 账号ID
	AccountName       string `gorm:"column:account_name;type:varchar(256);not null" json:"account_name"`              // 账号名
	AccountPwd        string `gorm:"column:account_pwd;type:varchar(256);not null" json:"account_pwd"`                // 账号密码
	AccountPermission int64  `gorm:"column:account_permission;type:int;not null;default:1" json:"account_permission"` // 账号权限
	AppID             int64  `gorm:"column:app_id;type:bigint;not null" json:"app_id"`                                // 应用名
}

// TableName Account's table name
func (*Account) TableName() string {
	return TableNameAccount
}

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTest = "test"

// Test mapped from table <test>
type Test struct {
	ID   *int64  `gorm:"column:id;type:int" json:"id"`
	Name *string `gorm:"column:name;type:varchar(20)" json:"name"`
}

// TableName Test's table name
func (*Test) TableName() string {
	return TableNameTest
}

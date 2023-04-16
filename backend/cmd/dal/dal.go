package dal

import (
	"backend/cmd/dal/query"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsn = "root:20010721zhl@tcp(localhost:3306)/profile?charset=utf8&parseTime=True&loc=Local"

var DB *gorm.DB

func Init() {
	DB = ConnectDB().Debug()
	query.SetDefault(DB)
}

func ConnectDB() (conn *gorm.DB) {
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Info),
		CreateBatchSize: 1000, // 批量插入大小 1000
	})
	if err != nil {
		panic(any(fmt.Errorf("cannot establish db connection: %w", err)))
	}
	return conn
}

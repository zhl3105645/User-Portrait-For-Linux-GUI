package file_input

import (
	"backend/biz/hadoop"
	"backend/cmd/dal/query"
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"time"
)

func Gene(appId int64, userId int64, extra string) {
	ctx := context.Background()

	extraStruct := struct {
		RecordIds []int64
		Events    []string
	}{}

	err := json.Unmarshal([]byte(extra), &extraStruct)
	if err != nil {
		logger.Error("json unmarshal failed. err=", err.Error())
		return
	}
	if len(extraStruct.Events) == 0 {
		return
	}

	// 写入hadoop
	ok := hadoop.WriteEvents(ctx, extraStruct.Events)
	// hadoop 写入失败, mysql 回滚
	if !ok {
		logger.Error("hadoop write failed.")
		_, _ = query.Record.WithContext(ctx).Where(query.Record.RecordID.In(extraStruct.RecordIds...)).Delete()
	}

	logger.Info("file upload end time = ", time.Now().UnixMilli())

	return
}

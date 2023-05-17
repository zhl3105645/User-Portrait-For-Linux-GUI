package upload

import (
	"backend/biz/entity/event_data"
	"backend/biz/hadoop"
	"backend/biz/microtype"
	"backend/cmd/dal"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/csv"
	"github.com/bytedance/gopkg/util/logger"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const appId = 2
const dirPath = "D:\\hadoop_data\\event_data\\2"

func TestName(t *testing.T) {
	ctx := context.Background()
	dal.Init()
	hadoop.Init(ctx)

	files, err := openDir()
	if err != nil {
		return
	}

	writeEvents := make([]string, 0)
	recordIds := make([]int64, 0)

	for _, file := range files {
		if file == nil {
			continue
		}

		fileName := file.Name()
		names := strings.Split(fileName, ".")
		if len(names) < 2 || names[1] != "csv" {
			logger.Error("file name wrong 1. err=", err.Error())
			continue
		}
		nums := strings.Split(names[0], "_")
		if len(nums) < 2 {
			logger.Error("file name wrong 2. err=", err.Error())
			continue
		}

		uIdStr := nums[0]
		uId, err := strconv.ParseInt(uIdStr, 10, 64)
		if err != nil {
			logger.Error("uid parse failed. err=", err.Error())
			continue
		}

		// 读文件
		events, err := readFile1(file)
		if err != nil {
			logger.Error("read file failed. err=", err.Error())
			continue
		}

		// 检查文件格式
		beginTime, useTime, ok := checkEvents(events)
		if !ok {
			logger.Error("file check failed")
			continue
		}

		// 插入使用记录
		mo := &model.Record{
			RecordID:  0,
			UserID:    uId,
			BeginTime: beginTime,
			UseTime:   useTime,
		}
		err = query.Record.WithContext(ctx).Create(mo)
		if err != nil {
			logger.Error("create record failed. err=", err.Error())
			continue
		}
		recordIds = append(recordIds, mo.RecordID)

		// 转换数据模型
		day := time.UnixMilli(beginTime).Format("2006-01-02")
		dbEvents := transEvent(mo.RecordID, day, events)
		writeEvents = append(writeEvents, dbEvents...)
	}

	logger.Info("file data extract end, time = ", time.Now().UnixMilli())
	// 写入hadoop
	ok := hadoop.WriteEvents(ctx, writeEvents)
	// hadoop 写入失败, mysql 回滚
	if !ok {
		logger.Error("hadoop write failed.")
		_, _ = query.Record.WithContext(ctx).Where(query.Record.RecordID.In(recordIds...)).Delete()
	}

	logger.Info("file upload end time = ", time.Now().UnixMilli())

	return
}

func transEvent(recordId int64, day string, events [][]string) []string {
	res := make([]string, 0, len(events))
	for _, event := range events {
		r := &hadoop.Event{
			RecordId: recordId,
			AppId:    int32(appId),
			Day:      day,
		}
		if event_data.EventTypeIndex < len(event) {
			eventType := event[event_data.EventTypeIndex]
			eventTypeVal, err := strconv.ParseInt(eventType, 10, 64)
			if err != nil {
				continue
			}
			r.EventType = int32(eventTypeVal)
		}
		if event_data.EventTimeIndex < len(event) {
			eventTime := event[event_data.EventTimeIndex]
			eventTimeVal, err := strconv.ParseInt(eventTime, 10, 64)
			if err != nil {
				continue
			}
			r.EventTime = eventTimeVal
		}
		if event_data.MousePosIndex < len(event) {
			r.MousePos = event[event_data.MousePosIndex]
		}
		if event_data.MouseClickTypeIndex < len(event) {
			mouseClickType := event[event_data.MouseClickTypeIndex]
			if mouseClickType == "" {
				r.MouseClickType = int32(0)
			} else {
				mouseClickTypeVal, err := strconv.ParseInt(mouseClickType, 10, 64)
				if err != nil {
					continue
				}
				r.MouseClickType = int32(mouseClickTypeVal)
			}
		}
		if event_data.MouseClickButtonIndex < len(event) {
			mouseClickBtn := event[event_data.MouseClickButtonIndex]
			if mouseClickBtn == "" {
				r.MouseClickBtn = int32(0)
			} else {
				mouseClickBtnVal, err := strconv.ParseInt(mouseClickBtn, 10, 64)
				if err != nil {
					continue
				}
				r.MouseClickBtn = int32(mouseClickBtnVal)
			}
		}
		if event_data.MouseMoveTypeIndex < len(event) {
			mouseMoveType := event[event_data.MouseMoveTypeIndex]
			if mouseMoveType == "" {
				r.MouseMoveType = int32(0)
			} else {
				mouseMoveTypeVal, err := strconv.ParseInt(mouseMoveType, 10, 64)
				if err != nil {
					continue
				}
				r.MouseMoveType = int32(mouseMoveTypeVal)
			}
		}
		if event_data.KeyClickTypeIndex < len(event) {
			keyClickType := event[event_data.KeyClickTypeIndex]
			if keyClickType == "" {
				r.KeyClickType = int32(0)
			} else {
				keyClickTypeVal, err := strconv.ParseInt(keyClickType, 10, 64)
				if err != nil {
					continue
				}
				r.KeyClickType = int32(keyClickTypeVal)
			}
		}
		if event_data.KeyCodeIndex < len(event) {
			r.KeyCode = event[event_data.KeyCodeIndex]
		}
		if event_data.ComponentNameIndex < len(event) {
			r.ComponentName = event[event_data.ComponentNameIndex]
		}
		if event_data.ComponentTypeIndex < len(event) {
			componentType := event[event_data.ComponentTypeIndex]
			if componentType == "" {
				r.ComponentType = int32(0)
			} else {
				componentTypeVal, err := strconv.ParseInt(componentType, 10, 64)
				if err != nil {
					continue
				}
				r.ComponentType = int32(componentTypeVal)
			}
		}
		if event_data.ComponentExtraIndex < len(event) {
			r.ComponentExtra = event[event_data.ComponentExtraIndex]
		}

		res = append(res, r.String())
	}

	return res
}

func openDir() ([]os.FileInfo, error) {
	dir, err := os.Open(dirPath)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, microtype.DirOpenFailed
	}
	defer dir.Close()

	infos, err := dir.Readdir(-1)
	if err != nil {
		return nil, microtype.DirReadFailed
	}

	return infos, nil
}

func readFile1(file os.FileInfo) ([][]string, error) {
	fileName := file.Name()
	path := dirPath + "\\" + fileName
	filePtr, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer filePtr.Close()

	reader := csv.NewReader(filePtr)

	reader.FieldsPerRecord = -1
	events, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return events, nil
}

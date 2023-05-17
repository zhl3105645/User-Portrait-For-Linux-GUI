package upload

import (
	"backend/biz/entity/event_data"
	"backend/biz/hadoop"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"encoding/csv"
	"github.com/bytedance/gopkg/util/logger"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
	"time"
)

type FileUpload struct {
	userId int64
	files  []*multipart.FileHeader

	appId int64
}

func NewFileUpload(userId int64, files []*multipart.FileHeader) *FileUpload {
	return &FileUpload{
		userId: userId,
		files:  files,
	}
}

func (f *FileUpload) Load(ctx context.Context) error {
	userMo, err := query.User.WithContext(ctx).Where(query.User.UserID.Eq(f.userId)).First()
	if err != nil {
		return err
	}
	f.appId = userMo.AppID

	logger.Info("file upload begin time = ", time.Now().UnixMilli())

	writeEvents := make([]string, 0)
	recordIds := make([]int64, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range f.files {
		if file == nil {
			continue
		}

		wg.Add(1)
		go func(fi *multipart.FileHeader) {
			defer wg.Done()
			// 读文件
			events, err := readFile(fi)
			if err != nil {
				logger.Error("read file failed. err=", err.Error())
				return
			}

			// 检查文件格式
			beginTime, useTime, ok := checkEvents(events)
			if !ok {
				logger.CtxWarnf(ctx, "file check failed")
				return
			}

			// 查询记录是否存在
			record, err := query.Record.WithContext(ctx).Where(query.Record.UserID.Eq(f.userId), query.Record.BeginTime.Eq(beginTime)).First()
			if err == nil && record != nil {
				logger.Warn("使用记录已存在")
				return
			}

			// 插入使用记录
			mo := &model.Record{
				RecordID:  0,
				UserID:    f.userId,
				BeginTime: beginTime,
				UseTime:   useTime,
			}
			mu.Lock()
			err = query.Record.WithContext(ctx).Create(mo)
			if err != nil {
				logger.Error("create record failed. err=", err.Error())
				return
			}
			recordIds = append(recordIds, mo.RecordID)
			mu.Unlock()

			// 转换数据模型
			day := time.UnixMilli(beginTime).Format("2006-01-02")
			dbEvents := f.transEvent(mo.RecordID, day, events)
			mu.Lock()
			writeEvents = append(writeEvents, dbEvents...)
			mu.Unlock()
		}(file)
	}

	wg.Wait()
	logger.Info("file data extract end, time = ", time.Now().UnixMilli())
	// 写入hadoop
	ok := hadoop.WriteEvents(ctx, writeEvents)
	// hadoop 写入失败, mysql 回滚
	if !ok {
		logger.Error("hadoop write failed.")
		_, _ = query.Record.WithContext(ctx).Where(query.Record.RecordID.In(recordIds...)).Delete()
	}

	logger.Info("file upload end time = ", time.Now().UnixMilli())

	return nil
}

func (f *FileUpload) transEvent(recordId int64, day string, events [][]string) []string {
	res := make([]string, 0, len(events))
	for _, event := range events {
		r := &hadoop.Event{
			RecordId: recordId,
			AppId:    int32(f.appId),
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

func (f *FileUpload) GetResp() *backend.UserDataUploadResp {
	return &backend.UserDataUploadResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

func readFile(file *multipart.FileHeader) ([][]string, error) {
	filePtr, err := file.Open()
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

func writeFile(filePath string, fileName string, events [][]string) error {
	_, err := os.Open(filePath)
	if os.IsNotExist(err) {
		err := os.Mkdir(filePath, os.ModeDir)
		if err != nil {
			return err
		}
	} else if err != nil {
		return microtype.DirOpenFailed
	}

	newFile, err := os.Create(filePath + "\\" + fileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(events)
	if err != nil {
		return err
	}

	return nil
}

// 开始时间，使用时长
func checkEvents(events [][]string) (int64, int64, bool) {
	// 大小 首行数据名
	if len(events) < 3 {
		return 0, 0, false
	}

	// 应用启动 应用关闭
	beginEvent := events[1]
	endEvent := events[len(events)-1]
	if len(beginEvent) < 2 || len(endEvent) < 2 || beginEvent[0] != string(event_data.AppStart) || endEvent[0] != string(event_data.AppQuit) {
		return 0, 0, false
	}

	beginTime, err1 := strconv.ParseInt(beginEvent[1], 10, 64)
	endTime, err2 := strconv.ParseInt(endEvent[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}

	return beginTime, endTime - beginTime, true
}

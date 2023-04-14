package upload

import (
	"backend/biz/entity/event_data"
	"backend/biz/entity/user"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"mime/multipart"
	"os"
	"strconv"
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
	ue := user.NewUser(user.IdQuery, nil, nil, f.userId)
	if err := ue.Load(ctx); err != nil {
		return err
	}

	f.appId = ue.GetIdQueryUser().AppID

	for _, file := range f.files {
		if file == nil {
			continue
		}

		// 读文件
		events, err := readFile(file)
		if err != nil {
			return err
		}

		// 检查
		beginTime, ok := checkEvents(events)
		if !ok {
			logger.CtxWarnf(ctx, "file check failed")
			continue
		}

		// 写文件
		fileName := fmt.Sprintf("%d_%d.csv", f.userId, beginTime)
		err = writeFile(fmt.Sprintf("%s\\%d", event_data.EventDataDirPath, f.appId), fileName, events)
		if err != nil {
			logger.CtxErrorf(ctx, "file write failed, err=%s", err.Error())
			continue
		}
	}

	return nil
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

type EventType string

const (
	AppStart EventType = "1"
	AppQuit  EventType = "2"
)

func checkEvents(events [][]string) (int64, bool) {
	// 大小 首行数据名
	if len(events) < 3 {
		return 0, false
	}

	// 应用启动 应用关闭
	beginEvent := events[1]
	endEvent := events[len(events)-1]
	if len(beginEvent) < 2 || len(endEvent) < 2 || beginEvent[0] != string(AppStart) || endEvent[0] != string(AppQuit) {
		return 0, false
	}

	// 开始时间
	beginTime, err1 := strconv.ParseInt(beginEvent[1], 10, 64)
	_, err2 := strconv.ParseInt(endEvent[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, false
	}

	return beginTime, true
}

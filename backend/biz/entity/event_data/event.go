package event_data

import (
	"backend/biz/microtype"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const EventDataDirPath = "D:\\hadoop_data\\event_data"

type Event struct {
	appId int64

	// 数据
	count    map[int64]int64 // uId -> event_record
	filePath []string        // 文件名
}

func NewEvent(appId int64) *Event {
	return &Event{
		appId:    appId,
		count:    make(map[int64]int64),
		filePath: make([]string, 0),
	}
}

func (e *Event) Load(ctx context.Context) error {
	files, err := e.openDir()
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		names := strings.Split(fileName, ".")
		if len(names) < 2 || names[1] != "csv" {
			continue
		}
		nums := strings.Split(names[0], "_")
		if len(nums) < 2 {
			continue
		}

		uIdStr := nums[0]
		uId, err := strconv.ParseInt(uIdStr, 10, 64)
		if err != nil {
			continue
		}

		if v, ok := e.count[uId]; ok {
			e.count[uId] = v + 1
		} else {
			e.count[uId] = 1
		}

		e.filePath = append(e.filePath, fmt.Sprintf(fmt.Sprintf("%s\\%d\\%s", EventDataDirPath, e.appId, fileName)))
	}

	return nil
}

func (e *Event) openDir() ([]os.FileInfo, error) {
	dir, err := os.Open(fmt.Sprintf("%s\\%d", EventDataDirPath, e.appId))
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

func (e *Event) GetEventsRecordNum() map[int64]int64 {
	return e.count
}

func (e *Event) GetFilePath() []string {
	return e.filePath
}

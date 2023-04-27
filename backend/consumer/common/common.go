package common

import (
	"backend/biz/entity/event_data"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"os"
)

func OpenFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	events, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func GetUserEventPath(ctx context.Context, appId int64) map[int64][]string {
	ed := event_data.NewEvent(appId)
	if err := ed.Load(ctx); err != nil {
		logger.Error("fileDataErr=", err.Error())
		return nil
	}

	return ed.GetUId2FilePath()
}

func WriteToDataToPath(m map[int64]map[string]int64, path string, colNames []string) {
	data := make([][]string, 0)
	data = append(data, colNames)
	for userId, desc2Cnt := range m {
		for desc, cnt := range desc2Cnt {
			data = append(data, []string{
				fmt.Sprintf("%d", userId),
				desc,
				fmt.Sprintf("%d", cnt),
			})
		}
	}

	file, err := os.Create(path)
	if err != nil {
		logger.Error("create file failed. err=", err.Error())
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		return
	}

	return
}

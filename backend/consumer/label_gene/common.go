package label_gene

import (
	"backend/biz/entity/event_data"
	"context"
	"encoding/csv"
	"github.com/bytedance/gopkg/util/logger"
	"os"
)

func openFile(path string) ([][]string, error) {
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

func getUserEventPath(ctx context.Context, appId int64) map[int64][]string {
	ed := event_data.NewEvent(appId)
	if err := ed.Load(ctx); err != nil {
		logger.Error("fileDataErr=", err.Error())
		return nil
	}

	return ed.GetUId2FilePath()
}

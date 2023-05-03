package seq_mining

import (
	"archive/zip"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
)

type ResultDownload struct {
	taskId int64

	//
	buf *bytes.Buffer
}

func NewResultDownload(taskId int64) *ResultDownload {
	return &ResultDownload{
		taskId: taskId,
	}
}

func (r *ResultDownload) Load(ctx context.Context) error {
	task, err := query.SeqMiningTask.WithContext(ctx).Where(query.SeqMiningTask.TaskID.Eq(r.taskId)).First()
	if err != nil {
		return err
	}

	if task.Event2number == nil || task.Result == nil {
		return err
	}

	customEvent2Number := make(map[string]int)
	res := make([]*Result, 0)

	err = json.Unmarshal([]byte(*task.Event2number), &customEvent2Number)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(*task.Result), &res)
	if err != nil {
		return err
	}
	// custom event
	eventNumbers := make([]*EventNumber, 0)
	for event, number := range customEvent2Number {
		eventNumbers = append(eventNumbers, &EventNumber{
			Number: number,
			Event:  event,
		})
	}
	sort.Slice(eventNumbers, func(i, j int) bool {
		return eventNumbers[i].Number < eventNumbers[j].Number
	})
	eventNumbersStr := make([][]string, 0)
	eventNumbersStr = append(eventNumbersStr, []string{"编号", "事件"})
	for _, eventNumber := range eventNumbers {
		eventNumbersStr = append(eventNumbersStr, []string{
			fmt.Sprintf("%d", eventNumber.Number),
			fmt.Sprintf("%s", eventNumber.Event),
		})
	}

	// res
	resStr := make([][]string, 0)
	resStr = append(resStr, []string{"出现次数", "编号序列"})
	for _, re := range res {
		resStr = append(resStr, []string{
			fmt.Sprintf("%d", re.Cnt),
			fmt.Sprintf("%v", re.Numbers),
		})
	}

	r.buf = new(bytes.Buffer)
	zipWriter := zip.NewWriter(r.buf)
	defer zipWriter.Close()

	// Add file 1
	file1Writer, _ := zipWriter.Create("event2number.csv")
	writer1 := csv.NewWriter(file1Writer)
	err = writer1.WriteAll(eventNumbersStr)
	if err != nil {
		return err
	}
	writer1.Flush()

	// Add file 2
	file2Writer, _ := zipWriter.Create("result.csv")
	writer2 := csv.NewWriter(file2Writer)
	err = writer2.WriteAll(resStr)
	if err != nil {
		return err
	}
	writer2.Flush()

	return nil
}

func (r *ResultDownload) GetBuf() *bytes.Buffer {
	return r.buf
}

func (r *ResultDownload) GetResp() *backend.SeqMiningResultDownloadResp {
	return &backend.SeqMiningResultDownloadResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

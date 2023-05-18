package seq_mining

import (
	"backend/biz/hadoop"
	"backend/biz/usecase/seq_mining"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"backend/optimize_prefixspan"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"sort"
)

const (
	// prefix db 二维最大长度
	maxSeqLength = 100
	// 过长序列长度
	sameSeqLength = 5
	// 最大数量 最多挖掘10000
	maxResultNum = 10000
)

var (
	// 事件数据 -> 编号
	customEvent2Number = make(map[string]int)
	// 编号起始位置
	curNumber = 0
	// 数据集
	db = make([]optimize_prefixspan.Sequence, 0)
)

func Gene(appId int64, taskId int64) {
	ctx := context.Background()
	logger.Info("seq mining begin.")
	// 初始化
	// 事件数据 -> 编号
	customEvent2Number = make(map[string]int)
	// 编号起始位置
	curNumber = 0
	// 数据集
	db = make([]optimize_prefixspan.Sequence, 0)

	// 查询记录
	seqMiningModel, err := query.SeqMiningTask.WithContext(ctx).Where(query.SeqMiningTask.TaskID.Eq(taskId)).First()
	if err != nil {
		return
	}

	// 更新任务状态
	updateMo := model.SeqMiningTask{
		Status: seq_mining.StatusRun,
	}
	_, err = query.SeqMiningTask.WithContext(ctx).Where(query.SeqMiningTask.TaskID.Eq(taskId)).Updates(updateMo)
	if err != nil {
		return
	}

	percent := seqMiningModel.Percent
	if percent <= 0 {
		return
	}

	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User

	// 查询 mysql 已有记录
	records, err := recordMO.Join(userDO, recordDO.UserID.EqCol(userDO.UserID)).
		Where(userDO.AppID.Eq(appId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query mysql failed. Err=", err.Error())
		return
	}

	for idx, rec := range records {
		// 读取行为原始数据
		events, err := hadoop.QueryEventsByRecordId(ctx, rec.RecordID)
		if err != nil {
			logger.Error(fmt.Sprintf("query hadoop failed. idx=%d, recordId=%d, err=%s", idx, rec.RecordID, err.Error()))
			continue
		}
		logger.Info(fmt.Sprintf("查询第%d次记录完成，reordId=%d, 记录长度=%d", idx, rec.RecordID, len(events)))

		if len(events) <= 0 {
			logger.Warn("event length = 0")
			continue
		}

		numbers := make([]int, 0, len(events))
		// 将原始数据转换为可区分 event 并对其编号
		// 将事件数据转换成编号序列 & 限制过长相同序列
		lastNumber := -1
		sameNumberLength := 0
		for _, event := range events {
			id, ok := getCustomEventNumber(event)
			if !ok {
				continue
			}
			if id == lastNumber {
				sameNumberLength++
				if sameNumberLength >= sameSeqLength {
					continue
				}
			} else {
				sameNumberLength = 0
			}
			lastNumber = id
			numbers = append(numbers, id)
		}
		logger.Info("number length = ", len(numbers))
		// 将编号序列按照最大大小进行划分
		number2Slice := funk.ChunkInts(numbers, maxSeqLength)
		seqs := make([]optimize_prefixspan.Sequence, 0, len(number2Slice))
		for _, slice := range number2Slice {
			seq := optimize_prefixspan.Sequence{}
			seq = append(seq, slice...)
			seqs = append(seqs, seq)
		}

		db = append(db, seqs...)
	}

	// 进行挖掘
	numbers, cnts := optimize_prefixspan.PrefixSpan(db, int(percent)*len(db)/100)
	if len(numbers) != len(cnts) {
		return
	}
	res := make([]*seq_mining.Result, 0, len(cnts))
	for idx := range numbers {
		res = append(res, &seq_mining.Result{
			Cnt:     cnts[idx],
			Numbers: numbers[idx],
		})
	}

	// 升序
	sort.Slice(res, func(i, j int) bool {
		return res[i].Cnt > res[j].Cnt
	})

	//fmt.Printf("%v", res)
	if len(res) > maxResultNum {
		res = res[:maxResultNum]
	}

	// 挖掘结果写入数据库中
	resBs, err := json.Marshal(res)
	if err != nil {
		return
	}

	event2numberBs, err := json.Marshal(customEvent2Number)
	if err != nil {
		return
	}

	updateMo2 := model.SeqMiningTask{
		Status:       seq_mining.StatusEnd,
		Event2number: proto.String(string(event2numberBs)),
		Result:       proto.String(string(resBs)),
	}

	_, err = query.SeqMiningTask.WithContext(ctx).Where(query.SeqMiningTask.TaskID.Eq(taskId)).Updates(updateMo2)
	if err != nil {
		return
	}
	logger.Info("seq mining end.")
}

func getCustomEventNumber(event *hadoop.Event) (int, bool) {
	customEvent := ""
	// 去除鼠标移动事件
	if event.EventType == hadoop.MouseMove {
		return 0, false
	}

	if event.EventType == hadoop.Shortcut {
		customEvent = fmt.Sprintf("%d|%d|%d|%d|%s|%s", event.EventType, event.MouseClickType, event.MouseClickBtn, event.KeyClickType, event.KeyCode, event.ComponentName)
	} else {
		customEvent = fmt.Sprintf("%d|%d|%d|%d|%s|%s", event.EventType, event.MouseClickType, event.MouseClickBtn, event.KeyClickType, "", event.ComponentName)
	}

	if number, ok := customEvent2Number[customEvent]; ok {
		return number, true
	}

	customEvent2Number[customEvent] = curNumber
	curNumber++
	return curNumber - 1, true
}

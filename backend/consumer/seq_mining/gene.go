package seq_mining

import (
	"backend/biz/entity/event_data"
	"backend/biz/usecase/seq_mining"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"backend/consumer/common"
	"backend/optimize_prefixspan"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
	"github.com/thoas/go-funk"
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

func frequentItems(db []optimize_prefixspan.Sequence, minSupport int) []int {
	var list []int // 满足最小支持度的item
	m := make(map[int]int)
	for _, seq := range db {
		exist := make(map[int]bool)
		for _, item := range seq {
			exist[item] = true
		}
		for item, _ := range exist {
			m[item] = m[item] + 1
		}
	}

	for item, cnt := range m {
		if cnt >= minSupport {
			list = append(list, item)
		}
	}

	return list
}

func Gene(appId int64, taskId int64) {
	ctx := context.Background()
	logger.Info("seq mining begin.")

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

	// 数据文件
	ed := event_data.NewEvent(appId)
	if err := ed.Load(ctx); err != nil {
		logger.Error("err=", err.Error())
		return
	}

	paths := ed.GetFilePath()

	for _, path := range paths {
		// 读取行为原始数据
		events, err := common.OpenFile(path)
		if err != nil {
			continue
		}
		numbers := make([]int, 0, len(events))
		// 将原始数据转换为可区分 event 并对其编号
		// 将事件数据转换成编号序列 & 限制过长相同序列
		lastNumber := -1
		sameNumberLength := 0
		for idx, event := range events {
			if idx == 0 { // 去除头部
				continue
			}
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

	list := frequentItems(db, int(percent)*len(db)/100)
	fmt.Printf("%v\n", list)

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

func getCustomEventNumber(event []string) (int, bool) {
	if len(event) <= event_data.ComponentNameIndex {
		return 0, false
	}
	idxs := []int{event_data.EventTypeIndex, event_data.MouseClickTypeIndex, event_data.MouseClickButtonIndex, event_data.KeyClickTypeIndex, event_data.KeyCodeIndex, event_data.ComponentNameIndex}
	customEvent := ""
	// 去除鼠标移动事件
	if event[event_data.EventTypeIndex] == string(event_data.MouseMove) {
		return 0, false
	}
	for _, idx := range idxs {
		// 若是键盘输入且不是快捷键，则忽略键盘输入值
		if idx == event_data.KeyCodeIndex && event[event_data.EventTypeIndex] != string(event_data.Shortcut) {
			continue
		}
		customEvent = customEvent + "|" + event[idx]
	}
	if number, ok := customEvent2Number[customEvent]; ok {
		return number, true
	}

	customEvent2Number[customEvent] = curNumber
	curNumber++
	return curNumber - 1, true
}

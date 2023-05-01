package rule_gene

import (
	"backend/biz/entity/event_data"
	"backend/biz/entity/rule"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"backend/consumer/config"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Gene(appId int64) {
	//defer geneDone(appId)

	ctx := context.Background()
	var (
		wg sync.WaitGroup
		// 错误
		dbDataErr   error
		fileDataErr error
		// 数据
		dbRecords                  []*model.Record             // db 记录
		fileRecord                 map[string]*model.Record    // 文件记录 userId_beginTime -> 记录
		userId2BehaviorDurationMap map[int64][]map[int64]int64 // user_id -> []行为时长map
	)

	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User

	// 查询数据库中的行为记录
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			_ = recover()
		}()

		// 查询已有记录
		res, err := recordMO.Join(userDO, recordDO.UserID.EqCol(userDO.UserID)).
			Where(userDO.AppID.Eq(appId)).Find()
		if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
			dbDataErr = err
			logger.Error("dbDataErr=", dbDataErr.Error())
			return
		}
		dbRecords = res

		s, _ := json.Marshal(dbRecords)
		logger.Info("dbRecords=", string(s))
	}()

	// 查询文件中的行为记录
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			_ = recover()
		}()

		// 文件结果
		fileRecord = make(map[string]*model.Record)
		// 行为时长记录
		userId2BehaviorDurationMap = make(map[int64][]map[int64]int64)

		// 事件规则以及行为规则
		eventRules, behaviorRules, err := rule.GetRuleModels(ctx, appId)
		if err != nil {
			fileDataErr = err
			logger.Error("fileDataErr=", fileDataErr.Error())
			return
		}
		// 规则描述map
		ruleDescMap := make(map[int64]string)
		for _, r := range eventRules {
			if r == nil {
				continue
			}

			ruleDescMap[r.RuleID] = r.RuleDesc
		}
		for _, r := range behaviorRules {
			if r == nil {
				continue
			}

			ruleDescMap[r.RuleID] = r.RuleDesc
		}

		// 数据文件
		ed := event_data.NewEvent(appId)
		if err := ed.Load(ctx); err != nil {
			fileDataErr = err
			logger.Error("fileDataErr=", fileDataErr.Error())
			return
		}

		paths := ed.GetFilePath()
		for _, path := range paths {
			uId, err := getUserId(path)
			if err != nil {
				continue
			}

			events, err := openFile(path)
			if err != nil {
				continue
			}

			record, behaviorDurationMap := process(events, eventRules, behaviorRules, ruleDescMap)
			if record == nil {
				continue
			}

			record.UserID = uId

			key := fmt.Sprintf("%d_%d", uId, record.BeginTime)
			fileRecord[key] = record
			if len(behaviorDurationMap) > 0 {
				if _, ok := userId2BehaviorDurationMap[uId]; !ok {
					userId2BehaviorDurationMap[uId] = make([]map[int64]int64, 0)
				}
				userId2BehaviorDurationMap[uId] = append(userId2BehaviorDurationMap[uId], behaviorDurationMap)
			}
		}

		s, _ := json.Marshal(fileRecord)
		logger.Info("fileRecord=", string(s))
	}()

	wg.Wait()
	if dbDataErr != nil || fileDataErr != nil {
		return
	}

	// 汇总 四类
	// 文件存在 : db不存在 db存在但无行为数据 db存在也有数据 （后两者都要更新）
	// 文件不存在 : 删除 （全量更新，这里做增量即可）
	set1 := make([]*model.Record, 0) // 添加
	set2 := make([]model.Record, 0)  // 更新

	for _, fileRec := range fileRecord {
		if fileRec == nil {
			continue
		}
		exist := false
		for _, dbRec := range dbRecords {
			// db存在
			if dbRec.UserID == fileRec.UserID && dbRec.BeginTime == fileRec.BeginTime {
				set2 = append(set2, *fileRec)
				exist = true
				break
			}
		}
		// db不存在
		if !exist {
			set1 = append(set1, fileRec)
		}
	}

	// 更新
	for _, rec := range set2 {
		_, err := recordMO.Where(recordDO.UserID.Eq(rec.UserID),
			recordDO.BeginTime.Eq(rec.BeginTime)).
			Updates(rec)
		if err != nil {
			logger.Error("update record failed. err=", err.Error())
		}
	}

	// 添加
	err := recordMO.Create(set1...)
	if err != nil {
		logger.Error("create record failed. err=", err.Error())
	}

	// 更新用户的平均时长 && 更新应用的平均使用时长和用户最长使用时长
	aveBehaviorDurationMap := make(map[int64]int64)
	maxBehaviorDurationMap := make(map[int64]int64)
	userNum := int64(0)

	for uId, durationMaps := range userId2BehaviorDurationMap {
		aveDurationMap := make(map[int64]int64)
		for _, durationMap := range durationMaps {
			for ruleId, duration := range durationMap {
				if cnt, ok := aveDurationMap[ruleId]; ok {
					aveDurationMap[ruleId] = cnt + duration
				} else {
					aveDurationMap[ruleId] = duration
				}
			}
		}
		for ruleId, duration := range aveDurationMap {
			aveDurationMap[ruleId] = duration / int64(len(durationMaps))
		}
		if len(aveDurationMap) == 0 {
			continue
		}
		bs, err := json.Marshal(aveDurationMap)
		if err != nil {
			logger.Error("json marshal ave duration map failed. err=", err.Error())
			continue
		}

		mo := model.User{
			BehaviorDurationMap: proto.String(string(bs)),
		}

		_, err = query.User.WithContext(ctx).Where(query.User.UserID.Eq(uId)).Updates(mo)
		if err != nil {
			logger.Error("update user behavior duration failed. err=", err.Error())
		}

		// 应用平均时长 & 用户最长
		userNum++
		for id, duration := range aveDurationMap {
			// 平均
			if cnt, ok := aveBehaviorDurationMap[id]; ok {
				aveBehaviorDurationMap[id] = cnt + duration
			} else {
				aveBehaviorDurationMap[id] = duration
			}

			if cnt, ok := maxBehaviorDurationMap[id]; ok {
				if cnt < duration {
					maxBehaviorDurationMap[id] = duration
				}
			} else {
				maxBehaviorDurationMap[id] = duration
			}
		}
	}

	// 应用平均
	for id, duration := range aveBehaviorDurationMap {
		aveBehaviorDurationMap[id] = duration / userNum
	}

	aveBs, err := json.Marshal(aveBehaviorDurationMap)
	if err != nil {
		logger.Error("json marshal app ave duration map failed. err=", err.Error())
		return
	}
	maxBs, err := json.Marshal(maxBehaviorDurationMap)
	if err != nil {
		logger.Error("json marshal max ave duration map failed. err=", err.Error())
		return
	}
	mo := model.App{
		AveBehaviorDurationMap: proto.String(string(aveBs)),
		MaxBehaviorDurationMap: proto.String(string(maxBs)),
	}

	_, err = query.App.WithContext(ctx).Where(query.App.AppID.Eq(appId)).Updates(mo)
	if err != nil {
		logger.Error("update app behavior failed. err=", err.Error())
		return
	}

	return
}

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

func getUserId(path string) (int64, error) {
	paths := strings.Split(path, "\\")
	if len(paths) == 0 {
		return 0, fmt.Errorf("路径格式错误")
	}
	fileName := paths[len(paths)-1]
	str := strings.Split(fileName, "_")
	if len(str) < 1 {
		return 0, fmt.Errorf("文件名格式错误")
	}
	id, err := strconv.ParseInt(str[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("用户ID解析失败")
	}
	return id, nil
}

func geneDone(appId int64) {
	// running -> stop
	config.StatusChan <- &config.StatusChange{
		AppId:    appId,
		TaskType: config.RuleGene,
		Status:   config.Stop,
	}
}

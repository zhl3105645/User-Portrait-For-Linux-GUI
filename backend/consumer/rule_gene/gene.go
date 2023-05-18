package rule_gene

import (
	"backend/biz/entity/rule"
	"backend/biz/hadoop"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
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
)

func Gene(appId int64) {
	ctx := context.Background()
	var (
		// 数据
		mysqlRecords               []*model.Record             // mysql 记录
		hadoopRecord               map[int64]model.Record      // hadoop记录 record_id -> 记录
		userId2BehaviorDurationMap map[int64][]map[int64]int64 // user_id -> []行为时长map
	)

	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User

	// 查询 mysql 已有记录
	res, err := recordMO.Join(userDO, recordDO.UserID.EqCol(userDO.UserID)).
		Where(userDO.AppID.Eq(appId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query mysql failed. Err=", err.Error())
		return
	}
	mysqlRecords = res

	s, _ := json.Marshal(mysqlRecords)
	logger.Info("mysqlRecords=", string(s))

	// 事件规则以及行为规则
	eventRules, behaviorRules, err := rule.GetRuleModels(ctx, appId)
	if err != nil {
		logger.Error("query rule failed. Err=", err.Error())
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

	// hadoop结果
	hadoopRecord = make(map[int64]model.Record)
	// 行为时长记录
	userId2BehaviorDurationMap = make(map[int64][]map[int64]int64)

	for idx, mo := range mysqlRecords {
		events, err := hadoop.QueryEventsByRecordId(ctx, mo.RecordID)
		if err != nil {
			logger.Error(fmt.Sprintf("query hadoop failed. idx=%d, recordId=%d, err=%s", idx, mo.RecordID, err.Error()))
			continue
		}
		logger.Info(fmt.Sprintf("查询第%d次记录完成，reordId=%d, 记录长度=%d", idx, mo.RecordID, len(events)))

		if len(events) <= 0 {
			logger.Warn("event length = 0")
			continue
		}

		eventRuleData, behaviorRuleData, behaviorDurationMap := process(events, eventRules, behaviorRules, ruleDescMap)
		if eventRuleData == "" || behaviorRuleData == "" {
			continue
		}

		hadoopRecord[mo.RecordID] = model.Record{
			EventRuleValue:    proto.String(eventRuleData),
			BehaviorRuleValue: proto.String(behaviorRuleData),
		}
		if len(behaviorDurationMap) > 0 {
			if _, ok := userId2BehaviorDurationMap[mo.UserID]; !ok {
				userId2BehaviorDurationMap[mo.UserID] = make([]map[int64]int64, 0)
			}
			userId2BehaviorDurationMap[mo.UserID] = append(userId2BehaviorDurationMap[mo.UserID], behaviorDurationMap)
		}
	}

	s, _ = json.Marshal(hadoopRecord)
	logger.Info("hadoopRecord=", string(s))

	// 更新
	for id, rec := range hadoopRecord {
		_, err = recordMO.Where(recordDO.RecordID.Eq(id)).Updates(rec)
		if err != nil {
			logger.Error("update record failed. err=", err.Error())
		}
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

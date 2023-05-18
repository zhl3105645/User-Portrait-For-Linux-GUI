package label_gene

import (
	"backend/biz/entity/rule"
	"backend/biz/hadoop"
	"backend/biz/util"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
	"strings"
)

func ProcessGitDesc(ctx context.Context, appId int64) map[int64]string {
	res := make(map[int64]string)

	recordDO := query.Record
	recordMO := recordDO.WithContext(ctx)
	userDO := query.User
	// 查询 mysql 已有记录
	records, err := recordMO.Join(userDO, recordDO.UserID.EqCol(userDO.UserID)).
		Where(userDO.AppID.Eq(appId)).Find()
	if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
		logger.Error("query mysql failed. Err=", err.Error())
		return nil
	}
	userId2Records := make(map[int64][]*model.Record)
	for _, rec := range records {
		if len(userId2Records[rec.UserID]) == 0 {
			userId2Records[rec.UserID] = make([]*model.Record, 0)
		}
		userId2Records[rec.UserID] = append(userId2Records[rec.UserID], rec)
	}

	// 采集各个用户的git数据
	gitMsgCntMap := make(map[int64]map[string]int64) // user_id -> git msg -> cnt
	for userId, recs := range userId2Records {
		gitMsgCnt := make(map[string]int64) // git msg -> cnt
		for _, rec := range recs {
			events, err := hadoop.QueryEventsByRecordId(ctx, rec.RecordID)
			if err != nil {
				logger.Error(fmt.Sprintf("query hadoop failed. recordId=%d, err=%s", rec.RecordID, err.Error()))
				continue
			}
			logger.Info(fmt.Sprintf("reordId=%d, 记录长度=%d", rec.RecordID, len(events)))

			for msg, cnt := range getGitMsgCntMap(events) {
				if v, ok := gitMsgCnt[msg]; ok {
					gitMsgCnt[msg] = v + cnt
				} else {
					gitMsgCnt[msg] = cnt
				}
			}
		}
		gitMsgCntMap[userId] = gitMsgCnt
	}

	// 写入文件，作为记录
	// common.WriteToDataToPath(gitMsgCntMap, "D:\\graudation2\\code\\backend\\consumer\\label_gene\\data\\git_msg.csv", []string{"user_id", "git_msg", "msg_cnt"})

	// 规则匹配
	aveScoreMap := make(map[int64]float64)
	for userId, desc2Cnt := range gitMsgCntMap {
		score := int64(0)
		total := int64(0)
		for desc, cnt := range desc2Cnt {
			if desc == "" {
				continue
			}
			score += cnt * scoreGitMsg(desc)
			total += cnt
		}
		if total <= 0 {
			continue
		}
		aveScoreMap[userId] = float64(score) / float64(total)
	}

	// 分数分级
	gradeMap := util.GradeByPercent(aveScoreMap, []float64{0.3, 0.7})
	for userId, grade := range gradeMap {
		res[userId] = fmt.Sprintf("%d", grade)
	}

	return res
}

// scoreGitMsg 根据长度和关键词计算权重分数
func scoreGitMsg(msg string) int64 {
	score := int64(1)
	keyWords := []string{
		"fix", "feat", "build",
		"chore", "ci", "docs",
		"style", "refactor", "perf",
		"test", "BREAKING CHANGE", "enhancement"}
	msg = strings.ToLower(msg)
	for _, keyWord := range keyWords {
		if strings.Contains(msg, keyWord) {
			score += 1
			msg = strings.Replace(msg, keyWord, "", 1)
			break
		}
	}

	length := util.GetCharNumberOfString(msg)
	if length >= 4 {
		score += 1
	}

	return score
}

func getGitMsgCntMap(events []*hadoop.Event) map[string]int64 {
	res := make(map[string]int64)

	msgLineEdit := []string{
		"(5|||1||MainWindow.<class_name=QInputDialog>.<class_name=QLineEdit>)",
		"(5|||2||MainWindow.<class_name=QInputDialog>.<class_name=QLineEdit>)",
		"(4|||||MainWindow.<class_name=QInputDialog>.<class_name=QLineEdit>)",
	}
	commitBtn := []string{"(3|1|1|||MainWindow.<class_name=QInputDialog>.<class_name=QDialogButtonBox>.<class_name=QPushButton>)"}

	curMsg := "" // 当前输入框内容
	for _, event := range events {
		if rule.MatchEvent(event, msgLineEdit) {
			// 更新输入框内容
			curMsg = event.ComponentExtra
		} else if rule.MatchEvent(event, commitBtn) {
			if curMsg == "" {
				continue
			}
			if cnt, ok := res[curMsg]; ok {
				res[curMsg] = cnt + 1
			} else {
				res[curMsg] = 1
			}

			curMsg = "" // 清空当前状态
		}
	}

	return res
}

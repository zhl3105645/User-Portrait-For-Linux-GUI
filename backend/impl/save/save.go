package save

import (
	"backend/cmd/dal"
	"backend/cmd/dal/query"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"os"
)

func LoadBehaviorDurationData() {
	//
	modelId := int64(15)

	dal.Init()
	res, err := query.ModelDatum.WithContext(context.Background()).Where(query.ModelDatum.ModelID.Eq(modelId)).Find()
	if err != nil {
		logger.Error("err=", err.Error())
		return
	}

	strs := make([][]string, 0)
	strs = append(strs, []string{
		"user_id",
		"data",
	})
	for _, r := range res {
		strs = append(strs, []string{
			fmt.Sprintf("%d", r.UserID),
			r.Data,
		})
	}

	path := "D:\\毕设2\\code\\ml\\cluster\\data\\origin_data.csv"
	file, err := os.Create(path)
	if err != nil {
		logger.Error("err=", err.Error())
		return
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(strs)
	if err != nil {
		logger.Error("err=", err.Error())
		return
	}

	return
}

func LoadEventSeqData() {
	//
	appId := int64(2)

	dal.Init()
	res, err := query.Record.WithContext(context.Background()).
		Join(query.User, query.User.UserID.EqCol(query.Record.UserID)).
		Where(query.User.AppID.Eq(appId)).Find()
	if err != nil {
		logger.Error("err=", err.Error())
		return
	}

	strs := make([][]string, 0)
	strs = append(strs, []string{
		"user_id",
		"data",
	})
	for _, r := range res {
		if r.BehaviorRuleValue == nil {
			continue
		}
		strs = append(strs, []string{
			fmt.Sprintf("%d", r.UserID),
			*r.EventRuleValue,
		})
	}

	path := "D:\\毕设2\\code\\ml\\predict\\data\\origin_data.csv"
	file, err := os.Create(path)
	if err != nil {
		logger.Error("err=", err.Error())
		return
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(strs)
	if err != nil {
		logger.Error("err=", err.Error())
		return
	}

	return
}

package label_gene

import (
	"backend/consumer/common"
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"strconv"
)

func processBehaviorPrefer(ctx context.Context) map[int64]string {
	res := make(map[int64]string)
	path := "D:\\graudation2\\code\\ml\\cluster\\data\\result.csv"
	data, err := common.OpenFile(path)
	if err != nil {
		logger.Error("open file failed. err=", err.Error())
		return res
	}

	for _, d := range data {
		if len(d) < 2 {
			logger.Error("data length < 2")
			continue
		}
		id, err1 := strconv.ParseInt(d[0], 10, 64)
		category, err2 := strconv.ParseInt(d[1], 10, 64)
		if err1 != nil || err2 != nil {
			logger.Error("data parse int failed.")
			continue
		}

		res[id] = fmt.Sprintf("%d", category+1)
	}

	return res
}

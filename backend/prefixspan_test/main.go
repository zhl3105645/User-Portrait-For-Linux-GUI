package main

import (
	"backend/non_cnt_optimize_prefixspan"
	"backend/optimize_prefixspan"
	"backend/prefixspan"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"io/ioutil"
	"math/rand"
	"time"
)

func main() {
	var db []optimize_prefixspan.Sequence

	db = getData()
	i := 10
	logger.Info("*************min support = ", i)
	optimize(i, db)
	origin(i, db)

}

func geneData(dbLength int, seqLength int) []optimize_prefixspan.Sequence {
	db := make([]optimize_prefixspan.Sequence, 0)
	for i := 0; i < dbLength; i++ {
		seq := optimize_prefixspan.Sequence{}
		for j := 0; j < seqLength; j++ {
			seq = append(seq, rand.Intn(100))
		}
		db = append(db, seq)
	}

	return db
}

func getData() []optimize_prefixspan.Sequence {
	data, err := ioutil.ReadFile("D:\\graudation2\\code\\ml\\prefixspan\\db.txt")
	if err != nil {
		return nil
	}
	db := make([]optimize_prefixspan.Sequence, 0)
	err = json.Unmarshal(data, &db)
	if err != nil {
		return nil
	}

	return db
}

func optimize(percent int, db []optimize_prefixspan.Sequence) {
	// 进行挖掘
	before := time.Now()
	res, _ := optimize_prefixspan.PrefixSpan(db, int(percent)*len(db)/100)
	after := time.Now()

	duration := after.Sub(before).Milliseconds()
	logger.Info("optimize duration = ", duration)
	logger.Info("optimize res length = ", len(res))
}

func origin(percent int, db []optimize_prefixspan.Sequence) {
	// 进行挖掘
	newDB := transDB(db)
	before := time.Now()
	res := prefixspan.PrefixSpan(newDB, int(percent)*len(db)/100)
	after := time.Now()

	duration := after.Sub(before).Milliseconds()
	logger.Info("origin duration = ", duration)
	logger.Info("origin res length = ", len(res))
}

func transDB2(db []optimize_prefixspan.Sequence) []non_cnt_optimize_prefixspan.Sequence {
	newDB := make([]non_cnt_optimize_prefixspan.Sequence, 0, len(db))
	for _, seq := range db {
		newSeq := non_cnt_optimize_prefixspan.Sequence{}
		for _, item := range seq {
			newSeq = append(newSeq, item)
		}
		newDB = append(newDB, newSeq)
	}
	return newDB
}

func transDB(db []optimize_prefixspan.Sequence) []prefixspan.Sequence {
	newDB := make([]prefixspan.Sequence, 0, len(db))
	for _, seq := range db {
		newSeq := prefixspan.Sequence{}
		for _, item := range seq {
			newSeq = append(newSeq, prefixspan.ItemSet{item})
		}
		newDB = append(newDB, newSeq)
	}
	return newDB
}

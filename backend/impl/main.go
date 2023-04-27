package main

import (
	"backend/impl/rule"
	"backend/impl/save"
	"context"
)

func main() {
	//path1 := "D:\\毕设2\\data\\20230405\\1680680589144.csv"
	//path2 := "D:\\毕设2\\data\\20230405\\1680680832299.csv"
	//all_component := "D:\\毕设2\\data\\all_component\\1680791385825.csv"
	//paths := make([]string, 0)
	//
	//dir := "D:\\毕设2\\data\\test_rule"
	//d, err := os.Open(dir)
	//if err != nil {
	//	panic(any(err.Error()))
	//}
	//defer d.Close()
	//
	//files, err := d.Readdir(-1)
	//for _, info := range files {
	//	if info == nil {
	//		continue
	//	}
	//	paths = append(paths, dir+"\\"+info.Name())
	//}
	//
	//// 组件
	//componentMap := make(map[string]*singleUse.QTComponent)
	//// 事件规则，行为规则
	//
	//eventRules, behaviorRules := rule.GetRules()
	//
	//UserUse.Process(paths, componentMap, eventRules, behaviorRules)
	//saveData()
	//Post()
	//loadDataSource()
}

func loadRule() {
	rule.LoadRuleToDatabase(context.Background())
}

func loadDataSource() {
	rule.LoadDataSourceToDatabase(context.Background(), 2)
}

func saveData() {
	//save.LoadBehaviorDurationData()
	save.LoadEventSeqData()
}

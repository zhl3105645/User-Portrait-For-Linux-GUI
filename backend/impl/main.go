package main

import (
	"backend/impl/UserUse"
	"backend/impl/rule"
	"backend/impl/singleUse"
)

func main() {
	path1 := "D:\\毕设2\\data\\20230405\\1680680589144.csv"
	path2 := "D:\\毕设2\\data\\20230405\\1680680832299.csv"
	all_component := "D:\\毕设2\\data\\all_component\\1680791385825.csv"
	paths := []string{path1, path2, all_component}

	// 组件
	componentMap := make(map[string]*singleUse.QTComponent)
	// 事件规则，行为规则

	eventRules, behaviorRules := rule.GetRules()

	UserUse.Process(paths, componentMap, eventRules, behaviorRules)
}

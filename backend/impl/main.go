package main

import (
	"backend/impl/UserUse"
	"backend/impl/rule"
	"backend/impl/singleUse"
	"os"
)

func main() {
	paths := make([]string, 0)

	dir := "D:\\毕设2\\data\\test_rule"
	d, err := os.Open(dir)
	if err != nil {
		panic(any(err.Error()))
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	for _, info := range files {
		if info == nil {
			continue
		}
		paths = append(paths, dir+"\\"+info.Name())
	}

	// 组件
	componentMap := make(map[string]*singleUse.QTComponent)
	// 事件规则，行为规则

	eventRules, behaviorRules := rule.GetRules()

	UserUse.Process(paths, componentMap, eventRules, behaviorRules)
}

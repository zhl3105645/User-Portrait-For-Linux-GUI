package UserUse

import (
	"backend/impl/rule"
	"backend/impl/singleUse"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
	"strings"
)

type UserBasicBehavior struct {
	Average        *singleUse.BasicBehavior   `yaml:"Average"`
	BasicBehaviors []*singleUse.BasicBehavior `yaml:"BasicBehaviors"`
}

func Process(paths []string, componentMap map[string]*singleUse.QTComponent, eventRules []*rule.EventRule, behaviorRules []*rule.BehaviorRule) {
	// 读取已有component
	comFile, err1 := os.Open("./impl/component_gene.yaml")
	if err1 != nil {
		return
	}
	defer comFile.Close()

	existUi := &singleUse.AppComponent{}
	uiDecoder := yaml.NewDecoder(comFile)
	if err := uiDecoder.Decode(existUi); err != nil {
		return
	}

	for _, info := range existUi.Components {
		if info == nil {
			continue
		}
		componentMap[info.Name] = info
	}

	// 开始执行
	basicBehaviors := make([]*singleUse.BasicBehavior, 0, len(paths))
	ave := &singleUse.BasicBehavior{}
	behaviorTimeMap := make(map[int64]int64)
	for _, path := range paths {
		basic := singleUse.Process(path, componentMap, eventRules, behaviorRules)
		if basic != nil {
			basicBehaviors = append(basicBehaviors, basic)

			ave.UseTimeMS += basic.UseTimeMS
			ave.MouseClickCnt += basic.MouseClickCnt
			ave.MouseMoveCnt += basic.MouseMoveCnt
			ave.MouseWheelCnt += basic.MouseWheelCnt
			ave.KeyClickCnt += basic.KeyClickCnt
			ave.ShortcutCnt += basic.ShortcutCnt
			ave.MouseMoveDis += basic.MouseMoveDis
			ave.KeyClickSpeed += basic.KeyClickSpeed
			for id, add := range basic.BehaviorTime {
				if v, ok := behaviorTimeMap[id]; ok {
					behaviorTimeMap[id] = v + add
				} else {
					behaviorTimeMap[id] = add
				}
			}
		}
	}

	for id, v := range behaviorTimeMap {
		behaviorTimeMap[id] = v / int64(len(basicBehaviors))
	}

	ave.UseTimeMS /= int64(len(basicBehaviors))
	ave.MouseClickCnt /= int64(len(basicBehaviors))
	ave.MouseMoveCnt /= int64(len(basicBehaviors))
	ave.MouseWheelCnt /= int64(len(basicBehaviors))
	ave.KeyClickCnt /= int64(len(basicBehaviors))
	ave.ShortcutCnt /= int64(len(basicBehaviors))
	ave.MouseMoveDis /= float64(len(basicBehaviors))
	ave.KeyClickSpeed /= float64(len(basicBehaviors))
	ave.BehaviorTime = behaviorTimeMap

	// 用户维度基础数据
	userBasicFile, err := os.Create("./impl/UserUse/multiUse.yaml")
	if err != nil {
		println(err.Error())
		return
	}
	defer userBasicFile.Close()

	data := &UserBasicBehavior{
		Average:        ave,
		BasicBehaviors: basicBehaviors,
	}

	encoderUser := yaml.NewEncoder(userBasicFile)

	if err := encoderUser.Encode(data); err != nil {
		println(err.Error())
	}

	// ui 组件数据
	uiComponentFile, err := os.Create("./impl/component_gene.yaml")
	if err != nil {
		println(err.Error())
		return
	}
	defer uiComponentFile.Close()

	coms := make([]*singleUse.QTComponent, 0, len(componentMap))
	for _, com := range componentMap {
		if com == nil {
			continue
		}

		coms = append(coms, com)
	}

	sort.Slice(coms, func(i, j int) bool {
		return strings.Compare(coms[i].Name, coms[j].Name) < 0
	})

	ui := &singleUse.AppComponent{
		Components: coms,
	}

	encoderComs := yaml.NewEncoder(uiComponentFile)

	if err := encoderComs.Encode(ui); err != nil {
		println(err.Error())
	}
}

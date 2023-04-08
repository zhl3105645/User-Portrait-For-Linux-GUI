package UserUse

import (
	"backend/impl/rule"
	"backend/impl/singleUse"
	"gopkg.in/yaml.v2"
	"os"
)

type UserBasicBehavior struct {
	Average        *singleUse.BasicBehavior   `yaml:"Average"`
	BasicBehaviors []*singleUse.BasicBehavior `yaml:"BasicBehaviors"`
}

func Process(paths []string, componentMap map[string]*singleUse.QTComponent, eventRules []*rule.EventRule, behaviorRules []*rule.BehaviorRule) {
	basicBehaviors := make([]*singleUse.BasicBehavior, 0, len(paths))
	ave := &singleUse.BasicBehavior{}
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
		}
	}

	ave.UseTimeMS /= int64(len(basicBehaviors))
	ave.MouseClickCnt /= int64(len(basicBehaviors))
	ave.MouseMoveCnt /= int64(len(basicBehaviors))
	ave.MouseWheelCnt /= int64(len(basicBehaviors))
	ave.KeyClickCnt /= int64(len(basicBehaviors))
	ave.ShortcutCnt /= int64(len(basicBehaviors))
	ave.MouseMoveDis /= float64(len(basicBehaviors))
	ave.KeyClickSpeed /= float64(len(basicBehaviors))

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
	uiComponentFile, err := os.Create("./impl/component.yaml")
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

	ui := &singleUse.AppComponent{
		Components: coms,
	}

	encoderComs := yaml.NewEncoder(uiComponentFile)

	if err := encoderComs.Encode(ui); err != nil {
		println(err.Error())
	}
}

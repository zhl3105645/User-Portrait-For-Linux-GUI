package rule

import (
	"gopkg.in/yaml.v2"
	"os"
)

type EventRules struct {
	Rules []*EventRule `yaml:"Rules"`
}

type EventRule struct {
	Id     int      `yaml:"Id"`
	Name   string   `yaml:"Name"`
	Events []string `yaml:"Events"`
}

type BehaviorRules struct {
	Rules []*BehaviorRule `yaml:"Rules"`
}

type BehaviorRule struct {
	Id        int      `yaml:"Id"`
	Name      string   `yaml:"Name"`
	Behaviors []string `yaml:"Behaviors"`
}

func GetRules() ([]*EventRule, []*BehaviorRule) {
	eventRuleFile, err1 := os.Open("./impl/rule/event_rule.yaml")
	behaviorRuleFile, err2 := os.Open("./impl/rule/behavior_rule.yaml")
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	defer eventRuleFile.Close()
	defer behaviorRuleFile.Close()

	eventRules := &EventRules{}
	eventDecoder := yaml.NewDecoder(eventRuleFile)
	if err := eventDecoder.Decode(eventRules); err != nil {
		return nil, nil
	}

	behaviorRules := &BehaviorRules{}
	behaviorDecoder := yaml.NewDecoder(behaviorRuleFile)
	if err := behaviorDecoder.Decode(behaviorRules); err != nil {
		return eventRules.Rules, nil
	}

	return eventRules.Rules, behaviorRules.Rules
}

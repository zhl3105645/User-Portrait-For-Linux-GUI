package config

import (
	"github.com/bytedance/gopkg/util/logger"
	"gopkg.in/yaml.v2"
	"os"
)

// 使用 yaml 文件作为消息队列，有同时写的风险

var StatusChan chan *StatusChange

type TaskType int

const (
	ComponentGene     TaskType = 1 // 生成组件数据
	BasicBehaviorGene TaskType = 2 // 生成基础行为数据
	RuleGene          TaskType = 3 // 生成规则数据
	ModelGene         TaskType = 4 // 生成模型数据
	LabelGene         TaskType = 5 // 生成标签数据
)

type Status int

// Stop -> Begin : 发出信号
// Begin -> Running : 开始运行
// Running -> Stop : 运行完成
const (
	Stop    Status = 0 // 停止
	Begin   Status = 1 // 开始
	Running Status = 2 // 运行
)

type Configs struct {
	Configs map[int64]*Config `yaml:"Configs"`
}

type Config struct {
	Config map[TaskType]Status `yaml:"Config"` // 状态
	Param  map[TaskType]int64  `yaml:"Param"`  // 参数
}

func ReadConfig() (*Configs, error) {
	configFile, err := os.Open("D:\\graudation2\\code\\backend\\consumer\\task.yaml")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	configs := &Configs{}

	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(configs); err != nil {
		return nil, err
	}

	return configs, nil
}

func WriteConfig(configs *Configs) error {
	configFile, err := os.Create("D:\\graudation2\\code\\backend\\consumer\\task.yaml")
	if err != nil {
		return err
	}
	defer configFile.Close()

	decoder := yaml.NewEncoder(configFile)
	if err := decoder.Encode(configs); err != nil {
		return err
	}

	return nil
}

type StatusChange struct {
	AppId    int64
	TaskType TaskType
	Status   Status
}

func ReceiveStatusChange() {
	for change := range StatusChan {
		if change == nil {
			continue
		}
		configs, err := ReadConfig()
		if err != nil {
			logger.Error("read yaml config failed. err=", err.Error())
			continue
		}

		if configs == nil || len(configs.Configs) == 0 ||
			configs.Configs[change.AppId] == nil || len(configs.Configs[change.AppId].Config) == 0 {
			continue
		}

		configs.Configs[change.AppId].Config[change.TaskType] = change.Status
		err = WriteConfig(configs)
		if err != nil {
			logger.Error("write yaml config failed. err=", err.Error())
			continue
		}
	}
}

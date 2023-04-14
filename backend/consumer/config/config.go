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
	ComponentGene TaskType = 1 // 整合组件
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
	AppConfigs map[int64]*AppConfig `yaml:"AppConfigs"`
}

type AppConfig struct {
	Config map[TaskType]Status `yaml:"Config"`
}

func ReadConfig() (*Configs, error) {
	configFile, err := os.Open("D:\\毕设2\\code\\backend\\consumer\\task.yaml")
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
	configFile, err := os.Create("D:\\毕设2\\code\\backend\\consumer\\task.yaml")
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

		if configs == nil || len(configs.AppConfigs) == 0 ||
			configs.AppConfigs[change.AppId] == nil || len(configs.AppConfigs[change.AppId].Config) == 0 {
			continue
		}

		configs.AppConfigs[change.AppId].Config[change.TaskType] = change.Status
		err = WriteConfig(configs)
		if err != nil {
			logger.Error("write yaml config failed. err=", err.Error())
			continue
		}
	}
}

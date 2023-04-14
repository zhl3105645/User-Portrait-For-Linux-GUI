package component_gene

import (
	"backend/biz/entity/component"
	"backend/biz/entity/event_data"
	"backend/cmd/dal/model"
	"backend/consumer/config"
	"context"
	"encoding/csv"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang/protobuf/proto"
	"os"
	"strconv"
	"sync"
)

func Gene(appId int64) {
	ctx := context.Background()
	var (
		wg sync.WaitGroup
		// 错误
		existsErr error
		curErr    error
		// 数据
		existComs []*model.Component          // 已有组件
		curComs   map[string]*model.Component // 目前组件, 组件名 -> 组件
	)
	// 查询已有组件信息
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			_ = recover()
		}()

		cp := component.NewComponent(component.QueryAll, appId, nil, nil)
		if err := cp.Load(ctx); err != nil {
			existsErr = err
			logger.Error("existsErr=", existsErr.Error())
			return
		}

		existComs = cp.GetQueryComponent()
		logger.Info("existComs=", existComs)
	}()

	// 目前的组件信息
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			_ = recover()
		}()

		curComs = make(map[string]*model.Component)

		ed := event_data.NewEvent(appId)
		if err := ed.Load(ctx); err != nil {
			curErr = err
			logger.Error("curErr=", curErr.Error())
			return
		}

		paths := ed.GetFilePath()
		for _, path := range paths {
			events, err := OpenFile(path)
			if err != nil {
				continue
			}

			for _, event := range events {
				if event_data.ComponentNameIndex > len(event)-1 {
					continue
				}
				com := &model.Component{
					AppID:         appId,
					ComponentName: event[event_data.ComponentNameIndex],
				}

				if event_data.ComponentTypeIndex <= len(event)-1 {
					typ, err := strconv.ParseInt(event[event_data.ComponentTypeIndex], 10, 64)
					if err != nil {
						continue
					}
					com.ComponentType = typ
				}
				if event_data.ComponentExtraIndex <= len(event)-1 {
					com.ComponentDesc = proto.String(event[event_data.ComponentExtraIndex])
				}

				if _, ok := curComs[com.ComponentName]; !ok {
					curComs[com.ComponentName] = com
				}
			}
		}

		logger.Info("curComs=", curComs)

		return
	}()

	wg.Wait()
	if existsErr != nil || curErr != nil {
		GeneDone(appId)
		return
	}

	// 汇总
	addCom := make([]*model.Component, 0, len(curComs)-len(existComs))
	for name, com := range curComs {
		if com == nil {
			continue
		}

		exist := false
		for _, com1 := range existComs {
			if com1 == nil {
				continue
			}

			if com1.ComponentName == name {
				exist = true
				break
			}
		}

		if !exist {
			addCom = append(addCom, com)
		}
	}

	cp := component.NewComponent(component.InsertBatch, appId, nil, &component.InsertParam{
		InsertMo: addCom,
	})
	_ = cp.Load(ctx)

	GeneDone(appId)
	return
}

func OpenFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	events, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func GeneDone(appId int64) {
	// running -> stop
	config.StatusChan <- &config.StatusChange{
		AppId:    appId,
		TaskType: config.ComponentGene,
		Status:   config.Stop,
	}
}

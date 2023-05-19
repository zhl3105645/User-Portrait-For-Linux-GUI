package hadoop

import (
	"backend/gohive"
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/thoas/go-funk"
	"os"
	"sort"
	"strings"
	"time"
)

var HiveConnection *gohive.Connection

type Event struct {
	RecordId       int64
	EventType      int32
	EventTime      int64
	MousePos       string
	MouseClickType int32
	MouseClickBtn  int32
	MouseMoveType  int32
	KeyClickType   int32
	KeyCode        string
	ComponentName  string
	ComponentType  int32
	ComponentExtra string
	AppId          int32
	Day            string
}

func (e *Event) String() string {
	return fmt.Sprintf("(%d, %d, %d, %q, %d,%d,%d,%d,%q,%q, %d, %q, %d,%q)",
		e.RecordId, e.EventType, e.EventTime, e.MousePos, e.MouseClickType, e.MouseClickBtn, e.MouseMoveType,
		e.KeyClickType, e.KeyCode, e.ComponentName, e.ComponentType, e.ComponentExtra, e.AppId, e.Day)
}

func Init(ctx context.Context) {
	configuration := gohive.NewConnectConfiguration()
	// If it's not set it will be picked up from the logged user
	configuration.Username = "zhl"
	// This may not be necessary
	configuration.Password = "zhl"
	connection, errConn := gohive.Connect("192.168.81.131", 10000, "NONE", configuration)

	if errConn != nil {
		logger.Error(errConn)
		return
	}

	// 使用profile数据库
	cursor := connection.Cursor()

	cursor.Exec(ctx, "use profile")
	if cursor.Err != nil {
		logger.Error(cursor.Err)
		return
	}

	HiveConnection = connection
}

func QueryEventsByRecordId(ctx context.Context, recordId int64) ([]*Event, error) {
	if recordId <= 0 {
		return nil, fmt.Errorf("recordId len <= 0")
	}
	cursor := HiveConnection.Cursor()
	defer cursor.Close()

	sql := fmt.Sprintf("select * from event where record_id = %d", recordId)
	cursor.Exec(ctx, sql)
	if cursor.Err != nil {
		return nil, cursor.Err
	}

	res := make([]*Event, 0)
	for cursor.HasMore(ctx) {
		var recordId2 int64
		var eventType int32
		var eventTime int64
		var mousePos string
		var mouseClickType int32
		var mouseClickBtn int32
		var mouseMoveType int32
		var keyClickType int32
		var keyCode string
		var componentName string
		var componentType int32
		var componentExtra string
		var appId int32
		var day string

		cursor.FetchOne(ctx,
			&recordId2, &eventType,
			&eventTime, &mousePos, &mouseClickType,
			&mouseClickBtn, &mouseMoveType, &keyClickType,
			&keyCode, &componentName, &componentType,
			&componentExtra, &appId, &day)
		if cursor.Err != nil {
			logger.Error("hive sql fetch one failed. err=", cursor.Err.Error())
			continue
		}

		event := &Event{
			RecordId:       recordId2,
			EventType:      eventType,
			EventTime:      eventTime,
			MousePos:       mousePos,
			MouseClickType: mouseClickType,
			MouseClickBtn:  mouseClickBtn,
			MouseMoveType:  mouseMoveType,
			KeyClickType:   keyClickType,
			KeyCode:        keyCode,
			ComponentName:  componentName,
			ComponentType:  componentType,
			ComponentExtra: componentExtra,
			AppId:          appId,
			Day:            day,
		}

		res = append(res, event)
	}

	// 按照事件排序
	sort.Slice(res, func(i, j int) bool {
		return res[i].EventTime < res[j].EventTime
	})

	return res, nil
}

func WriteEvents(ctx context.Context, events []string) bool {
	if len(events) <= 0 {
		return false
	}
	cursor := HiveConnection.Cursor()
	defer cursor.Close()

	errCnt := 0
	// 一次写入大小
	eventsSlice := funk.ChunkStrings(events, 5000)
	logger.Info("event slice length = ", len(eventsSlice))
	for idx, slice := range eventsSlice {
		sql := "insert into table event values" + strings.Join(slice, ",")

		cursor.Exec(ctx, sql)
		if cursor.Err != nil {
			logger.Error("cursor exec failed. err=", cursor.Err.Error())
			writeSQL(idx, sql)
			errCnt++
			continue
		}
		logger.Info(fmt.Sprintf("第%d次写入成功, time = %d", idx, time.Now().UnixMilli()))
	}

	return errCnt == 0
}

func writeSQL(idx int, sql string) {
	dir := "D:\\graudation2\\code\\backend\\biz\\hadoop\\sql"
	file, err := os.Create(fmt.Sprintf("%s\\sql%d.sql", dir, idx))
	if err != nil {
		logger.Error("create file failed. err=", err.Error())
		return
	}

	_, err = file.WriteString(sql)
	if err != nil {
		logger.Error("write file failed. err=", err.Error())
		return
	}

	return
}

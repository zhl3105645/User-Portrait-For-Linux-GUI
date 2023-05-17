package hadoop

import (
	"backend/gohive"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"log"
	"testing"
)

func TestName(t *testing.T) {
	ctx := context.Background()

	configuration := gohive.NewConnectConfiguration()
	// If it's not set it will be picked up from the logged user
	configuration.Username = "zhl"
	// This may not be necessary
	configuration.Password = "zhl"
	connection, errConn := gohive.Connect("192.168.81.131", 10000, "NONE", configuration)

	if errConn != nil {
		log.Fatal(errConn)
	}
	cursor := connection.Cursor()

	cursor.Exec(ctx, "use profile")
	if cursor.Err != nil {
		log.Fatal(cursor.Err)
	}

	cursor.Exec(ctx, "select * from event")
	if cursor.Err != nil {
		log.Fatal(cursor.Err)
	}

	var recordId int64
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
	for cursor.HasMore(ctx) {
		cursor.FetchOne(ctx,
			&recordId, &eventType,
			&eventTime, &mousePos, &mouseClickType,
			&mouseClickBtn, &mouseMoveType, &keyClickType,
			&keyCode, &componentName, &componentType,
			&componentExtra, &appId, &day)
		if cursor.Err != nil {
			log.Fatal(cursor.Err)
		}

	}

	cursor.Close()
	connection.Close()
}

func TestQuery(t *testing.T) {
	ctx := context.Background()
	Init(ctx)

	events, err := QueryEventsByRecordId(ctx, 1)
	if err != nil {
		logger.Info(err.Error())
	}

	logger.Info(events)

}

func TestWrite(t *testing.T) {
	ctx := context.Background()
	Init(ctx)

	e1 := &Event{
		RecordId:       4,
		EventType:      1,
		EventTime:      12,
		MousePos:       "21",
		MouseClickType: 2,
		MouseClickBtn:  1,
		MouseMoveType:  2,
		KeyClickType:   1,
		KeyCode:        "d",
		ComponentName:  "dsa",
		ComponentType:  1,
		ComponentExtra: "dsa",
		AppId:          2,
		Day:            "2022-02-01",
	}
	e2 := &Event{
		RecordId:       5,
		EventType:      1,
		EventTime:      12,
		MousePos:       "21",
		MouseClickType: 2,
		MouseClickBtn:  1,
		MouseMoveType:  2,
		KeyClickType:   1,
		KeyCode:        "d",
		ComponentName:  "dsa",
		ComponentType:  1,
		ComponentExtra: "dsa",
		AppId:          2,
		Day:            "2022-02-01",
	}

	_ = WriteEvents(ctx, []string{e1.String(), e2.String()})
}

func TestWriteSQL(t *testing.T) {
	writeSQL(1, "dsadasd")
}

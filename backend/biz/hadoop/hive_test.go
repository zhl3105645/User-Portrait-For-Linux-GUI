package hadoop

import (
	"backend/gohive"
	"context"
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

	var userId int64
	var beginTime int64
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
			&userId, &beginTime, &eventType,
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

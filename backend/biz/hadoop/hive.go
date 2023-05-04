package hadoop

import (
	"backend/gohive"
	"context"
	"log"
)

var HiveConnection *gohive.Connection

type Event struct {
	UserId         int64
	BeginTime      int64
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

func Init(ctx context.Context) {
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

	HiveConnection = connection
}

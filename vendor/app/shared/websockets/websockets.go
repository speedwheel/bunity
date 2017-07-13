package websockets

import (
	"github.com/kataras/iris/websocket"
)
var ws websocket.Server

func WebsocketInit() *websocket.Server {
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	
	return ws
}
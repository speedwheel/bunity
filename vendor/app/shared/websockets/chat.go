package websockets

import (
	"fmt"
	"github.com/kataras/iris/websocket"
)

func UserChat(c websocket.Connection) {
	c.On("chat", func(msg string) {
		// Print the message to the console, c.Context() is the iris's http context.
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		// Write message back to the client message owner:
		// c.Emit("chat", msg)
		c.To(websocket.Broadcast).Emit("chat", msg)
	})
}
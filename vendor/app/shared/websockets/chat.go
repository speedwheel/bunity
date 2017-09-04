package websockets

import (
	"fmt"
	"github.com/kataras/iris/websocket"
	"sync"
	"app/model"
	ses "app/shared/session"
	//"time"
)

var (
	Conn = make(map[websocket.Connection]string)
	mutex = new(sync.Mutex)
)

func UserChat(c websocket.Connection) {
	ctx := c.Context()
	session := ses.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	mutex.Lock()
	Conn[c] = userSession.Id.Hex()
	mutex.Unlock()
	fmt.Println(Conn)
	c.On("like", func(msg string) {
		for k, id := range Conn {
			if id == msg {
				k.Emit("like", "someone liked")
			}
		}
		// Print the message to the console, c.Context() is the iris's http context.
		//fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		// Write message back to the client message owner:
		// c.Emit("chat", msg)
		//c.Emit("chat", msg)
		//fmt.Println(c.ID())
		//c.To(c.ID()).Emit("chat", "Cli")
	})
	
	c.OnDisconnect(func() {
		mutex.Lock()
		delete(Conn, c)
		mutex.Unlock()
		fmt.Printf("\nConnection with ID: %s has been disconnected!\n", c.ID())
	})
	
	/*var delay = 1 * time.Second
	go func() {
		i := 0
		for {
			mutex.Lock()
			broadcast(Conn, fmt.Sprintf("aaaa %d\n", i))
			mutex.Unlock()
			time.Sleep(delay)
			i++
		}
	}()

	go func() {
		i := 0
		for {
			mutex.Lock()
			broadcast(Conn, fmt.Sprintf("aaaa2 %d\n", i))
			mutex.Unlock()
			time.Sleep(delay)
			i++
		}
	}()*/
}

func broadcast(Conn map[websocket.Connection]bool, message string) {
	for k := range Conn {
		fmt.Println(k.ID())
		k.Emit("chat", message)
	}
}
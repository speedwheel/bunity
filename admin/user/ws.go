package user

import(
	"github.com/kataras/iris/websocket"
	//"app/shared/websockets"
	"fmt"
	"sync"
	ses "app/shared/session"
)

var (
	Conn = make(map[websocket.Connection]string)
	mutex = new(sync.Mutex)
)

type Ws struct {
	Controller
}

/*func NewWsSource() *Ws {
	return &Ws{
		Conn: websockets.WebsocketInit(),
	}
}*/

func BusinessChatNotif(c websocket.Connection) {
	ctx := c.Context()
	session := ses.Sessions.Start(ctx)
	adminID := session.GetString("adminID")
	mutex.Lock()
	Conn[c] = adminID
	mutex.Unlock()
	fmt.Println(Conn)
	
	c.On("newBusinessChat", func(commentID string) {
	
		source := NewDataSource()
		comment := source.GetCommentsById(commentID)
		
		for k, id := range Conn {
			if id != adminID  {
				k.Emit("newBusinessChat", comment)
				
			}
		}
	})
	c.On("refreshCommBiz", func(businessID string) {
		source := NewDataSource()
		commentsBusiness := source.GetCommentsByBusiness(businessID)
		for k, id := range Conn {
			if id != adminID  {
				k.Emit("refreshCommBiz", commentsBusiness)
			}
		}
		
	})
	
	c.OnDisconnect(func() {
		mutex.Lock()
		delete(Conn, c)
		mutex.Unlock()
		fmt.Printf("\nConnection with ID: %s has been disconnected!\n", c.ID())
	})
}

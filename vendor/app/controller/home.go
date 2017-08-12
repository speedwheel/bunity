package controller

import(
	"github.com/kataras/iris/context"
	//"github.com/messagebird/go-rest-api"
	//"fmt"
)

func Home (ctx context.Context) {
	/*client := messagebird.New("qV8HkQdNlDD0UDa9Z3mrXnlXK")
	params := &messagebird.MessageParams{Reference: "MyReference"}

	message, _ := client.NewMessage(
	  "Edward",
	  []string{"+447884204004"},
	  "Test message really from KBN",
	  params)
	  
	 fmt.Println(message)*/
	// fmt.Println(err)
	//ctx.Gzip(true)
	ctx.View("home.html")
}

func Page (ctx context.Context) {
	
	ctx.HTML("<h1> Kazeli Page </h1><a href='/'>home</a>")
}
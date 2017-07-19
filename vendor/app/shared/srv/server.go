package srv

import(
	stdContext "context"
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"app/route"
	//"app/shared/general"
	"app/shared/websockets"
	"app/config"
	//"app/model"
	"app/model"
	"html/template"
	"strings"
	"app/shared/db"
)

var (
	app *iris.Application
)

func KazeliApp() *iris.Application {
	app = iris.New()
	//app.Use(model.IsAuth)
	app := iris.New()
	
	app.Use(func(ctx context.Context) {
		session := db.Sessions.Start(ctx)
	     //ctx.Gzip(true)
		ctx.ViewData("auth", session.Get("userAuth"))
		userSession := model.User{}
		if (session.Get("user") != nil) {
			userSession = session.Get("user").(model.User)
		}
		ctx.ViewData("userSession", userSession)
		ctx.Next()
	})
	ws := websockets.WebsocketInit()
	ws.OnConnection(websockets.UserChat)
	app.Get("/userchat", ws.Handler())
	tmpl := iris.HTML(config.GetAppPath()+"templates", ".html")
	tmpl.Reload(true)
	app.RegisterView(tmpl.Layout("layouts/default.html"))
	app.StaticWeb("/static", config.GetAppPath()+"resources")
	tmpl.AddFunc("getRatio", func(val string) string {
		newVal := val[len(val)-5:len(val)-4]
		if newVal == "1" {
			return "landscape"
		} else {
			return "portrait"
		}
	})
	
	tmpl.AddFunc("raw", func(char string) template.HTML {
		 return template.HTML(strings.Replace(char,"<br>","",-1))
	})
	
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx context.Context) {
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}

		ctx.Writef("(Unexpected) internal server error")
	})
	
	
	route.Routes(app)
	
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		// close all hosts
		app.Shutdown(ctx)
	})
	return app
}
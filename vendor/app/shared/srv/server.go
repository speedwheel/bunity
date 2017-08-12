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
	"encoding/gob"
	"gopkg.in/mgo.v2/bson"
)

var (
	app *iris.Application
)

func KazeliApp() *iris.Application {
	gob.Register(model.User{})
	app = iris.New()
	//app.Use(model.IsAuth)
	
	 app.Use(func(ctx context.Context) {
		session := db.Sessions.Start(ctx)
	    ctx.Gzip(true)
		auth := false
		if session.Get("userAuth") != nil {
			auth = true
		}
		ctx.ViewData("auth", auth)
		ctx.Values().Set("auth", auth)
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
	tmpl := iris.HTML("./templates", ".html")/*.Binary(general.Asset, general.AssetNames)*/.Layout("layouts/default.html")
	tmpl.Reload(false)
	app.RegisterView(tmpl)
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
	
	tmpl.AddFunc("add", func(x, y int) int {
		return x + y
	})
	
	tmpl.AddFunc("sub", func(x, y int) int {
		return x - y
	})
	
	tmpl.AddFunc("divisible83", func(x int) bool {
		if x > 0 {
			if x % 83 == 1 {
				return true
			}
			return false
		}
		return false
	})
	
	tmpl.AddFunc("countAr", func(x []bson.ObjectId) int {
		return len(x)
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
	
	app.OnErrorCode(404, func(ctx context.Context) {
	 	ctx.Writef("My Custom 404 error page ")
	 })
	 
	
	
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		// close all hosts
		app.Shutdown(ctx)
	})
	
	
	return app
}
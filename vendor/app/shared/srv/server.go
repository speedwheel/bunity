package srv

import(
	stdContext "context"
	"time"
	"github.com/kataras/iris"
	"app/route"
	//"app/shared/general"
	"app/shared/websockets"
	"github.com/kataras/iris/websocket"
	"app/config"
	//"app/model"
	"app/model"
	"html/template"
	"strings"
	ses "app/shared/session"
	"encoding/gob"
	"gopkg.in/mgo.v2/bson"
	"github.com/iris-contrib/middleware/cors"
	"github.com/speedwheel/bunity/admin"
)


var (
	app *iris.Application
)

func init() {
	gob.Register(model.User{})
	gob.Register(model.Business{})
}

func KazeliApp() *iris.Application {
	app = iris.New()
	
	app.WrapRouter(cors.WrapNext(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}))
	
	app.Use(func(ctx iris.Context) {
		session := ses.Sessions.Start(ctx)
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
		ctx.ViewData("domain", "http://bunity.com:8080")
		ctx.ViewData("port", "8080")

		if ctx.Path() != "/userchat" && ctx.Path() != "/iris-ws.js" && ctx.Path() != "/notifications" {
			ctx.Gzip(true)
		}
		ctx.Next()
	})	

	route.Routes(app)
	app.StaticWeb("/static", config.GetAppPath()+"resources")
	
	//app.Use(model.IsAuth)
			
	ws := websockets.WebsocketInit()
	
	ws.OnConnection(websockets.UserChat)
	app.Get("/userchat", ws.Handler())
	app.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})
	
	tmpl := iris.HTML(config.GetAppPath()+"templates", ".html")/*.Binary(general.Asset, general.AssetNames)*/.Layout("layouts/default.html")
	tmpl.Reload(true)
	app.RegisterView(tmpl)
	
	
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
	
	tmpl.AddFunc("inSlice", func(u []admin.AdminSubItem, p string) bool {
		for _, i := range u {
        if i.Url == p {
            return true
        }
    }
    return false
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
	
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}

		ctx.Writef("(Unexpected) internal server error")
	})
	
	 
	
	
	 app.OnErrorCode(404, func(ctx iris.Context) {
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
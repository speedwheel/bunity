package route

import(
	"github.com/kataras/iris"
	"app/controller"
	"app/model"
	"app/shared/general"
	"github.com/speedwheel/bunity/admin/user"
	"app/shared/session"
	"github.com/speedwheel/bunity/admin"
	"app/config"
	"github.com/kataras/iris/websocket"
	"app/shared/websockets"

)

func Routes(app *iris.Application) {
	app.Get("/", controller.Home)
	app.Get("/page", controller.Page)
	app.Get("/register", controller.Register)
	app.Post("/register", controller.RegisterUser)
	app.Get("/login", controller.Login)
	app.Post("/login", controller.LoginPost, model.SetUserSession)
	app.Post("/signupsocial", controller.SingupSocial, model.SetUserSession)
	app.Get("/profile", controller.Profile, general.Auth)
	
	app.Get("/{businessID:string}", controller.BusinessProfilePage)
	app.Get("/{businessID:string}/maps", controller.BusinessProfileMaps)
	app.Get("/{businessID:string}/webresults", controller.BusinessProfileWeb)
	app.Get("/{businessID:string}/webresults/{pageCount:int}", controller.BusinessProfileWeb)
	app.Get("/{businessID:string}/internalresults", controller.BusinessProfileInternal)
	app.Get("/{businessID:string}/internalresults/{pageCount:int}", controller.BusinessProfileInternal)

	app.Get("/users/set_password/:eal_exp", controller.CheckForgotToken)
	app.Post("/users/set_password/:eal_exp", controller.UpdatePassword)
	
	app.Post("/forgot_password", controller.ForgotPassword)
	app.Get("/logout", general.Logout)
	
	app.Post("/likes/{businessID:string}", controller.BusinessLike)
	
	app.Post("/livesearch", controller.LiveSearch)
	
	app.Any("/search/business", controller.BusinessSearch)
	app.Get("/search/business/{pageCount:int}", controller.BusinessSearch)
	
	//users edit
	users := app.Party("/users", model.IsAuthRedirect)
	
	users.Get("/edit", controller.UsersEdit)
	users.Post("/edit", controller.UsersEditUpdate) 
	users.Get("/payment_methods", controller.PaymentMethods)
	
	
	
	businesses := app.Party("/businesses", model.IsAuthRedirect)
	businesses.Get("/", controller.BusinessList)
	businesses.Any("/add", controller.BusinessAddStep1)
	businesses.Any("/add/step2", controller.BusinessAddStep2)
	businesses.Any("/add/step1/{businessID:string}", controller.BusinessAddStep1)
	businesses.Any("/add/step2/{businessID:string}", controller.BusinessAddStep2)
	businesses.Any("/add/step3/{businessID:string}", controller.BusinessAddStep3)
	businesses.Any("/add/step4/{businessID:string}", controller.BusinessAddStep4)
	businesses.Any("/add/step5/{businessID:string}", controller.BusinessAddStep5)
	businesses.Any("/add/step6/{businessID:string}", controller.BusinessAddStep6)
	businesses.Any("/add/step7/{businessID:string}", controller.BusinessAddStep7)
	
	businesses.Any("/delete/{businessID:string)}", controller.BusinessDelete)
	businesses.Post("/trackEvents", controller.BusinessEventsTracker, controller.BusinessAddFinish)
	businesses.Post("/addfiles", controller.UploadFiles)
	businesses.Post("/deletefile", controller.DeleteFile)
	businesses.Post("/sendsms", controller.SendSms)
	businesses.Post("/verifycode", controller.VerifyCode)
	
	businesses.Get("/categories", controller.BusinessAllCategPage)
	businesses.Get("/categories/{businessSlug:string}", controller.BussinessByCategPage)
	
	
	businesses.Post("/updatephotos", controller.UpdatePhotos)
	
	
	office := app.Party("office.", func(ctx iris.Context ) {
		adminMenu := admin.NewAdminMenu()
		ctx.ViewData("AdminMenu", adminMenu)
		ctx.Next()
	}).Layout("admin/layouts/default.html")
	{
		office.StaticWeb("/adminstatic", config.GetAppPath()+"admin/resources")
		office.Any("/iris-ws.js", func(ctx iris.Context) {
			ctx.Write(websocket.ClientSource)
		})
		users := user.NewDataSource()
		ws := websockets.WebsocketInit()
		office.Get("/notifications", ws.Handler())
		ws.OnConnection(user.BusinessChatNotif)
		
		office.Controller("/", new(user.Controller), session.Sessions, users)
	}
}
package controller

import(
	"github.com/kataras/iris/context"
	"app/model"
)

func Register(ctx context.Context) {
	ctx.View("register.html")
}

func RegisterUser(ctx context.Context) {
	usr := model.User{}
	usr.Firstname = ctx.FormValue("firstname")
	usr.Lastname = ctx.FormValue("lastname")
	usr.Email = ctx.FormValue("email")
	usr.Account.Password = ctx.FormValue("password")
	model.RegisterUser(&usr)
	ctx.View("register.html")
}
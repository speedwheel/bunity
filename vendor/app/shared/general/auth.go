package general

import (
	"github.com/kataras/iris/context"
	"app/shared/db"
)

func Auth(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	auth := session.Get("userAuth")
	if auth == false || auth == nil {
		ctx.Redirect("/")
	}
}

func Logout(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	session.Clear()
	db.Sessions.Destroy(ctx)
	ctx.Redirect("/")
}
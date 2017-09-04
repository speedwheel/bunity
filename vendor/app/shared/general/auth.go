package general

import (
	"github.com/kataras/iris/context"
	ses "app/shared/session"
)

func Auth(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	auth := session.Get("userAuth")
	if auth == false || auth == nil {
		ctx.Redirect("/")
	}
}

func Logout(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	session.Clear()
	ses.Sessions.Destroy(ctx)
	ctx.Redirect("/")
}
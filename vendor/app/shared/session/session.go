package session

import (
	"gopkg.in/kataras/iris.v6"

)

func Da(ctx *iris.Context) {

	ctx.Set("da", "da")
	ctx.Next()
}
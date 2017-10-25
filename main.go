package main

import (
	"app/shared/srv"
	"github.com/kataras/iris"
)

func main() {
	app := srv.KazeliApp()
	/*if err := app.Run(iris.AutoTLS("2message.com:443")); err != nil {
		panic(err)
	}*/

	if err := app.Run(iris.Addr("bunity.com:8080")); err != nil {
		panic(err)
	}
}
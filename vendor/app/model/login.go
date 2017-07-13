package model

import (
	"app/shared/db"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"app/config"
	"github.com/kataras/iris/context"
)



func CheckLogin(_email string, _pass string, ctx context.Context) {
	cfg := config.Init()
	SecretKey := cfg.User.SecretKey
	result := User{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Find(bson.M{"email": _email}).One(&result); err != nil {
		panic(err.Error())
	}
	Db.Close()
	hash := []byte(result.Account.Password)
	pass := []byte(_pass + SecretKey)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	auth := true
	if err == nil {
		auth = false
	}
	ctx.Values().Set("auth", auth)
	ctx.Values().Set("user", result)
	ctx.Next()
}

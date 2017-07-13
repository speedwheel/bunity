package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"app/shared/db"
	"golang.org/x/crypto/bcrypt"
	"app/config"
)

func RegisterUser(usr *User) {
	cfg := config.Init()
	SecretKey := cfg.User.SecretKey
	usrDB := User{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Find(bson.M{"account.email": usr.Email}).One(&usrDB); err != nil {
		usr.Id = bson.NewObjectId()
		usr.Account.Kind = "internal"
		password := []byte(usr.Account.Password + SecretKey)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			//panic(err.Error())
		}
		usr.Account.Password = string(hashedPassword)
		err = c.Insert(&usr)
		if err != nil {
			//panic(err)
		}
		fmt.Println(usr)
	}
	Db.Close()
	
	
}
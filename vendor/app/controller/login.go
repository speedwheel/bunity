package controller

import(
	"github.com/kataras/iris/context"
	"app/model"
	"app/shared/db"
	ses "app/shared/session"
	"gopkg.in/mgo.v2/bson"
	"app/shared/social/fb"
	"app/shared/mail"
	"github.com/dchest/passwordreset"
	"time"
	"log"
	"golang.org/x/crypto/bcrypt"
	"app/config"
)

var (
	usr *model.User
	passwordResetToken = []byte("secret key")
	auth = false
)



func Login(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	auth := session.Get("userAuth")
	if auth == true {
		ctx.Redirect("/profile") 
	}
	ctx.View("login.html")	
}

func LoginPost(ctx context.Context) {
	_email := string(ctx.FormValue("email"))
	_pass := string(ctx.FormValue("password"))
	cfg := config.Init()
	SecretKey := cfg.User.SecretKey
	result := model.User{}
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
	if err != nil {
		auth = false
		result = model.User{}
	}
	ctx.Values().Set("auth", auth)
	ctx.Values().Set("user", result)
	ctx.Next()
}

func SingupSocial(ctx context.Context) {
	kind := ctx.FormValue("kind")
	action := string(ctx.FormValue("action"))
	usr = &model.User{}
	
	if(kind == "fb") {
		token := ctx.FormValue("token")
		fbApi := fb.Fb{}
		fbApi.Init(token)
		fbUser := fbApi.GetUser()
		usr.Id = bson.NewObjectId()
		usr.Account.Kind = kind
		fbUser.DecodeField("id",&usr.Account.Uid)
		fbUser.DecodeField("picture.data.url",&usr.Image)
		fbUser.Decode(&usr)
	}
	
	if(kind == "google") {
		usr.Firstname = ctx.FormValue("firstname")
		usr.Lastname = ctx.FormValue("lastname")
		usr.Email = ctx.FormValue("email")
		usr.Image = ctx.FormValue("image")
		usr.Account.Kind = ctx.FormValue("kind")
		usr.Account.Uid = ctx.FormValue("uid")
		usr.Id = bson.NewObjectId()
	}
	usrDB := model.User{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Find(bson.M{"$or": []bson.M{bson.M{"email": usr.Email},bson.M{"account.uid": usr.Account.Uid}}}).One(&usrDB); err != nil {
		
		if action == "signup" {
			if err := c.Insert(&usr); err != nil {
				panic(err.Error())
			}
			if err := c.Find(bson.M{"account.uid": usr.Account.Uid}).One(&usrDB); err != nil {
				panic(err.Error())
			}
			auth = true
		}
	}
	if action == "login" {
		if err := c.Find(bson.M{"account.uid": usr.Account.Uid}).One(&usrDB); err != nil {

		}
		auth = true
	}
	Db.Close()

	ctx.Values().Set("auth", auth)
	ctx.Values().Set("user", usrDB)
	ctx.Next()
	ctx.JSON(map[string]*bool{"userAuth": &auth})
}

func ForgotPassword(ctx context.Context) {
	usr = &model.User{}
	Db := db.MgoDb{}
	Db.Init()
	email := ctx.FormValue("email")
	
	c := Db.C("users")
	
	if err := c.Find(bson.M{"account.kind": "internal", "email": email}).Select(bson.M{"_id": 0, "email": 1}).One(&usr); err != nil {
		log.Printf(err.Error())

	}

	pwdval, _ := getPasswordHash(usr.Email)
	
	resetPasswordToken := passwordreset.NewToken(usr.Email, 1 * time.Hour, pwdval, passwordResetToken)
	
	if err := c.Update(bson.M{"account.kind": "internal", "email": email}, bson.M{"$set": bson.M{"account.resetPasswordToken": resetPasswordToken}}); err != nil {
		log.Printf(err.Error())

	}
	Db.Close()
	
	ctx.Redirect("/")
	mailApi := mail.Mail{}
	mailApi.Init()
	to := []string{"edi.ultras@gmail.com"}
	content := `<a href="https://bunity.com:8080/users/set_password/`+resetPasswordToken+`">sdafdas</a>`
	mailApi.Send(to, content)
}

func CheckForgotToken(ctx context.Context) {
	token := ctx.Values().Get("eal_exp").(string)
	_, err := passwordreset.VerifyToken(token, getPasswordHash, passwordResetToken);
	if err != nil {
		// verification failed, don't allow password reset
		return
	}
	ctx.ViewData("token", token)
	ctx.View("forgot_password.html")
}

func UpdatePassword(ctx context.Context) {
	cfg := config.Init()
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	token := ctx.Values().Get("eal_exp").(string)

	_password := string(ctx.FormValue("password"))

	SecretKey := cfg.User.SecretKey
	password := []byte(_password + SecretKey)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf(err.Error())
	}
	if err := c.Update(bson.M{"account.kind": "internal", "account.resetPasswordToken": token}, bson.M{"$set": bson.M{"account.password": hashedPassword}}); err != nil {
		log.Printf(err.Error())
		return
	}
	ctx.Redirect("/")
}

func getPasswordHash(email string) ([]byte, error) {
	var err error
	usr2 := model.User{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err = c.Find(bson.M{"account.kind": "internal", "email": email}).Select(bson.M{"_id": 0, "account.password": 1}).One(&usr2); err != nil {
		panic(err.Error())
		return nil, err
	}
	pw := []byte(usr2.Account.Password)
	return pw, nil
}
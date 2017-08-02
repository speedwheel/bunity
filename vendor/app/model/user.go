package model

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/kataras/iris/context"
	"app/shared/db"
	"log"
	"html/template"
)


type (
	User struct {
		Id     bson.ObjectId `json:"id" bson:"_id"  form:"-" facebook:"-"`
		Firstname   string `json:"firstname" bson:"firstname"  form:"firstname" facebook:"first_name"`
		Lastname   string `json:"lastname" bson:"lastname"  form:"lastname" facebook:"last_name"`
		Email  string `json:"email" bson:"email"  form:"email" facebook:"email"`
		Image string `json:"image" bson:"image"  form:"image" facebook:"url"`
		Liked []bson.ObjectId `json:"liked" bson:"liked" form:"liked,omitempty" facebook:"-"`
		Account Account
		Businesses []bson.ObjectId `json:"businesses" bson:"businesses" form:"businesses,omitempty" facebook:"-"`
		
	}
	
	Account struct {
		Kind string `json:"kind" bson:"kind" form:"kind,omitempty" facebook:"-"`
		Uid string `json:"uid" bson:"uid,omitempty"  form:"-" facebook:""`
		Password string `json:"password" bson:"password,omitempty" form:"password,omitempty" facebook:"-"`
		ResetPasswordToken string `json:"resetPasswordToken" bson:"resetPasswordToken" form:"resetPasswordToken,omitempty" facebook:"-"`
	}
	

	
	Business struct {
		Id  bson.ObjectId `json:"id" bson:"_id"  form:"-" facebook:"-"`
		UserId  bson.ObjectId `json:"user_id" bson:"user_id"  form:"-" facebook:"-"`
		Name string `json:"name" bson:"name" form:"name,omitempty" facebook:"-"`
		Slug string `json:"slug" bson:"slug" form:"slug,omitempty" facebook:"-"`
		Phone string `json:"phone" bson:"phone" form:"phone,omitempty" facebook:"-"`
		Address string `json:"address" bson:"address" form:"address,omitempty" facebook:"-"`
		Address2 string `json:"address2" bson:"address2" form:"address2,omitempty" facebook:"-"`
		State string `json:"state" bson:"state" form:"state,omitempty" facebook:"-"`
		City string `json:"city" bson:"city" form:"city,omitempty" facebook:"-"`
		PostalCode string `json:"postalcode" bson:"postalcode" form:"postalcode,omitempty" facebook:"-"`
		Website string `json:"IsAuth" bson:"website" form:"website,omitempty" facebook:"-"`
		Area string `json:"area" bson:"area" form:"area,omitempty" facebook:"-"`
		Country string `json:"country" bson:"country" form:"country,omitempty" facebook:"-"`
		Industry string `json:"industry" bson:"industry" form:"industry,omitempty" facebook:"-"`
		YearsBusiness string `json:"yearsBusiness" bson:"yearsBusiness" form:"yearsBusiness,omitempty" facebook:"-"`
		NumberEmployees string `json:"numberEmployees" bson:"numberEmployees" form:"numberEmployees,omitempty" facebook:"-"`
		SizeBusiness string `json:"sizeBusiness" bson:"sizeBusiness" form:"sizeBusiness,omitempty" facebook:"-"`
		RelationshipBusiness string `json:"relationshipBusiness" bson:"relationshipBusiness" form:"relationshipBusiness,omitempty" facebook:"-"`
		HowYourHear string `json:"howYourHear" bson:"howYourHear" form:"howYourHear,omitempty" facebook:"-"`
		Description template.HTML `json:"description" bson:"description" form:"description,omitempty" facebook:"-"`
		Gallery []string `json:"gallery" bson:"gallery" form:"gallery,omitempty" facebook:"-"`
		Profile []string `json:"profile" bson:"profile" form:"profile,omitempty" facebook:"-"`
		Cover []string `json:"cover" bson:"cover" form:"profile,omitempty" facebook:"-"`
		SmsCode string `json:"smsCode" bson:"smsCode" form:"smsCode,omitempty" facebook:"-"`
		Likes []bson.ObjectId `json:"likes" bson:"likes" form:"likes,omitempty" facebook:"-"`
		Verified uint8 `json:"check" bson:"check" form:"check,omitempty" facebook:"-"`
		Premium uint8 `json:"pro" bson:"pro" form:"pro,omitempty" facebook:"-"`
		Social Social
		Map Map
	}
	
	Social struct {
		Facebook string `json:"facebook" bson:"facebook" form:"facebook,omitempty" facebook:"-"`
		Google string `json:"google" bson:"google" form:"google,omitempty" facebook:"-"`
		Instagram string `json:"instagram" bson:"instagram" form:"instagram,omitempty" facebook:"-"`
		Youtube string `json:"youtube" bson:"youtube" form:"youtube,omitempty" facebook:"-"`
		Pinterest string `json:"pinterest" bson:"pinterest" form:"pinterest,omitempty" facebook:"-"`
		Linkedin string `json:"linkedin" bson:"linkedin" form:"linkedin,omitempty" facebook:"-"`
		Twitter string `json:"twitter" bson:"twitter" form:"twitter,omitempty" facebook:"-"`
	}
	
	Map struct {
		Lat float64 `json:"lat" bson:"lat" form:"lat,omitempty" facebook:"-"`
		Lng float64 `json:"lng" bson:"lng" form:"lng,omitempty" facebook:"-"`
	}
)

func SetUserSession(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	//usrInterface := ctx.Get("user")
	//usr := usrInterface.(*User)
	session.Set("userAuth", ctx.Values().Get("auth"))
	session.Set("user", ctx.Values().Get("user"))
}


func IsAuthRedirect(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	_, err := session.Get("userAuth").(bool)
	if !err  {
		ctx.Redirect("/")
	}
	ctx.Next()
}

func GetUserByID(userID bson.ObjectId) User {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	user := User{}
	if err := c.Find(bson.M{"_id": userID}).One(&user); err != nil {
		panic(err)
	}
	Db.Close()
	return user
}

func GetBusinessByIDandUser(businessID bson.ObjectId, userID bson.ObjectId) Business{
	
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	business := Business{}
	//if err := c.Find(bson.M{"_id": userID, "businesses": bson.M{ "$elemMatch": bson.M{"_id":businessID}}}).Select(bson.M{"_id":0, "businesses.$": 1}).One(&business); err != nil {
	if err := c.Find(bson.M{"user_id": userID, "_id": businessID}).One(&business); err != nil {	
	/*oe := bson.M{
        "$match" :bson.M {"_id": userID, "businesses._id": businessID},
	}
	oa := bson.M{
        "$project" :bson.M {"_id": 0, "businesses": 1},
	}
	oc := bson.M{
        "$unwind":"$businesses",
	}
	pipe := c.Pipe([]bson.M{oe,oa, oc})*/
	
	//if err := pipe.One(&business); err != nil {	
		log.Printf(err.Error())
	}
	Db.Close()
	return business
}

func GetBusinessByID(businessID bson.ObjectId) (error, Business){
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	business := Business{}
	
	err := c.Find(bson.M{"_id": businessID}).One(&business)
	Db.Close()
	return err, business
}

func GetAllBusinessByUser(userID bson.ObjectId) []Business{
	
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	business := []Business{}
	if err := c.Find(bson.M{"user_id": userID}).All(&business); err != nil {	

		log.Printf(err.Error())
	}
	Db.Close()
	return business
}
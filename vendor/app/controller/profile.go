package controller

import(
	"github.com/kataras/iris/context"
	"app/model"
	"app/shared/db"
	"gopkg.in/mgo.v2/bson"
	"log"
	"github.com/logpacker/PayPal-Go-SDK"
	"os"
	"fmt"
	"app/config"
	//"io"
	//"io/ioutil"
	"path/filepath"
	"math/rand"
	"time"
	"github.com/nfnt/resize"
	"bytes"
    "image"
	"image/jpeg"
    "image/png"
	"strconv"
	"github.com/messagebird/go-rest-api"
	"html/template"
	"app/shared/general"
	"math"
	//"github.com/patrickmn/go-cache"
	"github.com/kr/pretty"

)

type FormError struct {
	Class string
	Message string
}

type LiveResults struct {
	Name string
	Image[] string
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const smsChars = "1234567890"

var (
	countries = []string {"Afghanistan","Albania","Algeria","Andorra","Angola","Anguilla","Antigua Barbuda","Argentina","Armenia","Aruba","Australia","Austria","Azerbaijan","Bahamas","Bahrain","Bangladesh","Barbados","Belarus","Belgium","Belize","Benin","Bermuda","Bhutan","Bolivia","Bosnia Herzegovina","Botswana","Brazil","British Virgin Islands","Brunei","Bulgaria","Burkina Faso","Burundi","Cambodia","Cameroon","Canada","Cape Verde","Cayman Islands","Chad","Chile","China","Colombia","Congo","Cook Islands","Costa Rica","Cote D Ivoire","Croatia","Cruise Ship","Cuba","Cyprus","Czech Republic","Denmark","Djibouti","Dominica","Dominican Republic","Ecuador","Egypt","El Salvador","Equatorial Guinea","Estonia","Ethiopia","Falkland Islands","Faroe Islands","Fiji","Finland","France","French Polynesia","French West Indies","Gabon","Gambia","Georgia","Germany","Ghana","Gibraltar","Greece","Greenland","Grenada","Guam","Guatemala","Guernsey","Guinea","Guinea Bissau","Guyana","Haiti","Honduras","Hong Kong","Hungary","Iceland","India","Indonesia","Iran","Iraq","Ireland","Isle of Man","Israel","Italy","Jamaica","Japan","Jersey","Jordan","Kazakhstan","Kenya","Kuwait","Kyrgyz Republic","Laos","Latvia","Lebanon","Lesotho","Liberia","Libya","Liechtenstein","Lithuania","Luxembourg","Macau","Macedonia","Madagascar","Malawi","Malaysia","Maldives","Mali","Malta","Mauritania","Mauritius","Mexico","Moldova","Monaco","Mongolia","Montenegro","Montserrat","Morocco","Mozambique","Namibia","Nepal","Netherlands","Netherlands Antilles","New Caledonia","New Zealand","Nicaragua","Niger","Nigeria","Norway","Oman","Pakistan","Palestine","Panama","Papua New Guinea","Paraguay","Peru","Philippines","Poland","Portugal","Puerto Rico","Qatar","Reunion","Romania","Russia","Rwanda","Saint Pierre Miquelon","Samoa","San Marino","Satellite","Saudi Arabia","Senegal","Serbia","Seychelles","Sierra Leone","Singapore","Slovakia","Slovenia","South Africa","South Korea","Spain","Sri Lanka","St Kitts Nevis","St Lucia","St Vincent","St. Lucia","Sudan","Suriname","Swaziland","Sweden","Switzerland","Syria","Taiwan","Tajikistan","Tanzania","Thailand","Timor L'Este","Togo","Tonga","Trinidad Tobago","Tunisia","Turkey","Turkmenistan","Turks Caicos","Uganda","Ukraine","United Arab Emirates","United Kingdom","United States","United States Minor Outlying Islands","Uruguay","Uzbekistan","Venezuela","Vietnam","Virgin Islands (US)","Yemen","Zambia","Zimbabwe"}
	statesUSA = []string {"Alabama","Alaska","Arizona","Arkansas","California","Colorado","Connecticut","Delaware","District of Columbia","Florida","Georgia","Hawaii","Idaho","Illinois","Indiana","Iowa","Kansas","Kentucky","Louisiana","Maine","Maryland","Massachusetts","Michigan","Minnesota","Mississippi","Missouri","Montana","Nebraska","Nevada","New Hampshire","New Jersey","New Mexico","New York","North Carolina","North Dakota","Ohio","Oklahoma","Oregon","Pennsylvania","Rhode Island","South Carolina","South Dakota","Tennessee","Texas","Utah","Vermont","Virginia","Washington","West Virginia","Wisconsin","Wyoming"}
	statesCanada = []string {"British Columbia","Ontario","Newfoundland and Labrador","Nova Scotia", "Prince Edward Island", "New Brunswick", "Quebec", "Manitoba", "Saskatchewan", "Alberta", "Northwest Territories", "Nunavut","Yukon Territory"}
	statesAustralia = []string {"New South Wales","Victoria","Queensland","Tasmania","South Australia","Western Australia","Northern Territory","Australian Capital Terrirory"}
	industry = []string {"Accounting","Advertising","Automotive","Computers","Construction","Consulting","Dental","Education","Entertainment","Entrepreneur","Financial","Health Care","Internet","Law","Manufacturing","Marketing","Medical","Printing","Publishing","Real Estate","Restaurant","Retail","Sales","Service","Telecommunications","Travel","Wholesale"}
	yearsBusiness = []string {"Start-Up","1 - 5 Years","6 - 10 Years","11 + Years"}
	nrEmployees = []string {"1-5","6-10","11-20","21-50","51-100","101-250","251+"}
	sizeBusiness = []string{"$100k - $1MM","$101MM +","$26MM - $100MM","$2MM - $5MM","$6MM - $25MM","Just Starting","Less than $100k"}
	relationshipBusiness = []string{"I'm the owner of this company.", "I work for this company.","I don't work here, but I'm acting on behalf of this company.","I'm a user of Zaphiri improving the business listing."}
	howYouHear = []string{"Business Breakthroughs International","Business Coach/Consultant","Business Mastery Event","Chet Holmes","Friend/Associate","Mitch Russo","Radio","Search Engine","Social Networking","Solution Provider","Tony Robbins","Trade Show","Twitter","Web Seminar"}
)



func Profile(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	ctx.Next()
	ctx.ViewData("user", session.Get("user"))
	ctx.View("profile.html")
}

func UsersEdit(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	ctx.ViewData("user", session.Get("user"))
	ctx.View("users_edit.html")
}

func UsersEditUpdate (ctx context.Context) {
	Db := db.MgoDb{}
	Db.Init()
	session := db.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	userSession.Firstname = ctx.FormValue("firstname")
	userSession.Lastname = ctx.FormValue("lastname")
	userSession.Email = ctx.FormValue("email")
	
	c := Db.C("users")
	
	if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$set": bson.M{"firstname": userSession.Firstname, "lastname": userSession.Lastname, "email": userSession.Email}}); err != nil {
		log.Printf(err.Error())
		return
	}
	Db.Close()
	session.Set("user", userSession)
	ctx.ViewData("user", userSession)
	
	ctx.View("users_edit.html")
}

func PaymentMethods(ctx context.Context) {
	c, err := paypalsdk.NewClient("AcV2C_ISVEx8aFrB1t6SpK3s7DSColZIHXcx8IapBO0dvCeSUk_8bC4S5lhVKVxvDnw7NS8eDNsS_9ge", "EOChceipRsgEVeJTcTK49XoATyptU7MKHR3Vb4M47YlJg0nR1OYhOyqFeNmRHyATnOEH5S5vbuUNDsYv", paypalsdk.APIBaseSandBox)
	c.SetLog(os.Stdout) // Set log to terminal stdout
	c.GetAccessToken()
	if err != nil {
	
	}
	amount := paypalsdk.Amount{
		Total:    "7.00",
		Currency: "USD",
	}
	redirectURI := "http://example.com/redirect-uri"
	cancelURI := "http://example.com/cancel-uri"
	description := "Description for this payment"
	paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)

	c.ExecuteApprovedPayment(paymentResult.ID, "7E7MGXCWTTKK2")
	ctx.View("payment_methods.html")
}

func BusinessList(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	ctx.ViewData("business", userSession.Businesses)
	ctx.View("businesses.html")	
}

func BusinessAddStep1(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	business := model.Business{}
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := session.Get("user").(model.User)
		business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil {
		business = session.Get("businessForm").(model.Business)
	}
	ctx.ViewData("business", business)
	ctx.ViewData("countries", countries)
	ctx.ViewData("statesUSA", statesUSA)
	ctx.ViewData("statesCanada", statesCanada)
	ctx.ViewData("statesAustralia", statesAustralia)
	
	ctx.View("business_add.html")	
}

func BusinessAddStep11(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	business := model.Business{}
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := session.Get("user").(model.User)
		business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil {
		business = session.Get("businessForm").(model.Business)
	}
	ctx.ViewData("business", business)
	ctx.ViewData("countries", countries)
	ctx.ViewData("statesUSA", statesUSA)
	ctx.ViewData("statesCanada", statesCanada)
	ctx.ViewData("statesAustralia", statesAustralia)
	
	ctx.View("business_add.html")	
}

func BusinessAddPost(ctx context.Context) {
	
	/*business := model.User{
		Businesses: []model.Business  {
			{
				bson.NewObjectId(),
				ctx.FormValue("business[name]"),
				ctx.FormValue("description"),
				ctx.FormValue("phone"),
				ctx.FormValue("website"),
				ctx.FormValue("address"),
			},
		},
	}*/
	
	

	//Db := db.MgoDb{}
	//Db.Init()
	//c := Db.C("users")

	//userSession := session.Get("user").(model.User)
	//if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$push": bson.M{"businesses": bson.M{"_id": business.Businesses[0].Id, "name": business.Businesses[0].Name, "description": business.Businesses[0].Description, "phone": business.Businesses[0].Phone, "website": business.Businesses[0].Website, "address": business.Businesses[0].Address}}}); err != nil {
	//	panic(err)
	//}
	//Db.Close()
	
	
}

func BusinessAddStep2(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	business := model.Business{}
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := session.Get("user").(model.User)
		business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil { 
		business = session.Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	
	ctx.ViewData("data", map[string][]string{
		"industries": industry,
		"yearsBusiness": yearsBusiness,
		"nrEmployees": nrEmployees,
		"sizeBusiness": sizeBusiness,
		"relationshipBusiness": relationshipBusiness,
		"howYouHear": howYouHear,
	})
	
	ctx.ViewData("business", business)
	ctx.View("business_add2.html")
}

func BusinessAddStep22(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	business := model.Business{}
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := session.Get("user").(model.User)
		business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil { 
		business = session.Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	
	ctx.ViewData("data", map[string][]string{
		"industries": industry,
		"yearsBusiness": yearsBusiness,
		"nrEmployees": nrEmployees,
		"sizeBusiness": sizeBusiness,
		"relationshipBusiness": relationshipBusiness,
		"howYouHear": howYouHear,
	})
	
	ctx.ViewData("business", business)
	ctx.View("business_add2.html")
}

func BusinessAddStep3(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	business := model.Business{}
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := session.Get("user").(model.User)
		business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	} else if businessSession != nil {
		business = session.Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	
	ctx.ViewData("business", business)
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.View("business_add3.html")
}

func BusinessAddStep4(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	business := model.Business{}
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := session.Get("user").(model.User)
		business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	} else if businessSession != nil {
		business = session.Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	
	ctx.ViewData("business", business)
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.View("business_add4.html")
}

func BusinessAddStep5(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	//business := model.Business{}
	userSession := session.Get("user").(model.User)
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		//userSession := session.Get("user").(model.User)
		//business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	} else if businessSession != nil {
		//business = session.Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	
	/*userFolder := config.GetAppPath()+"resources/uploads/"+userSession.Id.Hex()+"/"+ctx.Params().Get("businessID")+"/gallery/"
	files, _ := ioutil.ReadDir(userFolder)
	var imageSlice []string
    for _, f := range files {
            imageSlice = append(imageSlice, f.Name())
    }*/
	images := model.User{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Find(bson.M{"_id": userSession.Id, "businesses._id": bson.ObjectIdHex(ctx.Params().Get("businessID"))}).Select(bson.M{"_id": 0, "businesses.$":1}).One(&images); err != nil {
		panic(err)
	}
	galleryImages := images.Businesses[0].Gallery
	profileImages := images.Businesses[0].Profile
	coverImages := images.Businesses[0].Cover
	Db.Close()
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.ViewData("userID", userSession.Id.Hex())
	ctx.ViewData("galleryImages", galleryImages)
	ctx.ViewData("profileImages", profileImages)
	ctx.ViewData("coverImages", coverImages)
	ctx.View("business_add5.html")
}

func BusinessAddStep6(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	//userSession := session.Get("user").(model.User)
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		//userSession := session.Get("user").(model.User)
		//business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	} else if businessSession != nil {
		//business = session.Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.View("business_add6.html")
}

func BusinessAddStep7(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	//userSession := session.Get("user").(model.User)
	businessSession := session.Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		//userSession := session.Get("user").(model.User)
		//business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	} else if businessSession != nil {
		//business = session.Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.View("business_add7.html")
}

func BusinessDelete(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	if ctx.Method() == "GET" {
		businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
		Db := db.MgoDb{}
		Db.Init()
		c := Db.C("users")
		
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$pull": bson.M{"businesses": bson.M{"_id":businessID}}}); err != nil {
			panic(err)
		}
		if err := c.Find(bson.M{"_id": userSession.Id}).One(&userSession); err != nil {
			panic(err)
		}
		
	}
	session.Set("user", userSession)
	
	ctx.ViewData("business", userSession.Businesses)
	
	ctx.View("businesses.html")	
}



func BusinessEventsTracker(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	var mapAddress = ""
	business := model.Business{}
	businessSession := session.Get("businessForm")
	if businessSession != nil {
		business = session.Get("businessForm").(model.Business)
	}
	if ctx.Params().Get("businessID") != "" {
		userSession := session.Get("user").(model.User)
		business = model.GetBusinessByIDandUser(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	}
	formError := []FormError{}
	if ctx.FormValue("back") == "0" {
		setValues := bson.M{}
		if ctx.FormValue("step") == "1" {
			business.Website = ctx.FormValue("business[website]")
			setValues["businesses.$.website"] = business.Website
			
			business.Name = ctx.FormValue("business[name]") 
			if business.Name == "" {
				formError = append(formError, FormError{"businessName", "This field is required"})
			} else {
				setValues["businesses.$.name"] = business.Name
			}
			business.Phone = ctx.FormValue("business[phone]")
			if business.Phone == "" {
				formError = append(formError, FormError{"businessPhone", "This field is required"})
			} else {
				setValues["businesses.$.phone"] = business.Phone
			}
			business.Address = ctx.FormValue("business[address]")
			if business.Address == "" {
				formError = append(formError, FormError{"businessAddress", "This field is required"})
			} else {
				setValues["businesses.$.address"] = business.Address
				mapAddress += business.Address+","
			}
			
			business.Address2 = ctx.FormValue("business[address2]")
			setValues["businesses.$.address2"] = business.Address2
			
			business.Area = ctx.FormValue("business[area]")
			if business.Area == "" {
				formError = append(formError, FormError{"businessArea", "This field is required"})
			} else {
				setValues["businesses.$.area"] = business.Area
				mapAddress += business.Area+","
			}
			business.State = ctx.FormValue("business[state]")
			if business.State == "" {
				formError = append(formError, FormError{"businessStateControl", "This field is required"})
			} else {
				setValues["businesses.$.state"] = business.State
			}
			business.City = ctx.FormValue("business[city]")
			if business.City == "" {
				formError = append(formError, FormError{"businessCity", "This field is required"})
			} else {
				setValues["businesses.$.city"] = business.City
				mapAddress += business.City+","
			}
			business.PostalCode = ctx.FormValue("business[postal_code]")
			if business.PostalCode == "" {
				formError = append(formError, FormError{"businessPostalCode", "This field is required"})
			} else {
				setValues["businesses.$.postalcode"] = business.PostalCode
				mapAddress += business.PostalCode+","
			}
			business.Country = ctx.FormValue("business[country]")
			if business.Country == "" {
				formError = append(formError, FormError{"businessCountry", "This field is required"})
			} else {
				setValues["businesses.$.country"] = business.Country
				mapAddress += business.Country
				
			}
			fmt.Println(mapAddress)
			if mapAddress != "" {
				coor := general.MapsInit(mapAddress)
				if (coor) != nil {
					setValues["businesses.$.map.lat"] = coor[0].Geometry.Location.Lat
					setValues["businesses.$.map.lng"] = coor[0].Geometry.Location.Lng
				}
			}
		}
		if ctx.FormValue("step") == "2" {
			business.Industry = ctx.FormValue("business[industry]")
			if business.Industry == "" {
				formError = append(formError, FormError{"businessIndustry", "This field is required"})
			} else {
				setValues["businesses.$.industry"] = business.Industry
			}
			business.YearsBusiness = ctx.FormValue("business[yearsBusiness]")
			if business.YearsBusiness == "" {
				formError = append(formError, FormError{"businessYearsBusiness", "This field is required"})
			} else {
				setValues["businesses.$.yearsBusiness"] = business.YearsBusiness
			}
			business.NumberEmployees = ctx.FormValue("business[numberEmployees]")
			if business.NumberEmployees == "" {
				formError = append(formError, FormError{"businessNumberEmployees", "This field is required"})
			} else {
				setValues["businesses.$.numberEmployees"] = business.NumberEmployees
			}
			business.SizeBusiness = ctx.FormValue("business[sizeBusiness]")
			if business.SizeBusiness == "" {
				formError = append(formError, FormError{"businessSizeBusiness", "This field is required"})
			} else {
				setValues["businesses.$.sizeBusiness"] = business.SizeBusiness
			}
			business.RelationshipBusiness = ctx.FormValue("business[relationshipBusiness]")
			if business.RelationshipBusiness == "" {
				formError = append(formError, FormError{"businessRelationshipBusiness", "This field is required"})
			} else {
				setValues["businesses.$.relationshipBusiness"] = business.RelationshipBusiness
			}
			business.HowYourHear = ctx.FormValue("business[howYourHear]")
			if business.HowYourHear == "" {
				formError = append(formError, FormError{"businessHowYourHear", "This field is required"})
			} else {
				setValues["businesses.$.howYouHear"] = business.HowYourHear
			}
			
			
		}
		
		if ctx.FormValue("step") == "3" {
			business.Description = template.HTML(ctx.FormValue("business[description]"))
			
			if business.Description == "" {
				formError = append(formError, FormError{"businessDescription", "This field is required"})
			}
			
			fmt.Println(business.Description)
			Db := db.MgoDb{}
			Db.Init()
			c := Db.C("users")
			userSession := session.Get("user").(model.User)

			if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$set": bson.M{"businesses.$.description": &business.Description}}); err != nil {
				log.Printf(err.Error())
			}

			user := model.User{}
			if err := c.Find(bson.M{"_id": userSession.Id}).One(&user); err != nil {
				panic(err)
			}
			Db.Close()
			session.Set("user", user)
		}
		
		if ctx.FormValue("step") == "4" {
			setValues["businesses.$.social.facebook"] = ctx.FormValue("business[facebook]")
			setValues["businesses.$.social.google"] = ctx.FormValue("business[google]")
			setValues["businesses.$.social.instagram"] = ctx.FormValue("business[instagram]")
			setValues["businesses.$.social.youtube"] = ctx.FormValue("business[youtube]")
			setValues["businesses.$.social.pinterest"] = ctx.FormValue("business[pinterest]")
			setValues["businesses.$.social.linkedin"] = ctx.FormValue("business[linkedin]")
			setValues["businesses.$.social.twitter"] = ctx.FormValue("business[twitter]")
		}
		
		
		if(ctx.FormValue("businessID") != "") {
			Db := db.MgoDb{}
			Db.Init()
			c := Db.C("users")
			userSession := session.Get("user").(model.User)
			
			if ctx.FormValue("step") == "1" || ctx.FormValue("step") == "2" || ctx.FormValue("step") == "4" {
				if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$set": setValues}); err != nil {
					log.Printf(err.Error())
				}
			}
			
			if ctx.FormValue("step") == "3" {
				if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$set": bson.M{"businesses.$.description": business.Description}}); err != nil {
					log.Printf(err.Error())
				}
			}

			
			user := model.User{}
			if err := c.Find(bson.M{"_id": userSession.Id}).One(&user); err != nil {
				panic(err)
			}
			Db.Close()
			session.Set("user", user)
		}
		
		
	}
	
	
	if ctx.FormValue("businessID") == "" {
		fmt.Println(ctx.FormValue("businessID"))
		session.Set("businessForm", business)
	}
	
	if ctx.FormValue("add") == "1" {
		if len(formError) > 0 {
			ctx.JSON(formError)
			return
		}
		ctx.Next()
		return
	}
	ctx.JSON(formError)

	
}


func BusinessAddFinish(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	businessSession := session.Get("businessForm")
	business := model.Business{}
	if businessSession != nil {
		userSession := session.Get("user").(model.User)
		Db := db.MgoDb{}
		Db.Init()
		c := Db.C("users")
		
		business = businessSession.(model.Business)
		business.Id = bson.NewObjectId()
		business.Description = template.HTML(ctx.FormValue("business[description]"))
		business.Verified = 0
		business.Premium = 0
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$push": bson.M{"businesses": &business}}); err != nil {
			panic(err)
		}
		Db.Close()
		user:= model.GetUserByID(userSession.Id)
		session.Delete("businessForm")
		session.Set("user", user)
		//session.Set("businessForm", nil)
		
		}
	fmt.Println(business.Id.Hex())
	ctx.JSON(business.Id.Hex())
}

func UploadFiles(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	var folder string
	if ctx.FormValue("imageType") == "gallery" {
		folder = "gallery"
	} else if ctx.FormValue("imageType") == "profile" {
		folder = "profile"
	} else if ctx.FormValue("imageType") == "cover" {
		folder = "cover"
	}
	file, info, err := ctx.FormFile("file")
	userSession := session.Get("user").(model.User)
	if err != nil {
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}

	defer file.Close()
	var imageNew image.Image
	if ctx.FormValue("imageFormat") == "image/jpeg" {
		imageNew, _ = jpeg.Decode(file)
	} else if ctx.FormValue("imageFormat") == "image/png" {
		imageNew, _ = png.Decode(file)
	}
	
	
	
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, imageNew, nil)
	image, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	
	newImageResized := resize.Resize(800, 0, image, resize.Lanczos3)
	
	b := newImageResized.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y
	ratio := "1"
	if imgHeight >= imgWidth {
		ratio = "0"
	}
	
	fnameOld := info.Filename
	extension := filepath.Ext(fnameOld)
	extension = ".jpg"
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fname := ""
	for i := 0; i < 30; i++ {
		index := r.Intn(len(chars))
		fname += chars[index : index+1]
	}
	fname += fname+"="+ratio
	
	fname += extension
	
		
	var userFolder = config.GetAppPath()+"resources/uploads/"+userSession.Id.Hex()+"/"+ctx.FormValue("businessID")+"/"+folder+"/"
	if _, err := os.Stat(userFolder); os.IsNotExist(err) {
		os.MkdirAll(userFolder, 0711)
	}

	out, err := os.OpenFile(userFolder+fname,
		os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}
	defer out.Close()
	
	jpeg.Encode(out, newImageResized, nil)
	//io.Copy(out, file)
	
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$push": bson.M{"businesses.$."+folder: fname}}); err != nil {
		panic(err)
	}
	Db.Close()
	ctx.JSON(map[string]interface{}{"fname": fname, "url": "/static/uploads/"+userSession.Id.Hex()+"/"+ctx.FormValue("businessID")+"/"+folder+"/"+fname})
}
func DeleteFile(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	var folder string
	if ctx.FormValue("imageType") == "gallery" {
		folder = "gallery"
	} else if ctx.FormValue("imageType") == "profile" {
		folder = "profile"
	} else if ctx.FormValue("imageType") == "cover" {
		folder = "cover"
	}
	userSession := session.Get("user").(model.User)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$pull": bson.M{"businesses.$."+folder: ctx.FormValue("id")}}); err != nil {
		return
	}
	var path = config.GetAppPath()+"resources/uploads/"+userSession.Id.Hex()+"/"+ctx.FormValue("businessID")+"/"+folder+"/"+ctx.FormValue("id")
	err := os.Remove(path)
	if err != nil {
		return
	}
}

func SendSms(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	business := model.Business{}
	phoneSms := ctx.FormValue("smsCode")
	formError := []FormError{}
	if phoneSms == "" {
		formError = append(formError, FormError{"smsCode", "This field is required"})
	} else if _, err := strconv.Atoi(phoneSms); err != nil {
		formError = append(formError, FormError{"smsCode", "This field has to be a number"})
	}
	
	if len(formError) == 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < 4; i++ {
			index := r.Intn(len(smsChars))
			business.SmsCode += smsChars[index : index+1]
		}
		fmt.Println(business.SmsCode)
		Db := db.MgoDb{}
		Db.Init()
		c := Db.C("users")
		if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$set": bson.M{"businesses.$.smsCode": business.SmsCode}}); err != nil {
			log.Printf(err.Error())
		}
		Db.Close()
		client := messagebird.New("qV8HkQdNlDD0UDa9Z3mrXnlXK")
		params := &messagebird.MessageParams{Reference: "MyReference"}
		message, _ := client.NewMessage(
		  "Edward",
		  []string{phoneSms},
		  business.SmsCode,
		  params)
		  
		 fmt.Println(message)
	}
	
	ctx.JSON(formError)
}

func VerifyCode(ctx context.Context) {
	session := db.Sessions.Start(ctx)
	business := model.User{}	
	verificationCode := ctx.FormValue("verificationCode")
	userSession := session.Get("user").(model.User)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Find(bson.M{"_id": userSession.Id, "businesses._id": bson.ObjectIdHex(ctx.FormValue("businessID")), "businesses.smsCode": verificationCode}).Select(bson.M{"_id": 0, "businesses.$":1}).One(&business); err != nil {
		ctx.JSON(map[string]interface{}{"response": false, "message": "Invalid Verification Code"})
		return
	}
	ctx.JSON(map[string]bool{"response": true})
	return
}

func BusinessProfilePage(ctx context.Context) {
	user := model.User{}
	var err error
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user = model.GetBusinessByID(businessID)
	if err != nil {
		ctx.NotFound()
		return
	}
	
	session := db.Sessions.Start(ctx)
	userSessionC := session.Get("user")
	if userSessionC != nil {	
		userSession := session.Get("user").(model.User)
		liked := IsLike(businessID, userSession.Id)
		ctx.ViewData("liked", liked)
	}
	
	ctx.ViewData("nrLIkes", len(user.Businesses[0].Likes))
	ctx.ViewData("business", user)
	ctx.View("business_profile/index.html")
}

func UpdatePhotos(ctx context.Context) {
	fmt.Println("da")
	images := model.User{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(ctx.FormValue("userID")), "businesses._id": bson.ObjectIdHex(ctx.FormValue("businessID"))}).Select(bson.M{"_id": 0, "businesses.$":1}).One(&images); err != nil {
		panic(err)
	}
	Db.Close()
	galleryImages := images.Businesses[0].Gallery
	profileImages := images.Businesses[0].Profile
	coverImages := images.Businesses[0].Cover
	
	Db.Close()
	//ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	//ctx.ViewData("userID", userSession.Id.Hex())
	ctx.JSON(map[string][]string{"galleryImages": galleryImages, "profileImages": profileImages, "coverImages": coverImages})
}

func BusinessProfileMaps(ctx context.Context) {
	//c := cache.New(5*time.Minute, 10*time.Minute)
	user := model.User{}
	var err error
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user = model.GetBusinessByID(businessID)
	if err != nil {
		ctx.NotFound()
		return
	}
	
	ctx.ViewData("business", user)
	ctx.View("business_profile/map.html")
}

func BusinessProfileWeb(ctx context.Context) {
	user := model.User{}
	var err error
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user = model.GetBusinessByID(businessID)
	if err != nil {
		ctx.NotFound()
		return
	}
	page, err := ctx.Params().GetInt("pageCount")
	if err != nil && page != -1 {
		ctx.NotFound()
		return
	}
	if page == -1 {
		page = 1
	}
	start := (page - 1) * 10 + 1

	results, count := general.Run(start, "005405349541100282636:amgvfhrjtka", user.Businesses[0].Name)
	resultsPage := float64(10)
	countF := float64(count)
	pages := math.Ceil(countF / resultsPage)

	if pages > resultsPage {
		pages = float64(10)
	}
	
	
	var pagesSlice []int
	for i := 1; i <= int(pages); i++ {
        pagesSlice = append(pagesSlice, i)
    }
	ctx.ViewData("pagesCount", pagesSlice)
	ctx.ViewData("business", user)
	ctx.ViewData("results", results)
	ctx.View("business_profile/web.html")
}

func BusinessProfileInternal(ctx context.Context) {
	user := model.User{}
	var err error
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user = model.GetBusinessByID(businessID)
	if err != nil {
		ctx.NotFound()
		return
	}
	page, err := ctx.Params().GetInt("pageCount")
	if err != nil && page != -1 {
		ctx.NotFound()
		return
	}
	if page == -1 {
		page = 1
	}
	start := (page - 1) * 10 + 1
	results, count := general.Run(start, "005405349541100282636:zequohzqzru", user.Businesses[0].Name)
	resultsPage := float64(10)
	countF := float64(count)
	pages := math.Ceil(countF / resultsPage)

	if pages > resultsPage {
		pages = float64(10)
	}
	
	var pagesSlice []int
	for i := 1; i <= int(pages); i++ {
        pagesSlice = append(pagesSlice, i)
    }
	ctx.ViewData("pagesCount", pagesSlice)
	ctx.ViewData("business", user)
	ctx.ViewData("results", results)
	ctx.View("business_profile/internal.html")
}

func BusinessLike(ctx context.Context) {
	liked := false
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	session := db.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	
	flag := IsLike(businessID, userSession.Id)
	
	if flag == false {
		if err := c.Update(bson.M{"businesses._id": businessID}, bson.M{"$addToSet": bson.M{"businesses.$.likes": userSession.Id}}); err != nil {
		}
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$addToSet": bson.M{"liked": businessID}}); err != nil {
		}
		liked = true
	} else {
		if err := c.Update(bson.M{"businesses._id": businessID}, bson.M{"$pull": bson.M{"businesses.$.likes": userSession.Id}}); err != nil {
		}
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$pull": bson.M{"liked": businessID}}); err != nil {
		}
	}
	
	_, user := model.GetBusinessByID(businessID)
	

	Db.Close()
	ctx.JSON(map[string]interface{}{"success": liked, "count": len(user.Businesses[0].Likes)})
}

func IsLike(businessID bson.ObjectId, userID bson.ObjectId) bool {
	flag := true
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	
	if err := c.Find(bson.M{"_id": userID, "liked": businessID}).One(nil); err != nil {
		flag = false
	}

	Db.Close()
	return flag;
}

func LIveSearch(ctx context.Context) {
	searchStr := ctx.FormValue("keyword")
	business := []model.User{}
	session := db.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	/*if err := c.Find(bson.M{"businesses": bson.M{ "$elemMatch": bson.M{"name":bson.M{"$regex": searchStr, "$options": "i"}}}}).Select(bson.M{"_id": 0, "businesses.$":1}).All(&business); err != nil {
		panic(err)
	}*/
	
	/*if err := c.Find(bson.M{"businesses.name":bson.M{"$regex": searchStr, "$options": "i"}}).Select(bson.M{"_id": 0, "businesses": bson.M{ "$elemMatch": bson.M{"name":bson.M{"$regex": searchStr, "$options": "i"}}}}).All(&business); err != nil {
		panic(err)
	}*/
	
	oe := bson.M{
         "$match" : bson.M{"businesses.name": bson.M{"$regex": searchStr, "$options": "i"}, "businesses.likes": userSession.Id},
	}
	oc := bson.M{
        "$unwind":"$businesses",
	}
	ob := bson.M{
        "$group": bson.M{
			"_id": "$_id",
			"businesses": bson.M{"$push": "$businesses"},
		},
	}
	oa := bson.M{
        "$project" :bson.M {
			"businesses.name":1,
			"businesses.pro":1,
			"pro": "$businesses.pro",
			"check": "$businesses.check",
		},
	}
	
	os := bson.M{
        "$sort" :bson.M {
			"pro":-1,
			"check":-1,
		},
	}

	pipe := c.Pipe([]bson.M{oe, oc, oe, os, ob, oa})
	
	if err := pipe.All(&business); err != nil {	
		log.Printf(err.Error())
	}
	pretty.Println(business)
	Db.Close()
	
	searchResults := []LiveResults{}
	for _, item := range business {
		for _, subItem := range item.Businesses {
			searchResults = append(searchResults, LiveResults{subItem.Name, subItem.Profile})
		}
    }
	ctx.JSON(map[string]interface{}{"results": searchResults})
}
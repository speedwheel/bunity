package controller

import(
	"github.com/kataras/iris/context"
	"app/model"
	"app/shared/db"
	ses "app/shared/session"
	"gopkg.in/mgo.v2/bson"
	"log"
	"github.com/logpacker/PayPal-Go-SDK"
	"os"
	"fmt"
	"app/config"
	//"io"
	//"io/ioutil"
	"math/rand"
	"time"
	//"bytes"
    "image"
	"image/jpeg"
    "image/png"
	"strconv"
	"github.com/messagebird/go-rest-api"
	"html/template"
	"app/shared/general"
	"math"
	//"github.com/patrickmn/go-cache"
	//"github.com/kr/pretty"
	"strings"
	"github.com/disintegration/imaging"
	"math/big"
	"net"

)

type FormError struct {
	Class string
	Message string
}

type LiveResults struct {
	Name string
	Image[] string
	Url string
	Category string
	UserId string
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const smsChars = "1234567890"

var (
	Countries = []string {"Afghanistan","Albania","Algeria","Andorra","Angola","Anguilla","Antigua Barbuda","Argentina","Armenia","Aruba","Australia","Austria","Azerbaijan","Bahamas","Bahrain","Bangladesh","Barbados","Belarus","Belgium","Belize","Benin","Bermuda","Bhutan","Bolivia","Bosnia Herzegovina","Botswana","Brazil","British Virgin Islands","Brunei","Bulgaria","Burkina Faso","Burundi","Cambodia","Cameroon","Canada","Cape Verde","Cayman Islands","Chad","Chile","China","Colombia","Congo","Cook Islands","Costa Rica","Cote D Ivoire","Croatia","Cruise Ship","Cuba","Cyprus","Czech Republic","Denmark","Djibouti","Dominica","Dominican Republic","Ecuador","Egypt","El Salvador","Equatorial Guinea","Estonia","Ethiopia","Falkland Islands","Faroe Islands","Fiji","Finland","France","French Polynesia","French West Indies","Gabon","Gambia","Georgia","Germany","Ghana","Gibraltar","Greece","Greenland","Grenada","Guam","Guatemala","Guernsey","Guinea","Guinea Bissau","Guyana","Haiti","Honduras","Hong Kong","Hungary","Iceland","India","Indonesia","Iran","Iraq","Ireland","Isle of Man","Israel","Italy","Jamaica","Japan","Jersey","Jordan","Kazakhstan","Kenya","Kuwait","Kyrgyz Republic","Laos","Latvia","Lebanon","Lesotho","Liberia","Libya","Liechtenstein","Lithuania","Luxembourg","Macau","Macedonia","Madagascar","Malawi","Malaysia","Maldives","Mali","Malta","Mauritania","Mauritius","Mexico","Moldova","Monaco","Mongolia","Montenegro","Montserrat","Morocco","Mozambique","Namibia","Nepal","Netherlands","Netherlands Antilles","New Caledonia","New Zealand","Nicaragua","Niger","Nigeria","Norway","Oman","Pakistan","Palestine","Panama","Papua New Guinea","Paraguay","Peru","Philippines","Poland","Portugal","Puerto Rico","Qatar","Reunion","Romania","Russia","Rwanda","Saint Pierre Miquelon","Samoa","San Marino","Satellite","Saudi Arabia","Senegal","Serbia","Seychelles","Sierra Leone","Singapore","Slovakia","Slovenia","South Africa","South Korea","Spain","Sri Lanka","St Kitts Nevis","St Lucia","St Vincent","St. Lucia","Sudan","Suriname","Swaziland","Sweden","Switzerland","Syria","Taiwan","Tajikistan","Tanzania","Thailand","Timor L'Este","Togo","Tonga","Trinidad Tobago","Tunisia","Turkey","Turkmenistan","Turks Caicos","Uganda","Ukraine","United Arab Emirates","United Kingdom","United States","United States Minor Outlying Islands","Uruguay","Uzbekistan","Venezuela","Vietnam","Virgin Islands (US)","Yemen","Zambia","Zimbabwe"}
	StatesUSA = []string {"Alabama","Alaska","Arizona","Arkansas","California","Colorado","Connecticut","Delaware","District of Columbia","Florida","Georgia","Hawaii","Idaho","Illinois","Indiana","Iowa","Kansas","Kentucky","Louisiana","Maine","Maryland","Massachusetts","Michigan","Minnesota","Mississippi","Missouri","Montana","Nebraska","Nevada","New Hampshire","New Jersey","New Mexico","New York","North Carolina","North Dakota","Ohio","Oklahoma","Oregon","Pennsylvania","Rhode Island","South Carolina","South Dakota","Tennessee","Texas","Utah","Vermont","Virginia","Washington","West Virginia","Wisconsin","Wyoming"}
	StatesCanada = []string {"British Columbia","Ontario","Newfoundland and Labrador","Nova Scotia", "Prince Edward Island", "New Brunswick", "Quebec", "Manitoba", "Saskatchewan", "Alberta", "Northwest Territories", "Nunavut","Yukon Territory"}
	StatesAustralia = []string {"New South Wales","Victoria","Queensland","Tasmania","South Australia","Western Australia","Northern Territory","Australian Capital Terrirory"}
	category = []string {"Accounting","Advertising","Automotive","Computers","Construction","Consulting","Dental","Education","Entertainment","Entrepreneur","Financial","Health Care","Internet","Law","Manufacturing","Marketing","Medical","Printing","Publishing","Real Estate","Restaurant","Retail","Sales","Service","Telecommunications","Travel","Wholesale"}
	YearsBusiness = []string {"Start-Up","1 - 5 Years","6 - 10 Years","11 + Years"}
	NrEmployees = []string {"1-5","6-10","11-20","21-50","51-100","101-250","251+"}
	SizeBusiness = []string{"$100k - $1MM","$101MM +","$26MM - $100MM","$2MM - $5MM","$6MM - $25MM","Just Starting","Less than $100k"}
	RelationshipBusiness = []string{"I'm the owner of this company.", "I work for this company.","I don't work here, but I'm acting on behalf of this company.","I'm a user of Zaphiri improving the business listing."}
	//howYouHear = []string{"Business Breakthroughs International","Business Coach/Consultant","Business Mastery Event","Chet Holmes","Friend/Associate","Mitch Russo","Radio","Search Engine","Social Networking","Solution Provider","Tony Robbins","Trade Show","Twitter","Web Seminar"}
	phonePrefix = map[string]string {"AD": "376","AE": "971","AF": "93","AG": "1-268","AI": "1-264","AL": "355","AM": "374","AO": "244","AQ": "672","AR": "54","AS": "1-684","AT": "43","AU": "61","AW": "297","AX": "358-18","AZ": "994","BA": "387","BB": "1-246","BD": "880","BE": "32","BF": "226","BG": "359","BH": "973","BI": "257","BJ": "229","BL": "590","BM": "1-441","BN": "673","BO": "591","BQ": "599","BR": "55","BS": "1-242","BT": "975","BV": "","BW": "267","BY": "375","BZ": "501","CA": "1","CC": "61","CD": "243","CF": "236","CG": "242","CH": "41","CI": "225","CK": "682","CL": "56","CM": "237","CN": "86","CO": "57","CR": "506","CU": "53","CV": "238","CW": "599","CX": "61","CY": "357","CZ": "420","DE": "49","DJ": "253","DK": "45","DM": "1-767","DO": "1-809 and 1-829","DZ": "213","EC": "593","EE": "372","EG": "20","EH": "212","ER": "291","ES": "34","ET": "251","FI": "358","FJ": "679","FK": "500","FM": "691","FO": "298","FR": "33","GA": "241","GB": "44","GD": "1-473","GE": "995","GF": "594","GG": "44-1481","GH": "233","GI": "350","GL": "299","GM": "220","GN": "224","GP": "590","GQ": "240","GR": "30","GS": "500","GT": "502","GU": "1-671","GW": "245","GY": "592","HK": "852","HM": " ","HN": "504","HR": "385","HT": "509","HU": "36","ID": "62","IE": "353","IL": "972","IM": "44-1624","IN": "91","IO": "246","IQ": "964","IR": "98","IS": "354","IT": "39","JE": "44-1534","JM": "1-876","JO": "962","JP": "81","KE": "254","KG": "996","KH": "855","KI": "686","KM": "269","KN": "1-869","KP": "850","KR": "82","KW": "965","KY": "1-345","KZ": "7","LA": "856","LB": "961","LC": "1-758","LI": "423","LK": "94","LR": "231","LS": "266","LT": "370","LU": "352","LV": "371","LY": "218","MA": "212","MC": "377","MD": "373","ME": "382","MF": "590","MG": "261","MH": "692","MK": "389","ML": "223","MM": "95","MN": "976","MO": "853","MP": "1-670","MQ": "596","MR": "222","MS": "1-664","MT": "356","MU": "230","MV": "960","MW": "265","MX": "52","MY": "60","MZ": "258","NA": "264","NC": "687","NE": "227","NF": "672","NG": "234","NI": "505","NL": "31","NO": "47","NP": "977","NR": "674","NU": "683","NZ": "64","OM": "968","PA": "507","PE": "51","PF": "689","PG": "675","PH": "63","PK": "92","PL": "48","PM": "508","PN": "870","PR": "1-787 and 1-939","PS": "970","PT": "351","PW": "680","PY": "595","QA": "974","RE": "262","RO": "40","RS": "381","RU": "7","RW": "250","SA": "966","SB": "677","SC": "248","SD": "249","SE": "46","SG": "65","SH": "290","SI": "386","SJ": "47","SK": "421","SL": "232","SM": "378","SN": "221","SO": "252","SR": "597","SS": "211","ST": "239","SV": "503","SX": "599","SY": "963","SZ": "268","TC": "1-649","TD": "235","TF": "","TG": "228","TH": "66","TJ": "992","TK": "690","TL": "670","TM": "993","TN": "216","TO": "676","TR": "90","TT": "1-868","TV": "688","TW": "886","TZ": "255","UA": "380","UG": "256","UM": "1","US": "1","UY": "598","UZ": "998","VA": "379","VC": "1-784","VE": "58","VG": "1-284","VI": "1-340","VN": "84","VU": "678","WF": "681","WS": "685","XK": "","YE": "967","YT": "262","ZA": "27","ZM": "260","ZW": "263"}
)



func Profile(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	ctx.Next()
	ctx.ViewData("user", session.Get("user"))
	ctx.View("profile.html")
	
}

func UsersEdit(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	ctx.ViewData("user", session.Get("user"))
	ctx.View("users_edit.html")
}

func UsersEditUpdate (ctx context.Context) {
	Db := db.MgoDb{}
	Db.Init()
	session := ses.Sessions.Start(ctx)
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
	session := ses.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	business := model.GetAllBusinessByUser(userSession.Id)
	ctx.ViewData("business", business)
	ctx.View("businesses.html")	
}

func BusinessAddStep1(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
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
	ctx.ViewData("countries", Countries)
	ctx.ViewData("statesUSA", StatesUSA)
	ctx.ViewData("statesCanada", StatesCanada)
	ctx.ViewData("statesAustralia", StatesAustralia)
	
	ctx.View("business_add.html")	
}

func BusinessAddStep2(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
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
	
	categ := model.GetAllCategories(0)
	
	ctx.ViewData("data", map[string][]string{
		"industries": category,
		"yearsBusiness": YearsBusiness,
		"nrEmployees": NrEmployees,
		"sizeBusiness": SizeBusiness,
		"relationshipBusiness": RelationshipBusiness,
		//"howYouHear": howYouHear,
	})
	
	ctx.ViewData("categ", categ)
	
	ctx.ViewData("business", business)
	ctx.View("business_add2.html")
}


func BusinessAddStep3(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
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
	session := ses.Sessions.Start(ctx)
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
	session := ses.Sessions.Start(ctx)
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
	images := model.GetBusinessByIDandUser( bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	galleryImages := images.Gallery
	profileImages := images.Profile
	coverImages := images.Cover
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.ViewData("userID", userSession.Id.Hex())
	ctx.ViewData("galleryImages", galleryImages)
	ctx.ViewData("profileImages", profileImages)
	ctx.ViewData("coverImages", coverImages)
	ctx.View("business_add5.html")
}

func BusinessAddStep6(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
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
	
	addr := ctx.RemoteAddr()
	ipInt := ipToInt(addr)
	countryCode := model.GetIpByCountry(ipInt)
	
	
	ctx.ViewData("countryCode", countryCode)
	ctx.ViewData("businessID", ctx.Params().Get("businessID"))
	ctx.ViewData("phonePrefix", phonePrefix)
	
	ctx.View("business_add6.html")
}

func BusinessAddStep7(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
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
	ctx.ViewData("businessID", ctx.Params().Get("businessID"))
	ctx.View("business_add7.html")
}

func BusinessDelete(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	if ctx.Method() == "GET" {
		businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
		Db := db.MgoDb{}
		Db.Init()
		c := Db.C("businesses")
		
		if err := c.Remove(bson.M{"_id": businessID, "user_id": userSession.Id}); err != nil {
			log.Printf(err.Error())
		}
		c = Db.C("users")
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$pull": bson.M{"businesses":businessID}}); err != nil {
			log.Printf(err.Error())
		}
		if err := c.Find(bson.M{"_id": userSession.Id}).One(&userSession); err != nil {
			log.Printf(err.Error())
		}
		
	}
	session.Set("user", userSession)
	business := model.GetAllBusinessByUser(userSession.Id)
	ctx.ViewData("business", business)
	
	ctx.View("businesses.html")	
}



func BusinessEventsTracker(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
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
		setValues2 := bson.M{}
		if ctx.FormValue("step") == "1" {
			business.Website = ctx.FormValue("business[website]")
			setValues["web"] = business.Website
			
			business.Name = ctx.FormValue("business[name]") 
			if business.Name == "" {
				formError = append(formError, FormError{"businessName", "This field is required"})
			} else {
				setValues["name"] = business.Name
				setValues2["business.$.name"] = business.Name
				//setValues["slug"] = strings.ToLower(business.Name)
				business.NameSplit = strings.Fields(strings.ToLower(business.Name))
				setValues["namearr"] = business.NameSplit
				setValues2["business.$.namearr"] = business.NameSplit
				
				
			}
			business.Phone = ctx.FormValue("business[phone]")
			if business.Phone == "" {
				formError = append(formError, FormError{"businessPhone", "This field is required"})
			} else {
				setValues["phone"] = business.Phone
			}
			
			business.Address.Country = ctx.FormValue("business[country]")
			if business.Address.Country == "" {
				formError = append(formError, FormError{"businessCountry", "This field is required"})
			} else {
				setValues["address.country"] = business.Address.Country
				mapAddress += business.Address.Country
				
			}
			
			business.Address.Address = ctx.FormValue("business[address]")
			if business.Address.Address == "" {
				formError = append(formError, FormError{"businessAddress", "This field is required"})
			} else {
				setValues["address.add"] = business.Address.Address
				mapAddress += business.Address.Address+","
			}
			
			business.Address.Address2 = ctx.FormValue("business[address2]")
			setValues["address.add2"] = business.Address.Address2
			
			business.Address.Area = ctx.FormValue("business[area]")
			if business.Address.Area == "" {
				formError = append(formError, FormError{"businessArea", "This field is required"})
			} else {
				setValues["address.area"] = business.Address.Area
				mapAddress += business.Address.Area+","
			}
			business.Address.State = ctx.FormValue("business[state]")
			if business.Address.State == "" && (business.Address.Country == "United States" || business.Address.Country == "Canada" || business.Address.Country == "Australia") {
				formError = append(formError, FormError{"businessStateControl", "This field is required"})
			} else {
				setValues["address.state"] = business.Address.State
			}
			business.Address.City = ctx.FormValue("business[city]")
			if business.Address.City == "" {
				formError = append(formError, FormError{"businessCity", "This field is required"})
			} else {
				setValues["address.city"] = business.Address.City
				mapAddress += business.Address.City+","
			}
			business.Address.PostalCode = ctx.FormValue("business[postal_code]")
			if business.Address.PostalCode == "" {
				formError = append(formError, FormError{"businessPostalCode", "This field is required"})
			} else {
				setValues["address.postalcode"] = business.Address.PostalCode
				mapAddress += business.Address.PostalCode+","
			}
			
			if mapAddress != "" {
				coor := general.MapsInit(mapAddress)
				if (coor) != nil {
					setValues["map.lat"] = coor[0].Geometry.Location.Lat
					setValues["map.lng"] = coor[0].Geometry.Location.Lng
				}
			}
		}
		if ctx.FormValue("step") == "2" {
			business.Category = ctx.FormValue("business[category]")
			if business.Category == "" {
				formError = append(formError, FormError{"businessCateg", "This field is required"})
			} else {
				setValues["category"] = business.Category
			}
			business.Category2 = ctx.FormValue("business[category2]")
			if business.Category2 == "" {
				formError = append(formError, FormError{"businessCateg2", "This field is required"})
			} else {
				setValues["category2"] = business.Category2
			}
			business.YearsBusiness = ctx.FormValue("business[yearsBusiness]")
			if business.YearsBusiness == "" {
				formError = append(formError, FormError{"businessYearsBusiness", "This field is required"})
			} else {
				setValues["ybuss"] = business.YearsBusiness
			}
			business.NumberEmployees = ctx.FormValue("business[numberEmployees]")
			if business.NumberEmployees == "" {
				formError = append(formError, FormError{"businessNumberEmployees", "This field is required"})
			} else {
				setValues["emp"] = business.NumberEmployees
			}
			business.SizeBusiness = ctx.FormValue("business[sizeBusiness]")
			if business.SizeBusiness == "" {
				formError = append(formError, FormError{"businessSizeBusiness", "This field is required"})
			} else {
				setValues["sbuss"] = business.SizeBusiness
			}
			business.RelationshipBusiness = ctx.FormValue("business[relationshipBusiness]")
			if business.RelationshipBusiness == "" {
				formError = append(formError, FormError{"businessRelationshipBusiness", "This field is required"})
			} else {
				setValues["relbuss"] = business.RelationshipBusiness
			}
		}
		
		
		if ctx.FormValue("step") == "4" {
			setValues["social.facebook"] = ctx.FormValue("business[facebook]")
			setValues["social.google"] = ctx.FormValue("business[google]")
			setValues["social.instagram"] = ctx.FormValue("business[instagram]")
			setValues["social.youtube"] = ctx.FormValue("business[youtube]")
			setValues["social.pinterest"] = ctx.FormValue("business[pinterest]")
			setValues["social.linkedin"] = ctx.FormValue("business[linkedin]")
			setValues["social.twitter"] = ctx.FormValue("business[twitter]")
		}
		
		
		if(ctx.FormValue("businessID") != "") {
			Db := db.MgoDb{}
			Db.Init()
			c := Db.C("businesses")
			userSession := session.Get("user").(model.User)
			
			if ctx.FormValue("step") == "1" || ctx.FormValue("step") == "2" || ctx.FormValue("step") == "4" {
				if err := c.Update(bson.M{"user_id": userSession.Id, "_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}, bson.M{"$set": setValues}); err != nil {
					log.Printf(err.Error())
				}
				c := Db.C("users")
				if err := c.Update(bson.M{"_id": userSession.Id, "business._id":bson.ObjectIdHex(ctx.FormValue("businessID"))}, bson.M{"$set": setValues2}); err != nil {
					log.Printf(err.Error())
				}
			}
			
			if ctx.FormValue("step") == "3" {
				business.Description = template.HTML(ctx.FormValue("business[description]"))
				if business.Description == "" {
					formError = append(formError, FormError{"businessDescription", "This field is required"})
				}
				if err := c.Update(bson.M{"user_id": userSession.Id, "_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}, bson.M{"$set": bson.M{"desc": business.Description}}); err != nil {
					log.Printf(err.Error())
				}
			}

			
			c = Db.C("users")
			user := model.User{}
			if err := c.Find(bson.M{"_id": userSession.Id}).One(&user); err != nil {
				panic(err)
			}
			Db.Close()
			session.Set("user", user)
		}
		
		
	}
	
	
	if ctx.FormValue("businessID") == "" {
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

func splitWord() {

}


func BusinessAddFinish(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	businessSession := session.Get("businessForm")
	business := model.Business{}
	businessS := model.BusinessS{}
	if businessSession != nil {
		userSession := session.Get("user").(model.User)
		Db := db.MgoDb{}
		Db.Init()
		c := Db.C("businesses")
		
		business = businessSession.(model.Business)
		business.Id = bson.NewObjectId()
		business.UserId = userSession.Id
		business.Description = template.HTML(ctx.FormValue("business[description]"))
		business.Verified = 0
		business.Premium = 0
		//business.Slug = strings.ToLower(business.Name)
		if err := c.Insert(&business); err != nil {
			panic(err)
		}
		businessS.Id = business.Id
		businessS.Name = business.Name
		businessS.NameSplit = business.NameSplit
		c = Db.C("users")
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$push": bson.M{"business": businessS}}); err != nil {
			panic(err)
		}
		Db.Close()
		user:= model.GetUserByID(userSession.Id)
		session.Delete("businessForm")
		session.Set("user", user)
		//session.Set("businessForm", nil)
		
		}
	ctx.JSON(business.Id.Hex())
}

func UploadFiles(ctx context.Context) {
	var image image.Image
	var resizeWidth int
	session := ses.Sessions.Start(ctx)
	var folder string
	if ctx.FormValue("imageType") == "gallery" {
		folder = "gallery"
		resizeWidth = 1024
	} else if ctx.FormValue("imageType") == "profile" {
		folder = "profile"
		resizeWidth = 160
	} else if ctx.FormValue("imageType") == "cover" {
		folder = "cover"
		resizeWidth = 840
	}
	file, _, err := ctx.FormFile("file")
	userSession := session.Get("user").(model.User)
	if err != nil {
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}
	
	defer file.Close()
	
	
	if ctx.FormValue("imageFormat") == "image/jpeg" {
		image, _ = jpeg.Decode(file)
	} else if ctx.FormValue("imageFormat") == "image/png" {
		image, _ = png.Decode(file)
	}
	
	
	
	//buf := new(bytes.Buffer)
	//jpeg.Encode(buf, imageNew, nil)
	//image, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	

	b := image.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y
	thumbSize := imgHeight
	
	ratio := "1"
	
	//fnameOld := info.Filename
	//extension := filepath.Ext(fnameOld)
	extension := ".jpg"
	
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
	
	
	
	
	if imgHeight >= imgWidth {
		ratio = "0"
		thumbSize = imgWidth
	}
	//img1  := imaging.CropAnchor(image, thumbSize, thumbSize, imaging.Center)
	fmt.Println(thumbSize)
	
	newImageResized := imaging.Resize(image, resizeWidth, 0, imaging.Lanczos)
	err = imaging.Save(newImageResized, userFolder+fname)
	if err != nil {
		log.Println("Save failed: %v", err)
	}

	/*out, err := os.OpenFile(userFolder+fname,
		os.O_WRONLY|os.O_CREATE, 0666)*/

	/*if err != nil {
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}*/
	//defer out.Close()
	
	//jpeg.Encode(out, newImageResized, nil)
	//io.Copy(out, file)
	
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	if err := c.Update(bson.M{"_id": bson.ObjectIdHex(ctx.FormValue("businessID")), "user_id": userSession.Id}, bson.M{"$push": bson.M{folder: fname}}); err != nil {
		panic(err)
	}
	Db.Close()
	ctx.JSON(map[string]interface{}{"fname": fname, "url": "/static/uploads/"+userSession.Id.Hex()+"/"+ctx.FormValue("businessID")+"/"+folder+"/"+fname})
}
func DeleteFile(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
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
	c := Db.C("businesses")
	if err := c.Update(bson.M{"_id": bson.ObjectIdHex(ctx.FormValue("businessID")), "user_id": userSession.Id}, bson.M{"$pull": bson.M{folder: ctx.FormValue("id")}}); err != nil {
		log.Printf(err.Error())
		return
	}
	var path = config.GetAppPath()+"resources/uploads/"+userSession.Id.Hex()+"/"+ctx.FormValue("businessID")+"/"+folder+"/"+ctx.FormValue("id")
	err := os.Remove(path)
	if err != nil {
		return
	}
	
}

func SendSms(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	business := model.Business{}
	phoneSms := ctx.FormValue("smsCode")
	prefix := ctx.FormValue("prefix")
	phoneNr := prefix + phoneSms

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
		c := Db.C("businesses")
		if err := c.Update(bson.M{"_id": bson.ObjectIdHex(ctx.FormValue("businessID")), "user_id": userSession.Id}, bson.M{"$set": bson.M{"smsCode": business.SmsCode}}); err != nil {
			log.Printf(err.Error())
		}
		Db.Close()
		client := messagebird.New("qV8HkQdNlDD0UDa9Z3mrXnlXK")
		params := &messagebird.MessageParams{Reference: "MyReference"}
		message, _ := client.NewMessage(
		  "Edward",
		  []string{phoneNr},
		  business.SmsCode,
		  params)
		  
		 fmt.Println(message)
	}
	
	ctx.JSON(formError)
}

func VerifyCode(ctx context.Context) {
	session := ses.Sessions.Start(ctx)
	business := model.User{}	
	verificationCode := ctx.FormValue("verificationCode")
	userSession := session.Get("user").(model.User)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(ctx.FormValue("businessID")), "user_id": userSession.Id , "smsCode": verificationCode}).One(&business); err != nil {
		ctx.JSON(map[string]interface{}{"response": false, "message": "Invalid Verification Code"})
		return
	}
	if err := c.Update(bson.M{"_id": bson.ObjectIdHex(ctx.FormValue("businessID")), "user_id": userSession.Id}, bson.M{"$set": bson.M{"check":1}}); err != nil {
		log.Printf(err.Error())
	}
	ctx.JSON(map[string]bool{"response": true})
	return
}

func BusinessProfilePage(ctx context.Context) {
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}

	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user := model.GetBusinessByID(businessID)
	if err != nil {
		ctx.NotFound()
		return
	}
	
	session := ses.Sessions.Start(ctx)
	userSessionC := session.Get("user")
	if userSessionC != nil {	
		userSession := session.Get("user").(model.User)
		liked := IsLike(businessID, userSession.Id)
		ctx.ViewData("liked", liked)
	}
	
	randB := model.GetRandomBusinesses(5)
	
	
	ctx.ViewData("nrLIkes", len(user.Likes))
	ctx.ViewData("business", user)
	ctx.ViewData("ranBusiness", randB)
	ctx.View("business_profile/index.html")
}

func UpdatePhotos(ctx context.Context) {
	images := model.Business{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	if err := c.Find(bson.M{"user_id": bson.ObjectIdHex(ctx.FormValue("userID")), "_id": bson.ObjectIdHex(ctx.FormValue("businessID"))}).One(&images); err != nil {
		panic(err)
	}
	Db.Close()
	galleryImages := images.Gallery
	profileImages := images.Profile
	coverImages := images.Cover
	
	Db.Close()
	//ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	//ctx.ViewData("userID", userSession.Id.Hex())
	ctx.JSON(map[string][]string{"galleryImages": galleryImages, "profileImages": profileImages, "coverImages": coverImages})
}

func BusinessProfileMaps(ctx context.Context) {
	//c := cache.New(5*time.Minute, 10*time.Minute)
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user := model.GetBusinessByID(businessID)
	if err != nil {
		ctx.NotFound()
		return
	}
	
	ctx.ViewData("business", user)
	
	ctx.View("business_profile/map.html")
}

func BusinessProfileWeb(ctx context.Context) {
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user := model.GetBusinessByID(businessID)
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

	results, count := general.Run(start, "005405349541100282636:amgvfhrjtka", user.Name)
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
	if !bson.IsObjectIdHex(ctx.Params().Get("businessID")){
		ctx.NotFound()
		return
	}
	businessID := bson.ObjectIdHex(ctx.Params().Get("businessID"))
	err, user := model.GetBusinessByID(businessID)
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
	results, count := general.Run(start, "005405349541100282636:zequohzqzru", user.Name)
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
	session := ses.Sessions.Start(ctx)
	userSession := session.Get("user").(model.User)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	
	flag := IsLike(businessID, userSession.Id)
	
	if flag == false {
		if err := c.Update(bson.M{"_id": businessID}, bson.M{"$addToSet": bson.M{"likes": userSession.Id}}); err != nil {
		}
		c := Db.C("users")
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$addToSet": bson.M{"liked": businessID}}); err != nil {
		}
		liked = true
	} else {
		if err := c.Update(bson.M{"_id": businessID}, bson.M{"$pull": bson.M{"likes": userSession.Id}}); err != nil {
		}
		c := Db.C("users")
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$pull": bson.M{"liked": businessID}}); err != nil {
		}
	}
	
	_, user := model.GetBusinessByID(businessID)
	

	Db.Close()
	ctx.JSON(map[string]interface{}{"success": liked, "count": len(user.Likes)})
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

func LiveSearch(ctx context.Context) {
	query := bson.M{}
	sliceBson := []bson.M{}
	searchStr := strings.ToLower(ctx.FormValue("keyword"))
	splidWords := strings.Fields(searchStr)
	for _, w := range splidWords {
        sliceBson = append(sliceBson, bson.M{"namearr": bson.M{"$regex": "^"+w}})
    }
	query["$and"] = sliceBson
	//query["slug"] = bson.M{"$regex": "\\b"+searchStr+"\\w*"}
	//query["nameSplit"] = bson.M{"$regex": "^"+searchStr}
	business := []model.Business{}
	business2 := []model.Business{}
	
	auth := ctx.Values().Get("auth").(bool)
	if auth  {
		session := ses.Sessions.Start(ctx)
		userSession := session.Get("user").(model.User)
		query["likes"] = userSession.Id
	}
	pageSize := 8
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	
	/*if err := c.Find(query).Limit(8).Select(bson.M{"slug":1, "profile":1}).Sort("-pro", "-check").All(&business); err != nil {
		panic(err)
	}*/
	oe := bson.M{
        "$match" : query,
	}
	oa := bson.M{
        "$project": bson.M {"name": 1, "profile": 1, "check":1, "pro":1, "user_id":1, "categ":1},
	}
	ol := bson.M{
        "$limit" :pageSize,
	}
	or := bson.M{
		"$sort": bson.D{
			bson.DocElem{Name: "pro", Value: -1},
			bson.DocElem{Name: "check", Value: -1},
			bson.DocElem{Name: "name", Value: 1},
		 },
	}
	pipe := c.Pipe([]bson.M{oe, oa, or, ol })
	
	if err := pipe.All(&business); err != nil {	
		log.Printf(err.Error())
	}
	var excludeIds []bson.ObjectId	
	searchResults := []LiveResults{}
	for _, item := range business {
		searchResults = append(searchResults, LiveResults{item.Name, item.Profile, item.Id.Hex(), item.Category, item.UserId.Hex()})
		excludeIds = append(excludeIds, item.Id)
    }

	if auth  {
		if len(searchResults) < 8 {
			query2 := bson.M{}
			query2["_id"]  = bson.M{"$nin": excludeIds}
			query2["$and"] = sliceBson
			/*if err := c.Find(bson.M{"_id": bson.M{"$nin": excludeIds}, "slug": bson.M{"$regex": searchStr}}).Limit(8).Select(bson.M{"slug":1, "profile":1}).Sort("-pro", "-check").All(&business2); err != nil {
				panic(err)
			}*/
			oe := bson.M{
				"$match" :query2,
			}
			pipe := c.Pipe([]bson.M{oe, oa, or, ol })
			if err := pipe.All(&business2); err != nil {	
				log.Printf(err.Error())
			}
			for _, item := range business2 {
				if len(searchResults) < 8 {
					searchResults = append(searchResults, LiveResults{item.Name, item.Profile, item.Id.Hex(), item.Category, item.UserId.Hex()})
				}
			}
		}
	}
	Db.Close()
	
	
	ctx.JSON(map[string]interface{}{"results": searchResults})
}

func BusinessSearch(ctx context.Context) {
	var pageNum int
	var err error
	query := bson.M{}
	searchStr := strings.ToLower(ctx.FormValue("q"))
	businessCategory := ctx.FormValue("business_category")
	likedFriends := ctx.FormValue("liked_friends")
	country := ctx.FormValue("country")
	if country != "" {
		query["address.country"] = country
	}
	query["name"] = bson.M{"$regex": searchStr}
	if(businessCategory != "") {
		query["$or"] = []bson.M{bson.M{"categ": businessCategory},bson.M{"categ2": businessCategory}}
	}
	verified := ctx.FormValue("verified")
	if verified == "1" {
		query["check"] = 1
	}
	
	//query["$text"] = bson.M{"$search": searchStr}
	var business []bson.M
	var business2 []bson.M
	var businessExclude []bson.M
	auth := ctx.Values().Get("auth").(bool)
	if auth  {
		session := ses.Sessions.Start(ctx)
		userSession := session.Get("user").(model.User)
		query["likes"] = userSession.Id
		//liked := userSession.Liked
		//query["likes"] = bson.M{"$in": liked}
		//pretty.Println(query["likes"])
	}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	

	
	if(ctx.Method() == "GET") {
		pageNum, err = ctx.Params().GetInt("pageCount")
		if err != nil {
			if pageNum == -1 { 
				pageNum = 1
			} else {
				log.Printf(err.Error())
				ctx.NotFound()
				return
			}
		}
	} else if(ctx.Method() == "POST") {
		pageNum, err = strconv.Atoi(ctx.FormValue("countPage"))
		if err != nil {
			log.Printf(err.Error())
		}
	}
	
	pageSize := 40
	skips := pageSize * (pageNum -1)
	/*if err := c.Find(query).Skip(skips).Limit(pageSize).Select(bson.M{"name":1, "profile":1, "description":1, "user_id":1, "likes":1}).Sort("-pro", "-check").All(&business); err != nil {
		panic(err)
	}*/
	oe := bson.M{
        "$match" :query,
	}
	oa := bson.M{
        "$project": bson.M {"pro": 1, "check": 1, "name":1, "profile":1, "desc":1, "user_id":1, "likes":1, "nrLikes": bson.M{ "$size": "$likes" }, "address.city": 1, "address.country": 1, "categ": 1},
	}
	ol := bson.M{
        "$limit" :pageSize,
	}
	os := bson.M{
        "$skip" :skips,
	}
	or := bson.M{
		"$sort": bson.D{
			bson.DocElem{Name: "pro", Value: -1},
			bson.DocElem{Name: "check", Value: -1},
			bson.DocElem{Name: "name", Value: 1},
		 },
	}

	
	pipe := c.Pipe([]bson.M{oe, or, oa, os, ol })
	
	pipe2 := c.Pipe([]bson.M{oe})
	
	if err := pipe.All(&business); err != nil {	
		log.Printf(err.Error())
	}

	
	
	countTotal, err := c.Find(query).Count()
	if err != nil {
		panic(err)
	}
	var excludeIds []bson.ObjectId	
	pages := math.Ceil(float64(countTotal) / float64(pageSize))
	remaining := countTotal % pageSize
    rest := 0
	if remaining > 0 {
		rest = pageSize - remaining
	}
	//fmt.Println(countTotal, pageSize, countTotal % pageSize)

	if auth  {
		if err := pipe2.All(&businessExclude); err != nil {	
			log.Printf(err.Error())
		}
		for _, item := range businessExclude {
			excludeIds = append(excludeIds, item["_id"].(bson.ObjectId))
		}
		query2 := bson.M{}
		query2["name"] = bson.M{"$regex": searchStr}
		query2["_id"] = bson.M{"$nin": excludeIds}
		if country != "" {
			query2["country"] = country
		}
		if(businessCategory != "") {
			query2["$or"] = []bson.M{bson.M{"categ": businessCategory},bson.M{"categ2": businessCategory}}
		}
		if verified == "1" {
			query2["check"] = 1
		}
		countTotal2, err := c.Find(query2).Count()
		if err != nil {
			panic(err)
		}
		
		businessesCount := len(business)
		
		if businessesCount < pageSize {
			bussiness2Needed := pageSize - businessesCount
			remaining := pageSize
		
			if(pageNum == int(pages)) {
				remaining = bussiness2Needed
			}
			
			skips2 := 0
			//fmt.Println(pageNum, pages)
			if pageNum > int(pages) {
				skipsB1 := pageSize * (pageNum -1 - int(pages))
				
				skips2 = skipsB1 + rest
			}
			if businessesCount == 0 && pageNum == 1 {
				skips2 = 0
			}
			//if businessesCount == 0 && pageNum
			
			oe := bson.M{
				"$match" :query2,
			}
			ol := bson.M{
				"$limit" :remaining,
			}
			os := bson.M{
				"$skip" :skips2,
			}
			pipe := c.Pipe([]bson.M{oe, or, oa, os, ol })
			if err := pipe.All(&business2); err != nil {	
				log.Printf(err.Error())
			}
			for _, item := range business2 {
				business = append(business, item)
			}
			//fmt.Println(skips2, remaining, countTotal2)
			
		}
		
		countTotal3 := countTotal2 + countTotal
		pages = math.Ceil(float64(countTotal3) / float64(pageSize))
		

	}
	
	var pagesSlice []int
	for i := 1; i <= int(pages); i++ {
		pagesSlice = append(pagesSlice, i)
	}
	
	Db.Close()
	
	if(ctx.Method() == "GET") {
		categ := model.GetAllCategories(0)
		ctx.ViewData("searchStr", searchStr)
		ctx.ViewData("business", business)
		ctx.ViewData("pageNum", pagesSlice)
		ctx.ViewData("industries", categ)
		ctx.ViewData("businessCategory", businessCategory)
		ctx.ViewData("verified", verified)
		ctx.ViewData("likedFriends", likedFriends)
		ctx.ViewData("countries", Countries)
		ctx.ViewData("selectedCountry", country)
		
		
		ctx.View("search.html")
	} else if (ctx.Method() == "POST") {
		if(int(pageNum) > int(pages)) {
			business = []bson.M{}
		}		
		ctx.JSON(map[string]interface{}{"businesses": business})
	}
}

func BusinessAllCategPage(ctx context.Context) {
	categ := model.GetAllCategories(0)
	ctx.ViewData("categ", categ)
	ctx.View("businessCateg.html")
}

func BussinessByCategPage(ctx context.Context) {
	categSlug := ctx.Params().Get("businessSlug")
	businesses := model.GetBusinessByCateg(categSlug)
	ctx.ViewData("businesses", businesses)
	ctx.View("businessSingleCateg.html")
}

func ipToInt(IP string) *big.Int {
    IPvAddr := net.ParseIP(IP).To4()
    IPvInt := big.NewInt(0)
    IPvInt.SetBytes(IPvAddr)
    return IPvInt
} 
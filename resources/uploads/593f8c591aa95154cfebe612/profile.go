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
)

type FormError struct {
	Class string
	Message string
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
	ctx.Next()
	ctx.ViewData("user", ctx.Session().Get("user"))
	ctx.View("profile.html")
}

func UsersEdit(ctx context.Context) {
	ctx.ViewData("user", ctx.Session().Get("user"))
	ctx.View("users_edit.html")
}

func UsersEditUpdate (ctx context.Context) {
	Db := db.MgoDb{}
	Db.Init()
	
	userSession := ctx.Session().Get("user").(model.User)
	userSession.Firstname = ctx.FormValue("firstname")
	userSession.Lastname = ctx.FormValue("lastname")
	userSession.Email = ctx.FormValue("email")
	
	c := Db.C("users")
	
	if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$set": bson.M{"firstname": userSession.Firstname, "lastname": userSession.Lastname, "email": userSession.Email}}); err != nil {
		log.Printf(err.Error())
		return
	}
	Db.Close()
	ctx.Session().Set("user", userSession)
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
	userSession := ctx.Session().Get("user").(model.User)
	ctx.ViewData("business", userSession.Businesses)
	ctx.View("businesses.html")	
}

func BusinessAddStep1(ctx context.Context) {
	business := model.Business{}
	businessSession := ctx.Session().Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := ctx.Session().Get("user").(model.User)
		business = model.GetBusinessByID(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil {
		business = ctx.Session().Get("businessForm").(model.Business)
	}
	ctx.ViewData("business", business)
	ctx.ViewData("countries", countries)
	ctx.ViewData("statesUSA", statesUSA)
	ctx.ViewData("statesCanada", statesCanada)
	ctx.ViewData("statesAustralia", statesAustralia)
	
	ctx.View("business_add.html")	
}

func BusinessAddStep11(ctx context.Context) {
	business := model.Business{}
	businessSession := ctx.Session().Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := ctx.Session().Get("user").(model.User)
		business = model.GetBusinessByID(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil {
		business = ctx.Session().Get("businessForm").(model.Business)
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

	//userSession := ctx.Session().Get("user").(model.User)
	//if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$push": bson.M{"businesses": bson.M{"_id": business.Businesses[0].Id, "name": business.Businesses[0].Name, "description": business.Businesses[0].Description, "phone": business.Businesses[0].Phone, "website": business.Businesses[0].Website, "address": business.Businesses[0].Address}}}); err != nil {
	//	panic(err)
	//}
	//Db.Close()
	
	
}

func BusinessAddStep2(ctx context.Context) {
	business := model.Business{}
	businessSession := ctx.Session().Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := ctx.Session().Get("user").(model.User)
		business = model.GetBusinessByID(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil { 
		business = ctx.Session().Get("businessForm").(model.Business)
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
	business := model.Business{}
	businessSession := ctx.Session().Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := ctx.Session().Get("user").(model.User)
		business = model.GetBusinessByID(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
		ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	} else if businessSession != nil { 
		business = ctx.Session().Get("businessForm").(model.Business)
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
	business := model.Business{}
	businessSession := ctx.Session().Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		userSession := ctx.Session().Get("user").(model.User)
		business = model.GetBusinessByID(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	} else if businessSession != nil {
		business = ctx.Session().Get("businessForm").(model.Business)
	} else {
		ctx.Redirect("/business/add")
	}
	
	ctx.ViewData("business", business)
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.View("business_add3.html")
}

func BusinessAddStep4(ctx context.Context) {
	//business := model.Business{}
	userSession := ctx.Session().Get("user").(model.User)
	businessSession := ctx.Session().Get("businessForm")
	if ctx.Params().Get("businessID") != "" {
		//userSession := ctx.Session().Get("user").(model.User)
		//business = model.GetBusinessByID(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	} else if businessSession != nil {
		//business = ctx.Session().Get("businessForm").(model.Business)
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
	fmt.Println(profileImages)
	Db.Close()
	ctx.ViewData("businessID", bson.ObjectIdHex(ctx.Params().Get("businessID")).Hex())
	ctx.ViewData("userID", userSession.Id.Hex())
	ctx.ViewData("galleryImages", galleryImages)
	ctx.ViewData("profileImages", profileImages)
	ctx.View("business_add4.html")
}


func BusinessDelete(ctx context.Context) {
	userSession := ctx.Session().Get("user").(model.User)
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
	ctx.Session().Set("user", userSession)
	
	ctx.ViewData("business", userSession.Businesses)
	
	ctx.View("businesses.html")	
}



func BusinessEventsTracker(ctx context.Context) {
	business := model.Business{}
	businessSession := ctx.Session().Get("businessForm")
	if businessSession != nil {
		business = ctx.Session().Get("businessForm").(model.Business)
	}
	if ctx.Params().Get("businessID") != "" {
		userSession := ctx.Session().Get("user").(model.User)
		business = model.GetBusinessByID(bson.ObjectIdHex(ctx.Params().Get("businessID")), userSession.Id)
	}
	formError := []FormError{}
	if ctx.FormValue("back") == "0" {
		if ctx.FormValue("step") == "1" {
			business.Name = ctx.FormValue("business[name]")
			if business.Name == "" {
				formError = append(formError, FormError{"businessName", "This field is required"})
			}
			business.Phone = ctx.FormValue("business[phone]")
			if business.Phone == "" {
				formError = append(formError, FormError{"businessPhone", "This field is required"})
			}
			business.Address = ctx.FormValue("business[address]")
			if business.Address == "" {
				formError = append(formError, FormError{"businessAddress", "This field is required"})
			}
			business.Area = ctx.FormValue("business[area]")
			if business.Area == "" {
				formError = append(formError, FormError{"businessArea", "This field is required"})
			}
			business.State = ctx.FormValue("business[state]")
			if business.State == "" {
				formError = append(formError, FormError{"businessStateControl", "This field is required"})
			}
			business.PostalCode = ctx.FormValue("business[postal_code]")
			if business.PostalCode == "" {
				formError = append(formError, FormError{"businessPostalCode", "This field is required"})
			}
			business.Country = ctx.FormValue("business[country]")
			if business.Country == "" {
				formError = append(formError, FormError{"businessCountry", "This field is required"})
			}
		}
		if ctx.FormValue("step") == "2" {
			business.Industry = ctx.FormValue("business[industry]")
			if business.Industry == "" {
				formError = append(formError, FormError{"businessIndustry", "This field is required"})
			}
			business.YearsBusiness = ctx.FormValue("business[yearsBusiness]")
			if business.YearsBusiness == "" {
				formError = append(formError, FormError{"businessYearsBusiness", "This field is required"})
			}
			business.NumberEmployees = ctx.FormValue("business[numberEmployees]")
			if business.NumberEmployees == "" {
				formError = append(formError, FormError{"businessNumberEmployees", "This field is required"})
			}
			business.SizeBusiness = ctx.FormValue("business[sizeBusiness]")
			if business.SizeBusiness == "" {
				formError = append(formError, FormError{"businessSizeBusiness", "This field is required"})
			}
			business.RelationshipBusiness = ctx.FormValue("business[relationshipBusiness]")
			if business.RelationshipBusiness == "" {
				formError = append(formError, FormError{"businessRelationshipBusiness", "This field is required"})
			}
			business.HowYourHear = ctx.FormValue("business[howYourHear]")
			if business.HowYourHear == "" {
				formError = append(formError, FormError{"businessHowYourHear", "This field is required"})
			}
		}
		if(ctx.FormValue("businessID") != "") {
			Db := db.MgoDb{}
			Db.Init()
			c := Db.C("users")
			userSession := ctx.Session().Get("user").(model.User)
			
			if ctx.FormValue("step") == "1" {
				if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$set": bson.M{"businesses.$.name": business.Name, "businesses.$.phone": business.Phone, "businesses.$.state": business.State, "businesses.$.address": business.Address, "businesses.$.postalcode": business.PostalCode, "businesses.$.country": business.Country, "businesses.$.area": business.Area}}); err != nil {
					log.Printf(err.Error())
				}
			}
			
			if ctx.FormValue("step") == "2" {
				if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$set": bson.M{"businesses.$.industry": business.Industry, "businesses.$.yearsBusiness": business.YearsBusiness, "businesses.$.numberEmployees": business.NumberEmployees, "businesses.$.sizeBusiness": business.SizeBusiness, "businesses.$.relationshipBusiness": business.RelationshipBusiness, "businesses.$.howYourHear": business.HowYourHear}}); err != nil {
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
			ctx.Session().Set("user", user)
		}
		
		if ctx.FormValue("step") == "3" {
			business.Description = ctx.FormValue("business[description]")
			
			if business.Description == "" {
				formError = append(formError, FormError{"businessDescription", "This field is required"})
			}
			
			fmt.Println(business.Description)
			Db := db.MgoDb{}
			Db.Init()
			c := Db.C("users")
			userSession := ctx.Session().Get("user").(model.User)

			if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$set": bson.M{"businesses.$.description": &business.Description}}); err != nil {
				log.Printf(err.Error())
			}

			user := model.User{}
			if err := c.Find(bson.M{"_id": userSession.Id}).One(&user); err != nil {
				panic(err)
			}
			Db.Close()
			ctx.Session().Set("user", user)
		}
	}
	
	
	if ctx.FormValue("businessID") == "" {
		fmt.Println(ctx.FormValue("businessID"))
		ctx.Session().Set("businessForm", business)
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
	businessSession := ctx.Session().Get("businessForm")
	business := model.Business{}
	if businessSession != nil {
		userSession := ctx.Session().Get("user").(model.User)
		Db := db.MgoDb{}
		Db.Init()
		c := Db.C("users")
		
		business = businessSession.(model.Business)
		business.Id = bson.NewObjectId()
		business.Description = ctx.FormValue("business[description]")
		if err := c.Update(bson.M{"_id": userSession.Id}, bson.M{"$push": bson.M{"businesses": &business}}); err != nil {
			panic(err)
		}
		Db.Close()
		user:= model.GetUserByID(userSession.Id)
		ctx.Session().Delete("businessForm")
		ctx.Session().Set("user", user)
		//ctx.Session().Set("businessForm", nil)
		
		}
	fmt.Println(business.Id.Hex())
	ctx.JSON(business.Id.Hex())
}

func UploadFiles(ctx context.Context) {
	var folder string
	if ctx.FormValue("imageType") == "gallery" {
		folder = "gallery"
	} else if ctx.FormValue("imageType") == "profile" {
		folder = "profile"
	}
	file, info, err := ctx.FormFile("file")
	userSession := ctx.Session().Get("user").(model.User)
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
	
	newImageResized := resize.Resize(400, 0, image, resize.Lanczos3)
	
	fnameOld := info.Filename
	extension := filepath.Ext(fnameOld)
	extension = ".jpg"
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fname := ""
	for i := 0; i < 30; i++ {
		index := r.Intn(len(chars))
		fname += chars[index : index+1]
	}
	
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
//ctx.JSON(business.Id.Hex())
}
func DeleteFile(ctx context.Context) {
	userSession := ctx.Session().Get("user").(model.User)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("users")
	if err := c.Update(bson.M{"_id": userSession.Id, "businesses": bson.M{ "$elemMatch": bson.M{"_id":bson.ObjectIdHex(ctx.FormValue("businessID"))}}}, bson.M{"$pull": bson.M{"businesses.$.gallery": ctx.FormValue("id")}}); err != nil {
		panic(err)
	}
}

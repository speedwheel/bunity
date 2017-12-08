package user

import(
	"app/model"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"app/shared/db"
	"log"
	"net/url"
	"app/config"
	//"github.com/kr/pretty"
	"strconv"
	"strings"
	"os"
	"image"
	"image/jpeg"
    "image/png"
	"github.com/disintegration/imaging"
	//"fmt"
	"math/rand"
	"mime/multipart"
	"time"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Id     bson.ObjectId `json:"id" bson:"_id"  form:"-"`
	Username  string `json:"username" bson:"username"  form:"username" facebook:"username"`
	Email  string `json:"email" bson:"email"  form:"email" facebook:"email"`
	Password string `json:"password" bson:"password,omitempty" form:"password,omitempty"`
	Owned []bson.ObjectId `json:"owned,omitempty" bson:"owned,omitempty"  form:"owned,omitempty"`
}

const (
	userCol = "users"
	businessCol = "businesses"
	businessCategCol = "categories"
	adminBusinessChat = "adminBusinessChat"
	adminActivityType = "adminActivityType"
	adminUser = "admin"
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	userCols = []string{"", "", "firstname", "lastname", "email", "business.name"}
	busCols = []string{"slug", "phone", "country", "website"}
)

type DataSource struct {
	Users []model.User
	User model.User
	Businesses []model.Business

}

type ( 
	ChatMessage struct {
		Id  bson.ObjectId `json:"id" bson:"_id" form:"id"`
		ParentID  bson.ObjectId `json:"parent_id" bson:"parent_id,omitempty" form:"parent_id"`
		BusinessID  bson.ObjectId `json:"business_id" bson:"business_id" form:"business_id"`
		Time  time.Time `json:"time" bson:"time" form:"time"`
		Text  string `json:"text" bson:"text" form:"text"`
		Author Author `json:"author" bson:"author" form:"author"`
		ActivityType ActivityType `json:"activity_type" bson:"activity_type,omitempty" form:"activity_type"`
	}
	
	Author struct {
		Id  bson.ObjectId `json:"id" bson:"_id"  form:"id"`
		Name string `json:"name" bson:"name" form:"name"`
	}
)

type ActivityType struct {
	Id  bson.ObjectId `json:"id" bson:"_id" form:"id"`
	Name string `json:"name" bson:"name" form:"name"`
}


func NewDataSource() *DataSource {
	return &DataSource{
		Users: []model.User{},
		User: model.User{},
		Businesses: []model.Business{},
	}
}

func (d *DataSource) AdminLogin(username, password string) (Admin, error) {
	ok := true
	cfg := config.Init()
	var err error
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminUser)
	admin := Admin{}
	
	if err = c.Find(bson.M{"username": username}).One(&admin); err != nil {
		ok = false
	}
	Db.Close()
	if !ok {
		return Admin{}, err
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password+cfg.User.SecretKey)); err != nil {
		return Admin{}, err
	}
	return admin, nil
}

func (d *DataSource) GetAllUsers(urlQuery url.Values) ([]model.User, int, int) {
	Db := db.MgoDb{}
	query := bson.M{}
	var sortColInt int
	var stringI string
	sliceBson := []bson.M{}
	sortDoc := bson.D{}
	limit := 0
	limit, _ = strconv.Atoi(urlQuery["length"][0])
	skips, _ := strconv.Atoi(urlQuery["start"][0])
	sortValue := 1
	for i := 0; i < len(userCols) - 2; i++ {
	stringI = strconv.Itoa(i)
		if urlQuery["order["+stringI+"][column]"] != nil {
			sortColInt, _ = strconv.Atoi(urlQuery["order["+stringI+"][column]"][0])
			
			if urlQuery["order["+stringI+"][dir]"][0] == "desc" {
				sortValue = -1
			}
			sortDoc = append(sortDoc, bson.DocElem{Name: userCols[sortColInt], Value: sortValue})
		}
	}
	
	
	searchValue := urlQuery["search[value]"][0]
	if searchValue != "" {
		splidWords := strings.Fields(searchValue)
		for _, w := range splidWords {
			sliceBson = append(sliceBson, bson.M{"$or": []bson.M{bson.M{"firstname": bson.M{"$regex": "^"+w, "$options" : "i"}}, bson.M{"lastname": bson.M{"$regex": "^"+w, "$options" : "i"}},bson.M{"email": bson.M{"$regex": "^"+w, "$options" : "i"}}, bson.M{"business.namearr": bson.M{"$regex": "^"+w, "$options" : "i"}}}})
		}
	
		query["$and"] = sliceBson
	}
	
	pm := bson.M{
        "$match" :query,
	}
	
	/*pp := bson.M{
        "$project": bson.M {"_id": 0},
	}*/
	
	pl := bson.M{
        "$limit" :limit,
	}
	
	ps := bson.M{
        "$skip" :skips,
	}
	po := bson.M{
		"$sort": sortDoc,
	}

	
	Db.Init()
	c := Db.C(userCol)
	
	pipe := c.Pipe([]bson.M{pm, po, /*pp,*/ ps, pl })
	if err := pipe.All(&d.Users); err != nil {	
		log.Printf(err.Error())
	}
	
	CountFiltered, err := c.Find(query).Count()
	if err != nil {
		panic(err)
	}
	Count, err := c.Find(nil).Count()
	if err != nil {
		panic(err)
	}
	Db.Close()
	
	return d.Users, CountFiltered, Count
}

func (d *DataSource) GetAllBusinesses(urlQuery url.Values) ([]model.Business, int, int) {
	Db := db.MgoDb{}
	query := bson.M{}
	var sortColInt int
	var stringI string
	sliceBson := []bson.M{}
	sortDoc := bson.D{}
	limit, _ := strconv.Atoi(urlQuery["length"][0])
	skips, _ := strconv.Atoi(urlQuery["start"][0])
	sortValue := 1
	for i := 0; i < 3; i++ {
	stringI = strconv.Itoa(i)
		if urlQuery["order["+stringI+"][column]"] != nil {
			sortColInt, _ = strconv.Atoi(urlQuery["order["+stringI+"][column]"][0])
			
			if urlQuery["order["+stringI+"][dir]"][0] == "desc" {
				sortValue = -1
			}
			sortDoc = append(sortDoc, bson.DocElem{Name: busCols[sortColInt], Value: sortValue})
		}
	}
	
	
	searchValue := urlQuery["search[value]"][0]
	if searchValue != "" {
		splidWords := strings.Fields(searchValue)
		for _, w := range splidWords {
			sliceBson = append(sliceBson, bson.M{"$or": []bson.M{bson.M{"namearr": bson.M{"$regex": "^"+w}}}})
		}
	
		query["$and"] = sliceBson
	}
	
	pm := bson.M{
        "$match" :query,
	}
	
	/*pp := bson.M{
        "$project": bson.M {},
	}*/
	
	pl := bson.M{
        "$limit" :limit,
	}
	
	ps := bson.M{
        "$skip" :skips,
	}
	po := bson.M{
		"$sort": sortDoc,
	}

	
	Db.Init()
	c := Db.C(businessCol)
	
	pipe := c.Pipe([]bson.M{pm, po, /*pp,*/ ps, pl })
	if err := pipe.AllowDiskUse().All(&d.Businesses); err != nil {	
		log.Printf(err.Error())
	}
	
	CountFiltered, err := c.Find(query).Count()
	if err != nil {
		panic(err)
	}
	Count, err := c.Find(nil).Count()
	if err != nil {
		panic(err)
	}
	Db.Close()
	
	return d.Businesses, CountFiltered, Count
}

func (d *DataSource) GetBusinessByID(businessID string) model.Business {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(businessCol)
	business := model.Business{}
	
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(businessID)}).One(&business); err != nil {
		log.Printf(err.Error())
	}
	Db.Close()
	return business
}

func (d *DataSource) UpdateBusinessByID(businessID string, business bson.M, user bson.M, userID string) bool {
	businessIDHex := bson.ObjectIdHex(businessID)
	userIDHex := bson.ObjectIdHex(userID)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(businessCol)
	if err := c.Update(bson.M{"_id":businessIDHex}, bson.M{"$set": business}); err != nil {
		log.Printf(err.Error())
	}
	c = Db.C(userCol)
	if err := c.Update(bson.M{"_id": userIDHex, "business._id":businessIDHex}, bson.M{"$set": user}); err != nil {
		log.Printf(err.Error())
	}
	Db.Close()
	return true
}

func (d *DataSource) DeleteBusinessByID(businessID string/*, userID string*/) bool {
	businessIDHex := bson.ObjectIdHex(businessID)
	//userIDHex := bson.ObjectIdHex(userID)
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(businessCol)
	
	if err := c.Remove(bson.M{"_id": businessIDHex/*, "user_id": userID*/}); err != nil {
		log.Printf(err.Error())
		return false
	}
	
	c = Db.C(userCol)
	
	if err := c.Update(bson.M{/*"_id": userIDHex,*/ "business._id":businessIDHex}, bson.M{"$pull": bson.M{"business": bson.M{"_id": businessIDHex}}}); err != nil {
		log.Printf(err.Error())
	}
	return true;
}

func (d *DataSource) GetAllBusinessCategories() []model.Category {
	var categories []model.Category
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(businessCategCol)
	if err := c.Find(nil).All(&categories); err != nil {	
		log.Printf(err.Error())
	}
	Db.Close()
	return categories
}

func (d *DataSource) DeletePicture(businessID string, userID string, fileID string, imageType string) bool {
	var folder string
	if imageType == "gallery" {
		folder = "gallery"
	} else if imageType == "profile" {
		folder = "profile"
	} else if imageType == "cover" {
		folder = "cover"
	}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(businessCol)
	if err := c.Update(bson.M{"_id": bson.ObjectIdHex(businessID)}, bson.M{"$pull": bson.M{folder: fileID}}); err != nil {
		log.Printf(err.Error())
		return false
	}
	path := config.GetAppPath()+"resources/uploads/"+userID+"/"+businessID+"/"+folder+"/"+fileID
	err := os.Remove(path)
	if err != nil {
		//log.Printf(err.Error())
		Db.Close()
		return false
	}
	return true
}

func (d *DataSource) AddPicture(userID string, businessID string, imageFormat string, imageType string, file multipart.File) (string, string) {
	var image image.Image
	var resizeWidth int

	var folder string
	if imageType == "gallery" {
		folder = "gallery"
		resizeWidth = 1024
	} else if imageType == "profile" {
		folder = "profile"
		resizeWidth = 160
	} else if imageType == "cover" {
		folder = "cover"
		resizeWidth = 840
	}
	
	if imageFormat == "image/jpeg" {
		image, _ = jpeg.Decode(file)
	} else if imageFormat == "image/png" {
		image, _ = png.Decode(file)
	}

	b := image.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y
	//thumbSize := imgHeight
	
	ratio := "1"
	
	extension := ".jpg"
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fname := ""
	for i := 0; i < 30; i++ {
		index := r.Intn(len(chars))
		fname += chars[index : index+1]
	}
	fname += fname+"="+ratio
	
	fname += extension
	
		
	var userFolder = config.GetAppPath()+"resources/uploads/"+userID+"/"+businessID+"/"+folder+"/"
	if _, err := os.Stat(userFolder); os.IsNotExist(err) {
		os.MkdirAll(userFolder, 0711)
	}
	

	if imgHeight >= imgWidth {
		ratio = "0"
		//thumbSize = imgWidth
	}
	//img1  := imaging.CropAnchor(image, thumbSize, thumbSize, imaging.Center)
	
	newImageResized := imaging.Resize(image, resizeWidth, 0, imaging.Lanczos)
	err := imaging.Save(newImageResized, userFolder+fname)
	if err != nil {
		log.Println("Save failed: %v", err)
	}

	
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(businessCol)
	if err := c.Update(bson.M{"_id": bson.ObjectIdHex(businessID), "user_id": bson.ObjectIdHex(userID)}, bson.M{"$push": bson.M{folder: fname}}); err != nil {
		log.Printf(err.Error())
	}
	Db.Close()
	
	
	url := "/static/uploads/"+userID+"/"+businessID+"/"+folder+"/"+fname
	return fname, url
}

func (d *DataSource) InsertBusinessComment(chatMessage ChatMessage) bool {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminBusinessChat)
	if err := c.Insert(&chatMessage); err != nil {
		log.Printf(err.Error())
		Db.Close()
		return false
	}
	return true
}

func (d *DataSource) GetCommentsByBusiness(businessID string) []ChatMessage {
	chatMessage := []ChatMessage{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminBusinessChat)
	if err := c.Find(bson.M{"business_id": bson.ObjectIdHex(businessID)}).All(&chatMessage); err != nil {	
		log.Printf(err.Error())
	}
	return chatMessage
}

func (d *DataSource) GetCommentsById(Id string) ChatMessage {
	chatMessage := ChatMessage{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminBusinessChat)
	if err := c.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&chatMessage); err != nil {	
		log.Printf(err.Error())
	}
	return chatMessage
}

func (d *DataSource) UpdateBusinessCommentByID(chatID string, userID bson.ObjectId, msg string) bool {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminBusinessChat)
	if err := c.Update(bson.M{"_id":bson.ObjectIdHex(chatID), "author._id": userID}, bson.M{"$set": bson.M{"text": msg}}); err != nil {
		log.Printf(err.Error())
		return false
	}
	return true
}

func (d *DataSource) AddNewActivity(activity ActivityType) bool {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminActivityType)
	if err := c.Insert(&activity); err != nil {
		log.Printf(err.Error())
		return false
	}
	return true
}

func (d *DataSource) GetAllActivitiesTD(urlQuery url.Values) ([]ActivityType, int, int) {
	activityType := []ActivityType{}
	query := bson.M{}
	limit := 0
	limit, _ = strconv.Atoi(urlQuery["length"][0])
	skips, _ := strconv.Atoi(urlQuery["start"][0])
	searchValue := urlQuery["search[value]"][0]
	if searchValue != "" {
		query["name"] = bson.M{"$regex": "^"+searchValue}
	}
	//sortDoc := bson.D{}
	
	
	pm := bson.M{
        "$match" :query,
	}
	
	pl := bson.M{
        "$limit" :limit,
	}
	
	ps := bson.M{
        "$skip" :skips,
	}
	po := bson.M{

		"$sort": bson.D {
			bson.DocElem{Name: "name", Value: 1},
		},
	}

	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminActivityType)
	
	pipe := c.Pipe([]bson.M{pm, po, ps, pl })
	if err := pipe.All(&activityType); err != nil {	
		log.Printf(err.Error())
	}
	fmt.Println(activityType)
	CountFiltered, err := c.Find(query).Count()
	if err != nil {
		panic(err)
	}
	Count, err := c.Find(nil).Count()
	if err != nil {
		panic(err)
	}
	Db.Close()
	return activityType, CountFiltered, Count
}

func (d *DataSource) UpdateActivityType(activity ActivityType) bool {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminActivityType)
	if err := c.Update(bson.M{"_id": activity.Id,}, activity); err != nil {
		log.Printf(err.Error())
		Db.Close()
		return false
	}
	c = Db.C(adminBusinessChat)
	if _, err := c.UpdateAll(bson.M{"activity_type._id": activity.Id,}, bson.M{"$set": bson.M{"activity_type.name": activity.Name}}); err != nil {
		log.Printf(err.Error())
		Db.Close()
		return false
	}
	Db.Close()
	return true
}

func (d *DataSource) DeleteActivityType(activity bson.ObjectId) bool {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminActivityType)
	if err := c.Remove(bson.M{"_id": activity}); err != nil {
		log.Printf(err.Error())
		Db.Close()
		return false
	}
	Db.Close()
	return true
}

func (d *DataSource) GetAllActivities() []ActivityType {
	activityType := []ActivityType{}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminActivityType)
	if err := c.Find(nil).All(&activityType); err != nil {	
		log.Printf(err.Error())
	}
	Db.Close()
	return activityType
}

func (d *DataSource) GetActivityTypeByUser(userID string) []bson.M {
	results := []bson.M{}
	query := bson.M{}
	if userID == "" {
		userID = "59b018561aa95138d0269de2"
	}
	query["author._id"] = bson.ObjectIdHex(userID)
	
	pm := bson.M{
        "$match" :query,
	}
	
	pj := bson.M{
        "$project": bson.M {"week": bson.M{"$isoWeek": "$time"}, "year": bson.M{"$year": "$time"}, "activity_type": 1},
	}
	
	pp := bson.M{
        "$group": bson.M {
			//"_id": bson.M{"$week": "$time" },
			"_id": bson.M{"year":"$year", "week": "$week", "activity_type": "$activity_type.name"},
			//"activity_types": bson.M{"$push": "$activity_type.name" },
			"total": bson.M{"$sum":1},
			//"count": bson.M{"$sum": 1},
		},
	}
	
	pg := bson.M{
        "$group": bson.M {
			"_id": bson.M{"week":"$_id.week", "year":"$_id.year"},
			"activity_types": bson.M{"$addToSet": bson.M{"name": bson.M{"$ifNull": []interface{}{"$_id.activity_type", "Unspecified"}}, "sum": "$total"}},
		},
	}
	
	pe := bson.M{
		"$sort": bson.D {
			bson.DocElem{Name: "_id.year", Value: 1},
			bson.DocElem{Name: "_id.week", Value: 1},
			
		},
	}

	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminBusinessChat)
	
	pipe := c.Pipe([]bson.M{pm, pj, pp, pg, pe})
	if err := pipe.All(&results); err != nil {	
		log.Printf(err.Error())
	}
	return results
}

func (d *DataSource) GetAllAdmins() []bson.M {
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminUser)
	admins := []bson.M{}
	
	if err := c.Find(nil).Select(bson.M{"username": 1, "owned": 1}).All(&admins); err != nil {
		log.Printf(err.Error())
	}
	Db.Close()
	return admins
}

func (d *DataSource) GetBusinesses() []bson.M {
	var businesses []bson.M
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C("businesses")
	if err := c.Find(nil).Select(bson.M{"name": 1}).All(&businesses); err != nil {	
		log.Printf(err.Error())
	}
	
	Db.Close()
	return businesses
}

func (d *DataSource) GetOwnedBusinessesID(adminIDHex string) bson.M {
	adminID := bson.ObjectIdHex(adminIDHex)
	var businesses bson.M
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminUser)
	if err := c.Find(bson.M{"_id": adminID}).Select(bson.M{"owned": 1, "_id":0}).One(&businesses); err != nil {	
		log.Printf(err.Error())
	}
	
	Db.Close()
	fmt.Println(businesses)
	return businesses
}

func (d *DataSource) UpdateOwner(businessesHex []string, adminIDHex string) bool {
	ok := true
	adminID := bson.ObjectIdHex(adminIDHex)
	count := len(businessesHex)
	businesses := make([]bson.ObjectId, count, count)
	for i, value := range businessesHex {
		businesses[i] = bson.ObjectIdHex(value)
	}
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(adminUser)
	if err := c.Update(bson.M{"_id":adminID}, bson.M{"$set": bson.M{"owned": businesses}}); err != nil {
		log.Printf(err.Error())
	}
	c = Db.C(businessCol)
	

	if _, err := c.UpdateAll(bson.M{"owner": adminID}, bson.M{"$set": bson.M{"owner": ""}}); err != nil {
		ok = false
	}
	if _, err := c.UpdateAll(bson.M{"_id": bson.M{"$in": businesses}}, bson.M{"$set": bson.M{"owner": adminID}}); err != nil {
		ok = false
	}
	Db.Close()
	return ok
}

func (d *DataSource) UpdateAdminOwnerUsersPage(usersHex []string, businessesHex []string, adminIDHex string) bool {
	var err error
	var da  *mgo.ChangeInfo
	ok := true
	adminID := bson.ObjectIdHex(adminIDHex)
	count := len(usersHex)
	users := make([]bson.ObjectId, count, count)
	for i, value := range usersHex {
		users[i] = bson.ObjectIdHex(value)
	}
	count = len(businessesHex)
	businesses := make([]bson.ObjectId, count, count)
	for i, value := range businessesHex {
		businesses[i] = bson.ObjectIdHex(value)
	}
	
	Db := db.MgoDb{}
	Db.Init()
	c := Db.C(businessCol)
	if _, err = c.UpdateAll(bson.M{"_id": bson.M{"$in": businesses}}, bson.M{"$set": bson.M{"owner": adminID}}); err != nil {
		ok = false
	}
	
	c = Db.C(adminUser)
	if da, err = c.UpdateAll(bson.M{"owned": bson.M{"$in": businesses}}, bson.M{"$pullAll": bson.M{"owned": businesses}}); err != nil {
		ok = false
	}
	
	if err := c.Update(bson.M{"_id":adminID}, bson.M{"$addToSet": bson.M{"owned": bson.M{"$each": businesses}}}); err != nil {
		ok = false
	}
	fmt.Println(da, businesses)
	return ok
	
}
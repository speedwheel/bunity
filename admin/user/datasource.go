package user

import(
	"app/model"
	"gopkg.in/mgo.v2/bson"
	"app/shared/db"
	"log"
	"net/url"
	//"github.com/kr/pretty"
	"strconv"
	"strings"
)

const (
	userCol = "users"
	businessCol = "businesses"
	businessCategCol = "categories"
)

var (
	userCols = []string{"", "firstname", "lastname", "email"}
	busCols = []string{"slug", "phone", "country", "website"}
)

type DataSource struct {
	Users []model.User
	User model.User
	Businesses []model.Business
	Db db.MgoDb
}

func NewDataSource() *DataSource {
	return &DataSource{
		Users: []model.User{},
		User: model.User{},
		Businesses: []model.Business{},
	}
}

func (d *DataSource) GetAllUsers(urlQuery url.Values) ([]model.User, int, int) {
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
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
			sliceBson = append(sliceBson, bson.M{"$or": []bson.M{bson.M{"firstname": bson.M{"$regex": "^"+w, "$options" : "i"}}, bson.M{"lastname": bson.M{"$regex": "^"+w, "$options" : "i"}},bson.M{"email": bson.M{"$regex": "^"+w, "$options" : "i"}}}})
		}
	
		query["$and"] = sliceBson
	}
	
	pm := bson.M{
        "$match" :query,
	}
	
	pp := bson.M{
        "$project": bson.M {"_id": 0},
	}
	
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
	
	pipe := c.Pipe([]bson.M{pm, po, pp, ps, pl })
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
			sliceBson = append(sliceBson, bson.M{"$or": []bson.M{bson.M{"nameSplit": bson.M{"$regex": "^"+w}}}})
		}
	
		query["$and"] = sliceBson
	}
	
	pm := bson.M{
        "$match" :query,
	}
	
	pp := bson.M{
        "$project": bson.M {"_id": 0},
	}
	
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
	
	pipe := c.Pipe([]bson.M{pm, po, pp, ps, pl })
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
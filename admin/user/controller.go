package user

import(
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	"strings"
	data "app/controller"
)


const (
	PathHome = "admin/home"
	PathUsers = "admin/user/users"
	PathBusinesses = "admin/user/businesses"
)

type Controller struct {
	iris.SessionController
	Source *DataSource
}

func (c *Controller) BeginRequest(ctx iris.Context) {
	c.SessionController.BeginRequest(ctx)
	c.Data["Path"] = ctx.Path()
	//ctx.GetCurrentRoute().StaticPath()
}

func (c *Controller) Get() {
	c.Tmpl = PathHome+".html"
}

func (c *Controller) GetUsers() {
	c.Data["countries"] = data.Countries
	c.Data["statesUSA"] = data.StatesUSA
	c.Data["statesCanada"] = data.StatesCanada
	c.Data["statesAustralia"] = data.StatesAustralia
	c.Data["businessCategs"] = c.Source.GetAllBusinessCategories()
	
	c.Data["yearsBusiness"] = data.YearsBusiness
	c.Data["numberEmployees"] = data.NrEmployees
	c.Data["sizeBusiness"] = data.SizeBusiness
	c.Data["relationshipBusiness"] = data.RelationshipBusiness
	c.Tmpl = PathUsers+".html"
}

func (c *Controller) GetBusinesses() {
	c.Tmpl = PathBusinesses+".html"
}

func (c *Controller) GetUserlist() {
	if c.Ctx.IsAjax() {
		urlQuery := c.Ctx.Request().URL.Query()
		draw := urlQuery["draw"][0]
		Data, CountFiltered, Count := c.Source.GetAllUsers(urlQuery)
		c.Ctx.JSON(map[string]interface{}{"draw": draw, "recordsTotal": Count, "recordsFiltered": CountFiltered, "data": Data})
	}
}

func (c *Controller) GetBusinesslist() {
	if c.Ctx.IsAjax() {
		urlQuery := c.Ctx.Request().URL.Query()
		draw := urlQuery["draw"][0]
		Data, CountFiltered, Count := c.Source.GetAllBusinesses(urlQuery)
		c.Ctx.JSON(map[string]interface{}{"draw": draw, "recordsTotal": Count, "recordsFiltered": CountFiltered, "data": Data})
	}
}

func (c *Controller) GetBusinessBy(businessID string) {
	if c.Ctx.IsAjax() {	
		business := c.Source.GetBusinessByID(businessID)
		c.Ctx.JSON(map[string]interface{}{"business": business})
	}
}

func (c *Controller) PostBusinessBy(businessID string) {
	if c.Ctx.IsAjax() {
		business := bson.M{}
		user := bson.M{}
		name := c.Ctx.FormValue("b[name]")
		nameSplit := strings.Fields(strings.ToLower(name))
		
		business["name"] = name
		business["namearr"] = nameSplit
		business["phone"] = c.Ctx.FormValue("b[phone]")
		business["address.add"] = c.Ctx.FormValue("b[country]")
		business["address.area"] = c.Ctx.FormValue("b[area]")
		business["address.state"] = c.Ctx.FormValue("b[state]")
		business["address.city"] = c.Ctx.FormValue("b[city]")
		business["address.country"] = c.Ctx.FormValue("b[country]")
		business["address.add"] = c.Ctx.FormValue("b[address]")
		business["address.add2"] = c.Ctx.FormValue("b[address2]")
		business["address.postalcode"] = c.Ctx.FormValue("b[postal_code]")
		business["web"] = c.Ctx.FormValue("b[website]")
		business["categ"] = c.Ctx.FormValue("b[businesscateg]")
		business["categ2"] = c.Ctx.FormValue("b[businesscateg2]")
		
		business["ybuss"] = c.Ctx.FormValue("b[yearsBusiness]")
		business["nemp"] = c.Ctx.FormValue("b[numberEmployees]")
		business["sbuss"] = c.Ctx.FormValue("b[sizeBusiness]")
		business["relbuss"] = c.Ctx.FormValue("b[relationshipBusiness]")
		
		business["desc"] = c.Ctx.FormValue("b[businessdescription]")
		
		business["social.facebook"] = c.Ctx.FormValue("b[businessFacebook]")
		business["social.google"] = c.Ctx.FormValue("b[businessGoogle]")
		business["social.instagram"] = c.Ctx.FormValue("b[businessInstagram]")
		business["social.youtube"] = c.Ctx.FormValue("b[businessYoutube]")
		business["social.pinterest"] = c.Ctx.FormValue("b[businessPinterest]")
		business["social.linkedin"] = c.Ctx.FormValue("b[businessLinkedin]")
		business["social.twitter"] = c.Ctx.FormValue("b[businessTwitter]")
		
		
		user["business.$.name"] = name
		user["business.$.namearr"] = nameSplit
		userID := c.Ctx.FormValue("b[userid]")
		success := c.Source.UpdateBusinessByID(businessID, business, user, userID)
		c.Ctx.JSON(map[string]interface{}{"success": success})
	}
}

func (c *Controller) DeleteBusinessBy(/*userID string,*/ businessID string)  {
	res := false
	if c.Ctx.IsAjax() {
		res = c.Source.DeleteBusinessByID(businessID/*, userID*/)
	}
	c.Ctx.JSON(map[string]bool{"success": res})
}
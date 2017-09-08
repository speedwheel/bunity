package user

import(
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	"strings"
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
		businessCategs := c.Source.GetAllBusinessCategories()
		c.Ctx.JSON(map[string]interface{}{"business": business, "businessCateg": businessCategs})
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
		business["address.country"] = c.Ctx.FormValue("b[country]")
		business["address.area"] = c.Ctx.FormValue("b[area]")
		business["address.city"] = c.Ctx.FormValue("b[city]")
		business["address.address"] = c.Ctx.FormValue("b[address]")
		
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
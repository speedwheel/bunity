package user

import(
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"app/model"
	data "app/controller"
	"time"
)


const (
	PathHome = "admin/home"
	PathUsers = "admin/user/users"
	PathUserSettings = "admin/user/settings"
	PathUserStatistics = "admin/user/statistics"
)

type Controller struct {
	iris.SessionController
	Source *DataSource
	//Ws *Ws
}

func (c *Controller) BeginRequest(ctx iris.Context) {
	c.SessionController.BeginRequest(ctx)
	//c.Ws.Conn.OnConnection(c.Ws.BusinessChatNotif)
	c.Data["Path"] = ctx.Path()
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
	c.Data["userName"] = c.Session.Get("user").(model.User).Firstname //+ " " + c.Session.Get("user").(model.User).Lastname
	c.Data["userID"] = c.Session.Get("user").(model.User).Id
	c.Data["timeNow"] = time.Now()
	c.Data["activityType"] = c.Source.GetAllActivities()
	c.Tmpl = PathUsers+".html"
}

func (c *Controller) GetUsersSettings() {
	c.Tmpl = PathUserSettings+".html"
}

func (c *Controller) GetUsersStatistics() {
	results := c.Source.GetActivityTypeByUser("")
	c.Data["activityWeek"] = results
	c.Data["adminList"] = c.Source.GetAllAdmins()
	c.Tmpl = PathUserStatistics+".html"
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
		comments := c.Source.GetCommentsByBusiness(businessID)
		c.Ctx.JSON(map[string]interface{}{"business": business, "comments": comments})
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

func (c *Controller) PostPictureDelete() {
	if c.Ctx.IsAjax() {
		businessID := c.Ctx.FormValue("businessID")
		userID := c.Ctx.FormValue("userID")
		imageType := c.Ctx.FormValue("imageType")
		fileID := c.Ctx.FormValue("fileID")
		success := c.Source.DeletePicture(businessID, userID, fileID, imageType)
		c.Ctx.JSON(map[string]bool{"success": success})
	}
}

func (c *Controller) PostPictureAddBy(businessID string) {
	userID := c.Ctx.FormValue("userID")
	imageFormat := c.Ctx.FormValue("imageFormat")
	imageType := c.Ctx.FormValue("imageType")
	file, _, _ := c.Ctx.FormFile("file")
	fname, url := c.Source.AddPicture(userID, businessID, imageFormat, imageType, file)
	
	c.Ctx.JSON(map[string]string{"fname": fname, "url": url})
}

func (c *Controller) PostBusinessComment() {
	if c.Ctx.IsAjax() {
		user := c.Session.Get("user").(model.User)
		chatMessage := ChatMessage{}
		chatMessage.Id = bson.NewObjectId()
		chatMessage.Author.Id = user.Id
		chatMessage.Text = c.Ctx.FormValue("msg")
		chatMessage.Author.Name = user.Firstname + " " + user.Lastname
		chatMessage.BusinessID = bson.ObjectIdHex(c.Ctx.FormValue("bizID"))
		chatMessage.Time = time.Now()
		if c.Ctx.FormValue("parentID") != "" {
			chatMessage.ParentID = bson.ObjectIdHex(c.Ctx.FormValue("parentID"))
		}
		activityTypeID := c.Ctx.FormValue("activityTypeID")
		if activityTypeID != "" {
			chatMessage.ActivityType.Id = bson.ObjectIdHex(activityTypeID)
			chatMessage.ActivityType.Name = c.Ctx.FormValue("activityTypeName")
		}
		response := c.Source.InsertBusinessComment(chatMessage)
		if response == true {
			c.Ctx.JSON(map[string]interface{}{"success": response, "data": chatMessage})
			return
		}
		c.Ctx.JSON(map[string]bool{"success": response})
	}
}

func (c *Controller) PutBusinessComment() {
	if c.Ctx.IsAjax() {
		user := c.Session.Get("user").(model.User)
		chatID := c.Ctx.FormValue("postID")
		msg := c.Ctx.FormValue("msg")
		response := c.Source.UpdateBusinessCommentByID(chatID,user.Id, msg)
		
		c.Ctx.JSON(map[string]bool{"success": response})
	}
}

func (c *Controller) PostActivitytype() {
	if c.Ctx.IsAjax() {
		activity := ActivityType{}
		activity.Id = bson.NewObjectId()
		activity.Name = c.Ctx.FormValue("name")
		
		response := c.Source.AddNewActivity(activity)
		c.Ctx.JSON(map[string]bool{"success": response})
	}
}

func (c *Controller) GetActivitylist() {
	if c.Ctx.IsAjax() {
		urlQuery := c.Ctx.Request().URL.Query()
		draw := urlQuery["draw"][0]
		Data, CountFiltered, Count := c.Source.GetAllActivitiesTD(urlQuery)
		c.Ctx.JSON(map[string]interface{}{"draw": draw, "recordsTotal": Count, "recordsFiltered": CountFiltered, "data": Data})
	}
}

func (c *Controller) PutActivitylist() {
	if c.Ctx.IsAjax() {
		activity := ActivityType{}
		activity.Id = bson.ObjectIdHex(c.Ctx.FormValue("pk"))
		activity.Name = c.Ctx.FormValue("value")
		response := c.Source.UpdateActivityType(activity)
		c.Ctx.JSON(map[string]interface{}{"success": response, "value": activity.Name})
	}
}

func (c *Controller) DeleteActivitytypeBy(activity string) {
	if c.Ctx.IsAjax() {
		activityID := bson.ObjectIdHex(activity)
		response := c.Source.DeleteActivityType(activityID)
		c.Ctx.JSON(map[string]interface{}{"success": response})
	}
}

func (c *Controller) GetActivitytypeBy(userID string) {
	results := c.Source.GetActivityTypeByUser(userID)
	c.Ctx.JSON(map[string]interface{}{"data": results})
}


func (c *Controller) GetTestjson() {
	c.Ctx.JSON(map[string]interface{}{
		"my_statistic_1": map[string]interface{}{ "type": "integer", "value": 1, "label": "My Statistic 1", "order": 0},
		"my_statistic_2": map[string]interface{}{ "type": "percentage", "value": 0.5, "label": "My Statistic 2", "order": 1 },
		"my_statistic_3": map[string]interface{}{ "type": "percentage", "value": 0.25, "label": "My Statistic 3", "order": 2 },
		"my_statistic_4": map[string]interface{}{ "type": "percentage", "value": 0.75, "label": "My Statistic 4", "order": 3 },
	})
}
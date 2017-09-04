package user

import(
	"github.com/kataras/iris"
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
package admin

type (
	AdminItem struct {
		Name string
		Url string
		Icon string
		SubMenu []AdminSubItem 
	}
	
	AdminSubItem struct {
		Name string
		Url string
	}
)


func NewAdminMenu() *[]AdminItem {
	return &[]AdminItem{
		AdminItem {
			Name: "Data",
			Icon: "folder",
			SubMenu: []AdminSubItem {
				AdminSubItem {
					Name: "Users",
					Url: "/users",
				},
			},
		},
		AdminItem {
			Name: "Log Out",
			Url: "/logout",
			Icon: "settings_power",
		},
	}
}
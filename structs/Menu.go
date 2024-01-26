package structs

import "fmt"

type MenuItem struct {
	Id       int
	Menu     string
	Url      string
	IsActive bool
}

type MenuList struct {
	Menus []MenuItem
}

func (m *MenuList) SetActive(id int) {
	for i := range m.Menus {
		if m.Menus[i].Id == id {
			m.Menus[i].IsActive = true
		} else {
			m.Menus[i].IsActive = false
		}
	}
}

func PopulateMenu() MenuList {
	var menu MenuList
	menu.Menus = append(menu.Menus, MenuItem{Id: 0, Menu: "recurring bills", Url: "/recurring-debts", IsActive: true})
	menu.Menus = append(menu.Menus, MenuItem{Id: 1, Menu: "recurring debts", Url: "/recurring-debts", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 2, Menu: "long term debts", Url: "/long-term-debts", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 3, Menu: "assets", Url: "/assets", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 4, Menu: "credit utilization", Url: "/credit-utilization", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 5, Menu: "goals", Url: "/goals", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 6, Menu: "recomendations", Url: "/recomendations", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 7, Menu: "calendar", Url: "/calendar", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 8, Menu: "drip calculator", Url: "/drip-calculator", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 9, Menu: "time tables", Url: "/time-tables", IsActive: false})
	menu.Menus = append(menu.Menus, MenuItem{Id: 10, Menu: "Settings", Url: "/settings", IsActive: false})
	return menu
}

var Menu MenuList

func init() {
	fmt.Println("Menu.init(): Populating Menu...")
	Menu = PopulateMenu()
}

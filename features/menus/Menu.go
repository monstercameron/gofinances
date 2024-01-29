package menus

import (
	"database/sql"
	"fmt"
	"github.com/monstercameron/gofinances/database"
	"log"
)

type MenuItem struct {
	Id       int
	Menu     string
	Url      string
	IsActive bool
}

func (m *MenuItem) Save() {
	fmt.Println("MenuItem.Save()")
	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	var id int
	query := `SELECT id FROM menus WHERE id=?;`
	err = tx.QueryRow(query, m.Id).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		log.Fatal(err)
	}

	if err == sql.ErrNoRows {
		query = `INSERT INTO menus (menu, url, is_active) VALUES (?, ?, ?);`
		_, err = tx.Exec(query, m.Menu, m.Url, m.IsActive)
	} else {
		query = `UPDATE menus SET menu=?, url=?, is_active=? WHERE id=?;`
		_, err = tx.Exec(query, m.Menu, m.Url, m.IsActive, m.Id)
	}

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func SetActiveMenu(id int) int {
	var err error
	// set all menu items to inactive
	// then set the one with the id to active
	_, err = database.DB.Exec("UPDATE menus SET is_active=0")
	if err != nil {
		return -1
	} else {
		fmt.Println("Menu.SetActive(): \t\tAll menu items set to inactive.")
	}
	_, err = database.DB.Exec("UPDATE menus SET is_active=1 WHERE id=?", id)
	if err != nil {
		return -1
	} else {
		fmt.Println("Menu.SetActive(): \t\tMenu item set to active.")
	}
	return id
}

// return the id number of the first active menu item
func GetActiveMenu() int {
	var id int
	query := `SELECT id FROM menus WHERE is_active=1 limit 1;`
	err := database.DB.QueryRow(query).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		}
		log.Fatal(err)
	}
	return id
}

func PopulateMenu() {
	menus := []MenuItem{}
	menus = append(menus, MenuItem{Id: 0, Menu: "recurring bills", Url: "/recurring-debts", IsActive: true})
	menus = append(menus, MenuItem{Id: 1, Menu: "short term debts", Url: "/short-term-debts", IsActive: false})
	menus = append(menus, MenuItem{Id: 2, Menu: "long term debts", Url: "/long-term-debts", IsActive: false})
	menus = append(menus, MenuItem{Id: 3, Menu: "assets", Url: "/assets", IsActive: false})
	menus = append(menus, MenuItem{Id: 4, Menu: "credit utilization", Url: "/credit-utilization", IsActive: false})
	menus = append(menus, MenuItem{Id: 5, Menu: "goals", Url: "/goals", IsActive: false})
	menus = append(menus, MenuItem{Id: 6, Menu: "recomendations", Url: "/recomendations", IsActive: false})
	menus = append(menus, MenuItem{Id: 7, Menu: "calendar", Url: "/calendar", IsActive: false})
	menus = append(menus, MenuItem{Id: 8, Menu: "drip calculator", Url: "/drip-calculator", IsActive: false})
	menus = append(menus, MenuItem{Id: 9, Menu: "time tables", Url: "/time-tables", IsActive: false})
	menus = append(menus, MenuItem{Id: 10, Menu: "Settings", Url: "/settings", IsActive: false})

	// loop and save
	for _, menu := range menus {
		menu.Save()
	}
}

func GetMenus() []MenuItem {
	fmt.Println("Menu.GetMenus()")
	menus := []MenuItem{}
	rows, err := database.DB.Query("SELECT id, menu, url, is_active FROM menus")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var menu MenuItem
		if err := rows.Scan(&menu.Id, &menu.Menu, &menu.Url, &menu.IsActive); err != nil {
			log.Fatal(err)
		}
		// fmt.Println(menu)
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return menus
}

func init() {
	fmt.Println("Menu.init(): \t\t\tPopulating Menu...")
	// Menu = PopulateMenu()

	// Create table for menus
	var err error
	_, err = database.DB.Exec("CREATE TABLE IF NOT EXISTS menus (id INTEGER PRIMARY KEY, menu TEXT, url TEXT, is_active INTEGER)")
	if err != nil {
		log.Fatalf("Failed to create 'recurring_bills' table: %v", err)
	} else {
		fmt.Println("Database.Init(): \t\t'recurring_bills' table created.")
	}

	// check if menus are empty
	var count int
	query := `SELECT count(*) FROM menus;`
	err = database.DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		fmt.Println("Menu.init(): \t\t\tPopulating Menu...")
		PopulateMenu()
	}
}

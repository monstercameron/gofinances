package menus

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/monstercameron/gofinances/database"
	"github.com/monstercameron/gofinances/features/bills"
	"github.com/monstercameron/gofinances/features/debt"
	"github.com/monstercameron/gofinances/features/assets"
	"github.com/monstercameron/gofinances/features/settings"
	"github.com/monstercameron/gofinances/helpers"
)

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

// GetActiveMenu retrieves the ID of the active menu from the database.
// It executes a SQL query to select the ID of the menu where is_active is set to 1.
// If no active menu is found, it returns -1.
// If an error occurs during the query execution, it logs the error and terminates the program.
// The retrieved ID is returned as an integer.
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

// PopulateMenu populates the menu with predefined menu items.
// It creates a slice of MenuItem structs and appends each item to the slice.
// Then, it loops through the slice and saves each menu item.
func PopulateMenu() {
	menus := []MenuItem{}
	menus = append(menus, MenuItem{Id: 1, Menu: "recurring bills", Url: "/bills", IsActive: true})
	menus = append(menus, MenuItem{Id: 2, Menu: "short term debts", Url: "/debts", IsActive: false})
	menus = append(menus, MenuItem{Id: 3, Menu: "assets", Url: "/assets", IsActive: false})
	menus = append(menus, MenuItem{Id: 4, Menu: "credit utilization", Url: "/credit", IsActive: false})
	menus = append(menus, MenuItem{Id: 5, Menu: "goals", Url: "/goals", IsActive: false})
	menus = append(menus, MenuItem{Id: 6, Menu: "recomendations", Url: "/recomendations", IsActive: false})
	menus = append(menus, MenuItem{Id: 7, Menu: "calendar", Url: "/calendar", IsActive: false})
	menus = append(menus, MenuItem{Id: 8, Menu: "drip calculator", Url: "/drip", IsActive: false})
	menus = append(menus, MenuItem{Id: 9, Menu: "time tables", Url: "/timetables", IsActive: false})
	menus = append(menus, MenuItem{Id: 10, Menu: "Settings", Url: "/settings", IsActive: false})

	// loop and save
	for _, menu := range menus {
		menu.Save()
	}
}

// GetMenus retrieves a list of menus from the database.
// It executes a SQL query to fetch the menus and their details,
// and returns a slice of MenuItem structs representing the menus.
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

// MenuPicker handles the HTTP request for menu selection.
// It extracts the menu ID from the URL path, sets it as the active menu,
// and renders the main menu component to the response writer.
// The menu ID should be provided as a URL parameter in the format "/menu/{id}".
// If the URL parameter is invalid or cannot be converted to an integer, it returns a bad request error.
func GetMenu(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MenuPicker.MenuPicker(): r.URL.Path: ", r.URL.Path)

	// example /menu/1
	urlParam, err := helpers.ExtractSegmentFromPath(r.URL.Path, 2)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// convert to int
	id, err := strconv.Atoi(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	SetActiveMenu(id)

	component := MainMenuComponent(GetMenus())
	// serve text/html
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "menuSwitch")
	// render the component to the response writer
	component.Render(r.Context(), w)
}

// GetTab is a function that handles the logic for serving different tabs in a menu.
// It takes in an http.ResponseWriter and an http.Request as parameters.
// It retrieves the active menu ID and based on the ID, it renders the corresponding component.
// The rendered component is then served as a response with the appropriate Content-Type header.
// If the active menu ID is invalid, it returns a Bad Request error.
func GetTab(w http.ResponseWriter, r *http.Request) {
	id := GetActiveMenu()
	if id == -1 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	fmt.Println("MenuPicker.GetTab(): id: ", id)

	switch id {
	case 1:
		component := bills.RecurringBillsIndex()
		// serve text/html
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
		return
	case 2:
		component := debt.DebtsIndex()
		// send string response
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
		return
	case 3:
		component := assets.AssetsIndex()
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		component.Render(r.Context(), w)
		return
	case 4:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Credit utilization"))
		return
	case 5:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Goals"))
		return
	case 6:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Recommendations"))
		return
	case 7:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Calendar"))
		return
	case 8:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Drip calculator"))
		return
	case 9:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("time tables"))
		return
	case 10:
		users, err := settings.GetAllSettingsUsers()
		if err != nil {
			http.Error(w, "Error getting users", http.StatusInternalServerError)
			return
		}
		component := settings.SettingsPageIndex(users)
		w.Header().Set("Content-Type", "text/plain")
		component.Render(r.Context(), w)
		return
	default:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Default"))
		return
	}
}

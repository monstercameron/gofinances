package controller

import (
	"fmt"
	"net/http"

	"github.com/monstercameron/gofinances/features/bills"
	"github.com/monstercameron/gofinances/features/home"
	"github.com/monstercameron/gofinances/features/menus"
	"github.com/monstercameron/gofinances/features/settings"
)

func init() {
	fmt.Println("controller.init(): \t\tInitializing controller...")
}

// CreateRoutes sets up the routes for the HTTP server.
func CreateRoutes(server *http.ServeMux) {
	// Serve static files from the 'views/static' directory
	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up routes and associate them with handler functions
	/////////////////////// HOME ROUTES ///////////////////////
	server.HandleFunc("/", home.GetIndexPage) // Handler for the root route
	/////////////////////// MENU ROUTES ///////////////////////
	server.HandleFunc("/menu/", menus.MenuPicker) // Handler for the '/menu/' route
	server.HandleFunc("/pane", menus.GetTab)      // Handler for the '/pane' route
	/////////////////////// DEBT ROUTES ///////////////////////
	server.HandleFunc("/debts", bills.GetBillList)              // Handler for retrieving bills
	server.HandleFunc("/debts/new", bills.AddBills)             // Handler for adding new bills
	server.HandleFunc("/debts/update", bills.UpdateBills)       // Handler for updating bills
	server.HandleFunc("/debts/delete", bills.DeleteBills)       // Handler for deleting bills
	server.HandleFunc("/debts/total", bills.GetBillsTotalDebts) // Handler for retrieving total debts
	/////////////////////// SETTINGS ROUTES ///////////////////////
	server.HandleFunc("/settings", settings.GetSettingsPage)
	server.HandleFunc("/settings/test", settings.GetSettingsPage)         // Handler for the '/settings' route
	server.HandleFunc("/settings/adduser", settings.GetSettingsUserInput) // Handler for the '/settings/save' route
	server.HandleFunc("/settings/getusers", settings.GetSettingsUser)     // Handler for the '/settings/save' route
}

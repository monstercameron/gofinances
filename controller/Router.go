package controller

import (
	"fmt"
	"github.com/monstercameron/gofinances/features/home"
	"github.com/monstercameron/gofinances/features/menus"
	"github.com/monstercameron/gofinances/features/monthlydebts"
	"github.com/monstercameron/gofinances/features/settings"
	"net/http"
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
	server.HandleFunc("/debts", monthlydebts.GetBillList)              // Handler for retrieving bills
	server.HandleFunc("/debts/new", monthlydebts.AddBills)             // Handler for adding new bills
	server.HandleFunc("/debts/update", monthlydebts.UpdateBills)       // Handler for updating bills
	server.HandleFunc("/debts/delete", monthlydebts.DeleteBills)       // Handler for deleting bills
	server.HandleFunc("/debts/total", monthlydebts.GetBillsTotalDebts) // Handler for retrieving total debts
	/////////////////////// SETTINGS ROUTES ///////////////////////
	server.HandleFunc("/settings/test", settings.GetSettingsPage) // Handler for the '/settings' route
}

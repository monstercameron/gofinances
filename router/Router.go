package router

import (
	"fmt"
	"net/http"
	"github.com/monstercameron/gofinances/features/bills"
	"github.com/monstercameron/gofinances/features/home"
	"github.com/monstercameron/gofinances/features/menus"
	"github.com/monstercameron/gofinances/features/settings"
	// "github.com/monstercameron/gofinances/features/debts"
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
	server.HandleFunc("/bills", home.GetIndexPage) 		
	server.HandleFunc("/bills/bills", bills.GetBillList)              // Handler for retrieving bills
	server.HandleFunc("/bills/new", bills.AddBills)             // Handler for adding new bills
	server.HandleFunc("/bills/update", bills.UpdateBills)       // Handler for updating bills
	server.HandleFunc("/bills/delete", bills.DeleteBills)       // Handler for deleting bills
	server.HandleFunc("/bills/total", bills.GetBillsTotalDebts) // Handler for retrieving total debts
	/////////////////////// SETTINGS ROUTES ///////////////////////
	server.HandleFunc("/settings", home.GetIndexPage)
	server.HandleFunc("/settings/user", settings.GetSettingsUserActions) 	// Handler for the '/settings/save' route
	server.HandleFunc("/settings/users", settings.GetSettingsUser)     		// Handler for the '/settings/save' route
	/////////////////////// Debts ROUTES ///////////////////////
	server.HandleFunc("/debts", home.GetIndexPage)
}
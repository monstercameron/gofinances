package router

import (
	"fmt"
	"github.com/monstercameron/gofinances/features/bills"
	"github.com/monstercameron/gofinances/features/debts"
	"github.com/monstercameron/gofinances/features/home"
	"github.com/monstercameron/gofinances/features/menus"
	"github.com/monstercameron/gofinances/features/settings"
	"net/http"
)

func init() {
	fmt.Println("controller.init(): \t\tInitializing controller...")
}

// CreateRoutes sets up the routes for the HTTP server.
func CreateRoutes(server *http.ServeMux) *http.ServeMux {
	// Serve static files from the 'views/static' directory
	server.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Set up routes and associate them with handler functions
	/////////////////////// HOME ROUTES ///////////////////////
	server.HandleFunc("GET /", home.GetIndexPage) // Handler for the root route
	/////////////////////// MENU ROUTES ///////////////////////
	server.HandleFunc("GET /menu/", menus.GetMenu) // Handler for the '/menu/' route
	server.HandleFunc("GET /tab", menus.GetTab)      // Handler for the '/tab' route
	/////////////////////// DEBT ROUTES ///////////////////////
	server.HandleFunc("GET /bills", home.GetIndexPage)
	server.HandleFunc("GET /bills/bills", bills.GetBillList)        // Handler for retrieving bills
	server.HandleFunc("GET /bills/new", bills.AddBills)             // Handler for adding new bills
	server.HandleFunc("GET /bills/update", bills.UpdateBills)       // Handler for updating bills
	server.HandleFunc("GET /bills/delete", bills.DeleteBills)       // Handler for deleting bills
	server.HandleFunc("GET /bills/total", bills.GetBillsTotalDebts) // Handler for retrieving total debts
	/////////////////////// SETTINGS ROUTES ///////////////////////
	server.HandleFunc("GET /settings", home.GetIndexPage)
	server.HandleFunc("GET /settings/user", settings.GetSettingsUserActions) // Handler for the '/settings/save' route
	server.HandleFunc("GET /settings/users", settings.GetSettingsUser)       // Handler for the '/settings/save' route
	/////////////////////// Debts ROUTES ///////////////////////
	server.HandleFunc("GET /debts", home.GetIndexPage)
	// server.HandleFunc("/debts/add", debts.GetDebtItems)
	server.HandleFunc("GET /debts/update", debts.UpdateDebtItems)
	// server.HandleFunc("/debts/create", debts.GetDebtItems)
	// server.HandleFunc("/debts/create", debts.GetDebtItems)
	// server.HandleFunc("/debts/create", debts.GetDebtItems)
	return server
}

package router

import (
	"fmt"
	"net/http"

	"github.com/monstercameron/gofinances/features/bills"
	"github.com/monstercameron/gofinances/features/debt"
	"github.com/monstercameron/gofinances/features/home"
	"github.com/monstercameron/gofinances/features/menus"
	"github.com/monstercameron/gofinances/features/settings"
)

func init() {
	fmt.Println("controller.init(): \t\tInitializing controller...")
}

// CreateRoutes sets up the routes for the HTTP server.
func CreateRoutes(server *http.ServeMux) *http.ServeMux {
	// Serve static files from the 'views/static' directory
	server.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Set up routes and associate them with handler functions
	/////////////////////// HOME ROUTES /////////////////////////////////////////////////////////////////////////////
	server.HandleFunc("GET /", home.GetIndexPage) // Handler for the root route
	/////////////////////// MENU ROUTES /////////////////////////////////////////////////////////////////////////////
	server.HandleFunc("GET /menu/", menus.GetMenu) // Handler for the '/menu/' route
	server.HandleFunc("GET /tab", menus.GetTab)    // Handler for the '/tab' route
	/////////////////////// Bills ROUTES ////////////////////////////////////////////////////////////////////////////
	server.HandleFunc("GET /bills", home.GetIndexPage)
	server.HandleFunc("GET /bills/", bills.GetManyBills)                     // Handler for retrieving bills
	server.HandleFunc("GET /bills/{id}", bills.GetOneBill)                   // Handler for retrieving bills
	server.HandleFunc("POST /bills/", bills.AddBills)                        // Handler for adding new bills
	server.HandleFunc("GET /bills/add", bills.GetAddBillingComponent)        // Handler for updating bills
	server.HandleFunc("GET /bills/edit/{id}", bills.GetEditBillingComponent) // Handler for updating bills
	server.HandleFunc("POST /bills/{id}", bills.UpdateBills)                 // Handler for updating bills
	server.HandleFunc("DELETE /bills/{id}", bills.DeleteBills)               // Handler for deleting bills
	server.HandleFunc("GET /bills/total", bills.GetBillsTotalDebts)          // Handler for retrieving total debts
	/////////////////////// SETTINGS ROUTES ///////////////////////////////////////////////////////////////////////////
	server.HandleFunc("GET /settings", home.GetIndexPage)
	server.HandleFunc("GET /settings/user", settings.GetSettingsUserActions) // Handler for the '/settings/save' route
	server.HandleFunc("GET /settings/users", settings.GetSettingsUser)       // Handler for the '/settings/save' route
	/////////////////////// Debts ROUTES /////////////////////////////////////////////////////////////////////////////
	server.HandleFunc("GET /debts", home.GetIndexPage)
	// server.HandleFunc("/debts/add", debt.GetDebtItems)
	server.HandleFunc("GET /debts/update", debt.UpdateDebtItems)
	// server.HandleFunc("/debts/create", debt.GetDebtItems)
	// server.HandleFunc("/debts/create", debt.GetDebtItems)
	// server.HandleFunc("/debts/create", debt.GetDebtItems)
	return server
}

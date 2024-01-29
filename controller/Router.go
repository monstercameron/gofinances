package controller

import (
	"fmt"
	"net/http"

	"github.com/monstercameron/gofinances/features/home"
	"github.com/monstercameron/gofinances/features/menus"
	"github.com/monstercameron/gofinances/features/monthlydebts"
)

// CreateRoutes sets up the routes for the HTTP server.
func CreateRoutes(server *http.ServeMux) {
	// Serve static files from the 'views/static' directory
	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	// Set up routes and associate them with handler functions
	server.HandleFunc("/", home.GetIndexPage)                          // Handler for the root route
	server.HandleFunc("/menu/", menus.MenuPicker)                      // Handler for the '/menu/' route
	server.HandleFunc("/pane", menus.GetTab)                         // Handler for the '/pane' route
	server.HandleFunc("/debts", monthlydebts.GetBillList)              // Handler for retrieving bills
	server.HandleFunc("/debts/new", monthlydebts.AddBills)             // Handler for adding new bills
	server.HandleFunc("/debts/update", monthlydebts.UpdateBills)       // Handler for updating bills
	server.HandleFunc("/debts/delete", monthlydebts.DeleteBills)       // Handler for deleting bills
	server.HandleFunc("/debts/total", monthlydebts.GetBillsTotalDebts) // Handler for retrieving total debts
}

func init() {
	fmt.Println("controller.init(): \t\tInitializing controller...")
}

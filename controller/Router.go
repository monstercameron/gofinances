package controller

import (
	"net/http"
)

// CreateRoutes sets up the routes for the HTTP server.
func CreateRoutes(server *http.ServeMux) {
    // Serve static files from the 'views/static' directory
    server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

    // Set up routes and associate them with handler functions
    server.HandleFunc("/", GetIndexPage)        		// Handler for the root route
    server.HandleFunc("/menu/", MenuPicker)     		// Handler for the '/menu/' route
    server.HandleFunc("/pane", TestPage)        		// Handler for the '/pane' route
    server.HandleFunc("/debts", GetBills)       		// Handler for retrieving bills
    server.HandleFunc("/debts/new", AddBills)   		// Handler for adding new bills
    server.HandleFunc("/debts/update", UpdateBills) 	// Handler for updating bills
    server.HandleFunc("/debts/delete", DeleteBills) 	// Handler for deleting bills
    server.HandleFunc("/debts/total", GetTotalDebts) 	// Handler for deleting bills
}
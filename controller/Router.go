package controller

import (
	"fmt"
	"net/http"
)

func CreateRoutes(server *http.ServeMux) {
	// Set up your HTTP handlers
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	// Define a handler for a specific route
	server.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've hit the example route!")
	})
}

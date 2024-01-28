package main

import (
	"fmt"
	"github.com/monstercameron/gofinances/controller"
	"github.com/monstercameron/gofinances/database"
	"github.com/monstercameron/gofinances/helpers"
	"net/http"
)

func main() {
	fmt.Println("main()")
	// Initialize database
	database.SimpleTest(database.DB)
	
	// Create a new HTTP server
	server := http.NewServeMux()

	// Set up HTTP routes
	controller.CreateRoutes(server)

	// Setup signal handling and receive shutdown signal
	httpServer, done := helpers.SetupSignalHandling(server)

	// Start the server in a goroutine
	go func() {
		// Start the server and check for errors
		fmt.Println("Starting server on :3000")
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe: %v\n", err)
		}
	}()
	fmt.Println("Server started on port 3000")

	// Block until a shutdown signal is received
	<-done
	fmt.Println("Server gracefully stopped")
}

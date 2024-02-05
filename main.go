package main

import (
	"fmt"
	"github.com/monstercameron/gofinances/router"
	"github.com/monstercameron/gofinances/database"
	"github.com/monstercameron/gofinances/helpers"
	"net/http"
)

func main() {
	// Initialize database
	database.SimpleTest(database.DB)
	
	// Create a new HTTP server
	server := http.NewServeMux()

	// Set up HTTP routes
	router.CreateRoutes(server)

	// Setup signal handling and receive shutdown signal
	httpServer, done := helpers.SetupSignalHandling(server)

	// Start the server in a goroutine
	go func() {
		// Start the server and check for errors
		fmt.Println("\nMain: \t\t\t\tStarting server...")
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe: %v\n", err)
		}
	}()
	fmt.Println("Main: \t\t\t\tServer started.")

	// Block until a shutdown signal is received
	<-done
	fmt.Println("Main: \t\t\t\tShutting down server...")
}

func init() {
	fmt.Println("main.init(): \t\t\tInitializing main...")
}
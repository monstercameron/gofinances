package home

import (
	"net/http"
)

// GetIndexPage handles HTTP GET requests for the index page.
func GetIndexPage(w http.ResponseWriter, r *http.Request) {
	// Uncomment the following line to log incoming requests
	// helpers.Slogger.Info("Received request", "method", r.Method, "url", r.URL.String(), "protocol", r.Proto)

	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create the index page component
	component := IndexPage("goFinances")

	// Set the Content-Type of the response to text/html
	w.Header().Set("Content-Type", "text/html")

	// Render the index page component to the response writer
	component.Render(r.Context(), w)
}

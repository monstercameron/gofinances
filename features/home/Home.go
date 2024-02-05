package home

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/monstercameron/gofinances/features/menus"
	"github.com/monstercameron/gofinances/features/bills"
	"github.com/monstercameron/gofinances/features/settings"
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

func GetStartingPage() templ.Component {

	id := menus.GetActiveMenu()
	if id == -1 {
		fmt.Println("home.GetStartingPage(): id == -1")
		return bills.RecurringBillsIndex()
	}

	switch id {
	case 0:
		return bills.RecurringBillsIndex()
	case 10:
		return settings.SettingsPageIndex(settings.GetAllSettingsUsers())
	default:
		fmt.Println("home.GetStartingPage(): id == -1")
		return bills.RecurringBillsIndex()
	}
}

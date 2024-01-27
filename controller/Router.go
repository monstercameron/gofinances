package controller

import (
	"net/http"

	"github.com/monstercameron/gofinances/structs"
	"github.com/monstercameron/gofinances/views/pages"
)

func CreateRoutes(server *http.ServeMux) {
	// Set up your HTTP handlers
	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	server.HandleFunc("/", GetIndexPage)
	server.HandleFunc("/menu/", MenuPicker)
	server.HandleFunc("/pane", TestPage)
	server.HandleFunc("/debts", GetBills)
	server.HandleFunc("/debts/update", UpdateBills)
	server.HandleFunc("/debts/new", AddBills)
}

func GetIndexPage(w http.ResponseWriter, r *http.Request) {
	// helpers.Slogger.Info("Received request", "method", r.Method, "url", r.URL.String(), "protocol", r.Proto)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	component := pages.IndexPage("goFinances", structs.Menu)
	// serve text/html
	w.Header().Set("Content-Type", "text/html")
	// render the component to the response writer
	component.Render(r.Context(), w)
}

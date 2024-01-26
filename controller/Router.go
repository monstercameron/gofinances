package controller

import (
	"github.com/monstercameron/gofinances/structs"
	"github.com/monstercameron/gofinances/views/pages"
	"net/http"
	// "github.com/monstercameron/gofinances/views/components"
)

func CreateRoutes(server *http.ServeMux) {
	// Set up your HTTP handlers
	server.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	server.HandleFunc("/", GetIndexPage)
}

func GetIndexPage(w http.ResponseWriter, r *http.Request) {
	// helpers.Slogger.Info("Received request", "method", r.Method, "url", r.URL.String(), "protocol", r.Proto)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	component := pages.IndexPage("My Todo List", structs.Menu)
	// serve text/html
	w.Header().Set("Content-Type", "text/html")
	// render the component to the response writer
	component.Render(r.Context(), w)
}

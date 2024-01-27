package controller

import (
	"fmt"
	"github.com/monstercameron/gofinances/helpers"
	"github.com/monstercameron/gofinances/structs"
	"github.com/monstercameron/gofinances/views/components"
	"net/http"
	"strconv"
	// "github.com/monstercameron/gofinances/views/pages"
)

func init() {
	fmt.Println("MenuPicker.init(): starting MenuPicker route...")
}

func MenuPicker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("MenuPicker.MenuPicker(): r.URL.Path: ", r.URL.Path)

	// example /menu/1
	urlParam, err := helpers.ExtractSegmentFromPath(r.URL.Path, 2)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// convert to int
	id, err := strconv.Atoi(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	structs.Menu.SetActive(id)

	component := components.MainMenuComponent(structs.Menu.Menus)
	// serve text/html
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "menuSwitch")
	// render the component to the response writer
	component.Render(r.Context(), w)
}

func TestPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// helpers.Slogger.Info("Received request", "method", r.Method, "url", r.URL.String(), "protocol", r.Proto)

	// component := pages.IndexPage("My Todo List", structs.Menu)
	// serve text/html
	w.Header().Set("Content-Type", "text/html")
	// render the component to the response writer
	// component.Render(r.Context(), w)
	w.Write([]byte("Hello World"))
}
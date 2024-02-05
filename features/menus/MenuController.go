package menus

import (
	"fmt"
	"github.com/monstercameron/gofinances/features/monthlydebts"
	"github.com/monstercameron/gofinances/features/settings"
	"github.com/monstercameron/gofinances/helpers"
	"net/http"
	"strconv"
)

func init() {
	fmt.Println("MenuPicker.init(): \t\tInitializing MenuPicker...")
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

	SetActiveMenu(id)

	component := MainMenuComponent(GetMenus())
	// serve text/html
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "menuSwitch")
	// render the component to the response writer
	component.Render(r.Context(), w)
}

func GetTab(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := GetActiveMenu()
	if id == -1 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	fmt.Println("MenuPicker.GetTab(): id: ", id)

	switch id {
	case 0:
		component := monthlydebts.RecurringBillsIndex()
		// serve text/html
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
		return
	case 1:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Short term debts"))
		return
	case 2:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Long term debts"))
		return
	case 3:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Assets"))
		return
	case 4:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Credit utilization"))
		return
	case 5:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Goals"))
		return
	case 6:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Recommendations"))
		return
	case 7:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Calendar"))
		return
	case 8:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Drip calculator"))
		return
	case 9:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("time tables"))
		return
	case 10:
		component := settings.SettingsPageIndex(settings.GetAllSettingsUsers())
		w.Header().Set("Content-Type", "text/plain")
		component.Render(r.Context(), w)
		return
	default:
		// send string response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Default"))
		return
	}
}

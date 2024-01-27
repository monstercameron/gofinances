package controller

import (
	"fmt"
	"github.com/monstercameron/gofinances/structs"
	"github.com/monstercameron/gofinances/views/components"
	"net/http"
)

func AddDebt() {
	fmt.Println("test")
}

func SortRecurringBills(w http.ResponseWriter, r *http.Request) {

	sort := r.URL.Query().Get("sort")
	if sort != "" {
		structs.RecurringBills.SortBy(sort)
	}

	component := components.RecurringBillsComponent(structs.RecurringBills)
	// serve text/html
	w.Header().Set("Content-Type", "text/html")
	// render the component to the response writer
	component.Render(r.Context(), w)
}

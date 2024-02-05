package debts

import (
	"fmt"
	"net/http"
)

func init() {
	// Register the debts feature
	fmt.Println("Registering the debts feature")
}

func GetDebtsIndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting the debts index")
	component := DebtsIndex()
	// text/html
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

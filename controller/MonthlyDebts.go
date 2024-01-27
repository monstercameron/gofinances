package controller

import (
	"fmt"
	"github.com/monstercameron/gofinances/structs"
	"github.com/monstercameron/gofinances/views/components"
	"net/http"
	"strconv"
	"time"
)

func AddDebt() {
	fmt.Println("test")
}

func GetBills(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		// Convert id to integer
		intID, err := strconv.Atoi(id)
		if err != nil {
			// Handle error if id is not a valid integer
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		// Get bill by ID
		bill := structs.RecurringBills.GetByID(intID)
		if bill == nil {
			// Handle error if bill not found
			http.Error(w, "Bill not found", http.StatusNotFound)
			return
		}

		// Create a list with a single bill
		bills := structs.RecurringBillList{Bills: []structs.RecurrinBill{*bill}}

		// Serve text/html
		w.Header().Set("Content-Type", "text/html")
		// Render the component to the response writer
		components.RecurringBillsComponent(bills).Render(r.Context(), w)
	} else {

		sort := r.URL.Query().Get("sort")
		if sort != "" {
			structs.RecurringBills.SortBy(sort)
		}

		// Serve all bills when no ID is provided
		component := components.RecurringBillsComponent(structs.RecurringBills)
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
	}
}

func UpdateBills(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET and an ID is provided
	if r.Method == "GET" {
		fmt.Println("GET")
		id := r.URL.Query().Get("id")
		if id != "" {
			// Convert id to integer
			intID, err := strconv.Atoi(id)
			if err != nil {
				// Handle error if id is not a valid integer
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
			}
			bill := structs.RecurringBills.GetByID(intID)
			if err != nil {
				// Handle error if bill not found
				http.Error(w, "Bill not found", http.StatusNotFound)
			}
			w.Header().Set("Content-Type", "text/html")
			components.EditRecurringBillsComponent(*bill, true).Render(r.Context(), w)
		} else {
			newId := structs.RecurringBills.GetLastID()
			bill := structs.RecurrinBill{
				Id:         newId + 1,
				Name:       "",
				Amount:     0,
				DayOfMonth: 0,
				Owner:      "",
				Notes:      "",
			}
			w.Header().Set("Content-Type", "text/html")
			components.EditRecurringBillsComponent(bill, false).Render(r.Context(), w)
		}
		return
	}

	// 1. Parse URL Parameters
	billID := r.URL.Query().Get("id")
	if billID == "" {
		http.Error(w, "No bill ID provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(billID)
	if err != nil {
		http.Error(w, "Invalid bill ID", http.StatusBadRequest)
		return
	}

	// 2. Parse Form Data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	notes := r.FormValue("notes")

	// Convert date to day of month
	dayOfMonth, err := getDayOfMonth(date)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	// 3. Find and Update Bill
	bill := structs.RecurringBills.GetByID(id) // Assuming this function exists
	if err != nil {
		http.Error(w, "Bill not found", http.StatusNotFound)
		return
	}

	bill.Name = name
	bill.Amount = amount
	bill.DayOfMonth = dayOfMonth
	bill.Notes = notes

	// 4. Add to Bill List
	bills := structs.RecurringBillList{Bills: []structs.RecurrinBill{*bill}}

	// 5. Render Component
	w.Header().Set("Content-Type", "text/html")
	components.RecurringBillsComponent(bills).Render(r.Context(), w)
}

func AddBills(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Parse Form Data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "No name provided", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	notes := r.FormValue("notes")

	// Convert date to day of month
	dayOfMonth, err := getDayOfMonth(date)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	id := structs.RecurringBills.GetLastID() + 1

	bill := structs.RecurrinBill{
		Id:         id,
		Name:       name,
		Amount:     amount,
		DayOfMonth: dayOfMonth,
		Owner:      "cameron",
		Notes:      notes,
	}

	structs.RecurringBills.Bills = append(structs.RecurringBills.Bills, bill)

	// send headers
	w.Header().Set("HX-Trigger", "newBill")
	// reurn a 200 OK response
	w.WriteHeader(http.StatusOK)
	// write a message to the response writer
	fmt.Fprintf(w, "")
}

func DeleteBills(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse URL Parameters
	billID := r.URL.Query().Get("id")
	if billID == "" {
		http.Error(w, "No bill ID provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(billID)
	if err != nil {
		http.Error(w, "Invalid bill ID", http.StatusBadRequest)
		return
	}

	// 2. Find and Delete Bill
	bill := structs.RecurringBills.GetByID(id) // Assuming this function exists
	if bill == nil {
		http.Error(w, "Bill not found", http.StatusNotFound)
		return
	}

	// 3. Remove from Bill List
	structs.RecurringBills.RemoveByID(id)

	// 5. Render Component
	// reurn a 200 OK response
	w.WriteHeader(http.StatusOK)
	// write a message to the response writer
	fmt.Fprintf(w, "")
}

// Helper function to convert date string to day of month
func getDayOfMonth(dateStr string) (int, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, err
	}
	return date.Day(), nil
}

package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/monstercameron/gofinances/structs"
	"github.com/monstercameron/gofinances/views/components"
)

// GetBills handles the HTTP request to retrieve bill information.
// It can return a specific bill if an ID is provided, or all bills otherwise.
func GetBills(w http.ResponseWriter, r *http.Request) {
	// Extract the 'id' query parameter
	id := r.URL.Query().Get("id")

	// If 'id' parameter is provided
	if id != "" {
		// Convert 'id' to an integer
		intID, err := strconv.Atoi(id)
		if err != nil {
			// Respond with error if 'id' is not a valid integer
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		// Retrieve the bill by its ID
		bills := structs.RecurringBillList{Bills: []structs.RecurringBill{}}
		bills.Bills = append(bills.Bills, *bills.GetByID(intID))
		bill := bills.GetByID(intID)
		if bill == nil {
			// Respond with error if the bill is not found
			http.Error(w, "Bill not found", http.StatusNotFound)
			return
		}

		// Prepare a list containing the found bill
		bills.Bills = []structs.RecurringBill{*bill}

		// Set the Content-Type of the response to text/html
		w.Header().Set("Content-Type", "text/html")
		// Render the bill information as HTML to the response writer
		components.RecurringBillsComponent(bills).Render(r.Context(), w)
	} else {

		// Respond with all bills when no 'id' is provided
		bills := structs.RecurringBillList{Bills: []structs.RecurringBill{}}

		// Extract the 'sort' query parameter
		column := r.URL.Query().Get("column")
		order := r.URL.Query().Get("order")
		// Sort bills if 'sort' parameter is provided
		if order != "" {
			if order != "asc" && order != "desc" {
				http.Error(w, "Invalid sort order", http.StatusBadRequest)
				return
			}
		} else {
			// Default to ascending order if no order is provided
			order = "asc"
		}
		if column != "" {
			bills.SortBy(column, order)
		} else {
			// Default to sorting by name in ascending order
			bills.SortBy("name", "asc")
		}

		component := components.RecurringBillsComponent(bills)
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
	}
}

// UpdateBills handles the HTTP request for updating bill information.
// For a GET request, it fetches a specific bill for editing or prepares a new bill.
// For other request types, it updates a bill's details.
func UpdateBills(w http.ResponseWriter, r *http.Request) {
	// Handle GET request
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		if id != "" {
			// Fetch and edit an existing bill
			intID, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}
			// Retrieve the bill by its ID
			bills := structs.RecurringBillList{Bills: []structs.RecurringBill{}}
			bills.Bills = append(bills.Bills, *bills.GetByID(intID))
			bill := bills.GetByID(intID)
			if bill == nil {
				http.Error(w, "Bill not found", http.StatusNotFound)
				return
			}
			// trigger on page updates
			w.Header().Set("Hx-Trigger", "billsAction")
			w.Header().Set("Content-Type", "text/html")
			components.EditRecurringBillsComponent(*bill, true).Render(r.Context(), w)
		} else {
			// Set up a new bill
			// Retrieve the bill by its ID
			bills := structs.RecurringBillList{Bills: []structs.RecurringBill{}}
			newID := bills.GetLastID() + 1
			bill := structs.RecurringBill{
				Id:         newID,
				Name:       "",
				Amount:     0,
				DayOfMonth: 0,
				Owner:      "",
				Notes:      "",
			}
			// trigger on page updates
			w.Header().Set("Hx-Trigger", "billsAction")
			w.Header().Set("Content-Type", "text/html")
			components.EditRecurringBillsComponent(bill, false).Render(r.Context(), w)
		}
		return
	}

	// Handle POST request for updating a bill
	// Parse and validate the bill ID
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

	// Parse form data
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

	// Find and update the bill
	// Retrieve the bill by its ID
	bills := structs.RecurringBillList{Bills: []structs.RecurringBill{}}
	bills.Bills = append(bills.Bills, *bills.GetByID(id))
	bill := bills.GetByID(id)
	if bill == nil {
		http.Error(w, "Bill not found", http.StatusNotFound)
		return
	}
	bill.Name = name
	bill.Amount = amount
	bill.DayOfMonth = dayOfMonth
	bill.Notes = notes

	// Save the updated bill to the database
	bill.Save()

	// Render updated bill information
	w.Header().Set("Content-Type", "text/html")
	components.RecurringBillsComponent(structs.RecurringBillList{Bills: []structs.RecurringBill{*bill}}).Render(r.Context(), w)
}

// AddBills handles the HTTP POST request to add a new bill
// It parses the form data and adds a new bill to the list
func AddBills(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data from the request
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	// Retrieve and validate the 'name' field
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "No name provided", http.StatusBadRequest)
		return
	}

	// Retrieve and validate the 'amount' field
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Retrieve other form fields
	date := r.FormValue("date")
	notes := r.FormValue("notes")

	// Convert the date to a day of the month
	dayOfMonth, err := getDayOfMonth(date)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	// Generate a new ID for the bill
	bills := structs.RecurringBillList{Bills: []structs.RecurringBill{}}
	newID := bills.GetLastID() + 1

	// Create a new bill instance
	bill := structs.RecurringBill{
		Id:         newID,
		Name:       name,
		Amount:     amount,
		DayOfMonth: dayOfMonth,
		Owner:      "cameron", // This could be dynamic based on context
		Notes:      notes,
	}

	// save new bill to database
	bill.Save()

	// Set response headers and status
	w.Header().Set("HX-Trigger", "newBill, billsAction")
	w.WriteHeader(http.StatusOK)
	// Optional: Write a confirmation message to the response
	fmt.Fprintln(w, "")
}

// DeleteBills handles the HTTP DELETE request to remove a specific bill.
// It parses the bill ID from the URL parameters and removes the bill from the list.
func DeleteBills(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse and validate the bill ID from URL parameters
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

	// Find the bill by ID and handle if not found
	bills := structs.RecurringBillList{Bills: []structs.RecurringBill{}}
	bills.Bills = append(bills.Bills, *bills.GetByID(id))
	bill := bills.GetByID(id)
	if bill == nil {
		http.Error(w, "Bill not found", http.StatusNotFound)
		return
	}

	// Remove the bill from the list
	bills.RemoveByID(id)

	// trigger on page updates
	w.Header().Set("Hx-Trigger", "billsAction")

	// Set the response status to show a resource was deleted
	w.WriteHeader(http.StatusOK)
	// Optional: Write a confirmation message to the response
	fmt.Fprintln(w, "")
}

// Helper function to convert date string to day of month
func getDayOfMonth(dateStr string) (int, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, err
	}
	return date.Day(), nil
}

// returns total debts
func GetTotalDebts(w http.ResponseWriter, r *http.Request) {
	total := structs.GetTotalCost()
	// Set the Content-Type of the response to text/html
	w.Header().Set("Content-Type", "text/html")
	// Render the bill information as HTML to the response writer
	// set to 2 decimal places
	fmt.Fprintf(w, "%.2f", total)
}

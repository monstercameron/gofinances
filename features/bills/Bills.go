package bills

import (
	"database/sql"
	"fmt"
	"github.com/monstercameron/gofinances/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

// init is executed when the package is imported.
// It checks if the 'recurring_bills' table is empty and populates it with initial data if needed.
func init() {
	fmt.Println("RecuringBills.init(): \t\tchecking if 'recurring_bills' table is empty...")

	// Create table for recurring bills
	var err error
	_, err = database.DB.Exec("CREATE TABLE IF NOT EXISTS recurring_bills (id INTEGER PRIMARY KEY, name TEXT, amount REAL, day_of_month INTEGER, owner TEXT, notes TEXT)")
	if err != nil {
		log.Fatalf("Failed to create 'recurring_bills' table: %v", err)
	} else {
		fmt.Println("Database.Init(): \t\t'recurring_bills' table created.")
	}

	// Count the number of records in the 'recurring_bills' table
	var count int
	query := `SELECT count(*) FROM recurring_bills;`
	err = database.DB.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Printf("RecurringBills.init(): \t\tFailed to count records in 'recurring_bills' table: %v", err)
	}

	// If the table is empty (count = 0), populate it with initial data
	if count == 0 {
		fmt.Println("RecurringBills.init(): \t\ttable is empty. Populating with initial data...")
		InsertRecurringBillSamplesIntoDB()
	} else {
		fmt.Println("RecurringBills.init(): \t\ttable is not empty. Skipping population with initial data...")
	}
}

// RecurringBill represents a single recurring bill.
// It is used to store information about a bill that is paid on a regular basis.
type RecurringBill struct {
	Id         int
	Name       string
	Amount     float64
	DayOfMonth int
	Owner      string
	Notes      string
}

func (m *RecurringBill) Save() {
	fmt.Println("RecurringBill.Save()")
	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	var id int
	query := `SELECT id FROM recurring_bills WHERE id=?;`
	err = tx.QueryRow(query, m.Id).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		log.Fatal(err)
	}

	if err == sql.ErrNoRows {
		query = `INSERT INTO recurring_bills (name, amount, day_of_month, owner, notes) VALUES (?, ?, ?, ?, ?);`
		_, err = tx.Exec(query, m.Name, m.Amount, m.DayOfMonth, m.Owner, m.Notes)
	} else {
		query = `UPDATE recurring_bills SET name=?, amount=?, day_of_month=?, owner=?, notes=? WHERE id=?;`
		_, err = tx.Exec(query, m.Name, m.Amount, m.DayOfMonth, m.Owner, m.Notes, m.Id)
	}

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// RecurringBillList represents a list of RecurringBill.
type RecurringBillList struct {
	Bills []RecurringBill
}

// GetByID returns a pointer to the RecurringBill with the given ID.
// Returns nil if no bill is found.
func (m *RecurringBillList) GetByID(id int) *RecurringBill {
	var bill RecurringBill
	query := `SELECT id, name, amount, day_of_month, owner, notes FROM recurring_bills WHERE id=?;`
	err := database.DB.QueryRow(query, id).Scan(&bill.Id, &bill.Name, &bill.Amount, &bill.DayOfMonth, &bill.Owner, &bill.Notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}
	return &bill
}

func (m *RecurringBillList) GetAll() RecurringBillList {
	var bills []RecurringBill
	query := `SELECT id, name, amount, day_of_month, owner, notes FROM recurring_bills;`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var bill RecurringBill
		err := rows.Scan(&bill.Id, &bill.Name, &bill.Amount, &bill.DayOfMonth, &bill.Owner, &bill.Notes)
		if err != nil {
			log.Fatal(err)
		}
		bills = append(bills, bill)
	}
	err = rows.Err()
	if err != nil {
		log.Print(err)
		return RecurringBillList{}
	}
	return RecurringBillList{Bills: bills}
}

// RemoveByID removes the bill with the specified ID from the list.
func (m *RecurringBillList) RemoveByID(id int) {
	// delete bill from database
	query := `DELETE FROM recurring_bills WHERE id=?;`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		log.Fatal(err)
	}
}

// SortBy sorts the bills in the 'recurring_bills' table based on the provided criterion: name, amount, owner, or day of the month.
// The second parameter 'order' determines whether the sorting is in ascending (asc) or descending (desc) order.
// It populates the RecurringBillList with the sorted results.
func (m *RecurringBillList) SortBy(sortBy string, order string) {
	var query string

	// Validate the order parameter
	if order != "asc" && order != "desc" {
		log.Fatalf("Invalid order value: %s. It must be either 'asc' or 'desc'.", order)
	}

	// Construct the query based on the sortBy parameter
	switch sortBy {
	case "name", "amount", "owner", "day_of_month":
		query = fmt.Sprintf("SELECT * FROM recurring_bills ORDER BY %s %s;", sortBy, order)
	default:
		log.Fatalf("Invalid sortBy value: %s. It must be either 'name', 'amount', 'owner', or 'day_of_month'.", sortBy)
	}

	// Execute the query
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	// Clear existing bills
	m.Bills = []RecurringBill{}

	// Populate RecurringBillList with the sorted results
	for rows.Next() {
		var bill RecurringBill
		err := rows.Scan(&bill.Id, &bill.Name, &bill.Amount, &bill.DayOfMonth, &bill.Owner, &bill.Notes) // Ensure the order of fields matches your table schema
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		m.Bills = append(m.Bills, bill)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// GetLastID returns the ID of the last bill in the list.
func (m *RecurringBillList) GetLastID() int {
	// query reccuring_bills table for the last ID
	var id int
	query := `SELECT id FROM recurring_bills ORDER BY id DESC LIMIT 1;`
	err := database.DB.QueryRow(query).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Fatal(err)
	}
	return id
}

// PopulateRecurringBills creates and returns a RecurringBillList with initial data.
func PopulateRecurringBills() RecurringBillList {
	var bills RecurringBillList
	bills.Bills = append(bills.Bills, RecurringBill{Id: 0, Name: "Rent", DayOfMonth: 17, Amount: 1000, Owner: "Cameron", Notes: "Monthly rent for apartment"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 1, Name: "Car Payment", DayOfMonth: 2, Amount: 300, Owner: "Cameron", Notes: "Car loan payment"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 2, Name: "Car Insurance", DayOfMonth: 3, Amount: 100, Owner: "Cameron", Notes: "Quarterly car insurance"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 3, Name: "Internet Bill", DayOfMonth: 12, Amount: 60, Owner: "Alex", Notes: "Monthly internet service fee"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 4, Name: "Electricity Bill", DayOfMonth: 20, Amount: 75, Owner: "Taylor", Notes: "Monthly electricity usage"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 5, Name: "Gym Membership", DayOfMonth: 5, Amount: 40, Owner: "Jordan", Notes: "Annual gym subscription"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 6, Name: "Streaming Service", DayOfMonth: 25, Amount: 15, Owner: "Morgan", Notes: "Monthly Netflix subscription"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 7, Name: "Water Bill", DayOfMonth: 18, Amount: 50, Owner: "Casey", Notes: "Monthly water usage charge"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 8, Name: "Grocery Delivery", DayOfMonth: 1, Amount: 120, Owner: "Jamie", Notes: "Monthly grocery delivery service"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 9, Name: "Cell Phone", DayOfMonth: 11, Amount: 85, Owner: "Pat", Notes: "Monthly cell phone plan fee"})
	return bills
}

// InsertRecurringBillSamplesIntoDB creates and saves a RecurringBillList with initial data.
func InsertRecurringBillSamplesIntoDB() {
	RecurringBills := PopulateRecurringBills()

	tx, err := database.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, bill := range RecurringBills.Bills {
		bill.Save() // Assuming Save method is modified to accept a transaction
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// GetTotalCost calculates and returns the total cost of all bills.
// It returns 0 if there are no bills or in case of an error, logging the error internally.
func GetTotalCost() float64 {
	var totalCost sql.NullFloat64
	query := `SELECT SUM(amount) FROM recurring_bills;`
	err := database.DB.QueryRow(query).Scan(&totalCost)

	if err != nil {
		// Log the error for internal tracking
		log.Printf("error querying total cost: %v", err)
		return 0
	}

	// Check if totalCost is valid; if not, it means there are no rows/bills
	if !totalCost.Valid {
		return 0
	}

	return totalCost.Float64
}

// GetBills returns a RecurringBillList with all bills.
// It is used to populate the RecurringBillList with all bills.
func GetBills() RecurringBillList {
	// query recurring_bills table for all bills
	var bills RecurringBillList
	query := `SELECT * FROM recurring_bills;`
	rows, err := database.DB.Query(query)
	if err != nil {
		// Log the error and return an empty list
		log.Printf("error querying database: %v", err)
		return RecurringBillList{}
	}
	defer rows.Close()

	for rows.Next() {
		var bill RecurringBill
		err := rows.Scan(&bill.Id, &bill.Name, &bill.Amount, &bill.DayOfMonth, &bill.Owner, &bill.Notes)
		if err != nil {
			// Log the error and return an empty list
			log.Printf("error scanning row: %v", err)
			return RecurringBillList{}
		}
		bills.Bills = append(bills.Bills, bill)
	}

	err = rows.Err()
	if err != nil {
		// Log the error and return an empty list
		log.Printf("error after scanning rows: %v", err)
		return RecurringBillList{}
	}
	return bills
}

// GetBills handles the HTTP request to retrieve bill information.
// It can return a specific bill if an ID is provided, or all bills otherwise.
func GetManyBills(w http.ResponseWriter, r *http.Request) {
	// Respond with all bills when no 'id' is provided
	bills := RecurringBillList{Bills: []RecurringBill{}}

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

	component := RecurringBillsComponent(bills)
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

// GetBills handles the HTTP request to retrieve bill information.
// It can return a specific bill if an ID is provided, or all bills otherwise.
func GetOneBill(w http.ResponseWriter, r *http.Request) {
	// Extract the 'id' query parameter
	id := r.PathValue("id")

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
		bills := RecurringBillList{Bills: []RecurringBill{}}
		bills.Bills = append(bills.Bills, *bills.GetByID(intID))
		bill := bills.GetByID(intID)
		if bill == nil {
			// Respond with error if the bill is not found
			http.Error(w, "Bill not found", http.StatusNotFound)
			return
		}

		// Prepare a list containing the found bill
		bills.Bills = []RecurringBill{*bill}

		// Set the Content-Type of the response to text/html
		w.Header().Set("Content-Type", "text/html")
		// Render the bill information as HTML to the response writer
		RecurringBillsComponent(bills).Render(r.Context(), w)
	} 
}

func GetEditBillingComponent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id != "" {
		// Fetch and edit an existing bill
		intID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}
		// Retrieve the bill by its ID
		bills := RecurringBillList{Bills: []RecurringBill{}}
		bills.Bills = append(bills.Bills, *bills.GetByID(intID))
		bill := bills.GetByID(intID)
		if bill == nil {
			http.Error(w, "Bill not found", http.StatusNotFound)
			return
		}
		// trigger on page updates
		w.Header().Set("Hx-Trigger", "billsAction")
		w.Header().Set("Content-Type", "text/html")
		EditRecurringBillsComponent(*bill, true).Render(r.Context(), w)
	} else {
		// Set up a new bill
		// Retrieve the bill by its ID
		bills := RecurringBillList{Bills: []RecurringBill{}}
		newID := bills.GetLastID() + 1
		bill := RecurringBill{
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
		EditRecurringBillsComponent(bill, false).Render(r.Context(), w)
	}
	return
}

// UpdateBills handles the HTTP request for updating bill information.
// For a GET request, it fetches a specific bill for editing or prepares a new bill.
// For other request types, it updates a bill's details.
func UpdateBills(w http.ResponseWriter, r *http.Request) {
	// Handle POST request for updating a bill
	// Parse and validate the bill ID
	billID := r.PathValue("id")
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
	bills := RecurringBillList{Bills: []RecurringBill{}}
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
	w.Header().Set("HX-Trigger", "billsAction")
	w.Header().Set("Content-Type", "text/html")
	RecurringBillsComponent(RecurringBillList{Bills: []RecurringBill{*bill}}).Render(r.Context(), w)
}

// AddBills handles the HTTP POST request to add a new bill
// It parses the form data and adds a new bill to the list
func AddBills(w http.ResponseWriter, r *http.Request) {
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
	owner := r.FormValue("owner")

	// Convert the date to a day of the month
	dayOfMonth, err := getDayOfMonth(date)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	// Generate a new ID for the bill
	bills := RecurringBillList{Bills: []RecurringBill{}}
	newID := bills.GetLastID() + 1

	// Create a new bill instance
	bill := RecurringBill{
		Id:         newID,
		Name:       name,
		Amount:     amount,
		DayOfMonth: dayOfMonth,
		Owner:      owner, // This could be dynamic based on context
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
	// Parse and validate the bill ID from URL parameters
	billID := r.PathValue("id")
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
	bills := RecurringBillList{Bills: []RecurringBill{}}
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
func GetBillsTotalDebts(w http.ResponseWriter, r *http.Request) {
	total := GetTotalCost()
	// Set the Content-Type of the response to text/html
	w.Header().Set("Content-Type", "text/html")
	// Render the bill information as HTML to the response writer
	// set to 2 decimal places
	fmt.Fprintf(w, "%.2f", total)
}

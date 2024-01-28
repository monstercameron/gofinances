package structs

import (
	"database/sql"
	"fmt"
	"github.com/monstercameron/gofinances/database"
	"log"
)

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

// init is executed when the package is imported.
// It checks if the 'recurring_bills' table is empty and populates it with initial data if needed.
func init() {
	fmt.Println("Populating RecurringBills...")

	// Count the number of records in the 'recurring_bills' table
	var count int
	query := `SELECT count(*) FROM recurring_bills;`
	err := database.DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query the count from recurring_bills: %v", err)
	}

	// If the table is empty (count = 0), populate it with initial data
	if count == 0 {
		fmt.Println("RecurringBills is empty. Populating...")
		InsertRecurringBillSamplesIntoDB()
	} else {
		fmt.Println("RecurringBills is not empty. Skipping population.")
	}
}

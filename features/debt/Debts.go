package debt

import (
	"database/sql"
	"fmt"
	"github.com/monstercameron/gofinances/database"
	"net/http"
	"strconv"
)

func init() {
	// Register the debts feature
	fmt.Println("Registering the debts feature")
	// create debt table if it doesn't exist
	_, err := database.DB.Exec(`CREATE TABLE IF NOT EXISTS debts (
		id INTEGER PRIMARY KEY,
		debt_name TEXT,
		debt_owner TEXT,
		start_date TEXT,
		end_date TEXT,
		initial_amount REAL,
		current_amount REAL,
		interest_rate REAL,
		credit_limit REAL,
		notes TEXT \
	);`)
	if err != nil {
		fmt.Println("Failed to create 'debts' table: ", err)
	}
}

type Debt struct {
	ID            int     `json:"id"`
	DebtName      string  `json:"debtName"`
	DebtOwner     string  `json:"debtOwner"`
	StartDate     string  `json:"startDate"`
	EndDate       string  `json:"endDate"`
	InitialAmount float64 `json:"initialAmount"`
	CurrentAmount float64 `json:"currentAmount"`
	InterestRate  float64 `json:"interestRate"` // Added to reflect the structure from SQL
	CreditLimit   float64 `json:"creditLimit"`  // Added to reflect the structure from SQL
	Notes         string  `json:"notes"`
}

func (d *Debt) Save() error {
	fmt.Println("Debt.Save()")
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Printf("Error starting transaction: %v\n", err)
		return err
	}
	defer tx.Rollback() // Safely handle rollback if an error occurs

	var id int
	query := `SELECT id FROM debts WHERE id=?;`
	err = tx.QueryRow(query, d.ID).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("Error querying debt: %v\n", err)
		return err
	}

	if id == 0 {
		query := `INSERT INTO debts (debt_name, debt_owner, start_date, end_date, initial_amount, current_amount, interest_rate, credit_limit, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`
		_, err = tx.Exec(query, d.DebtName, d.DebtOwner, d.StartDate, d.EndDate, d.InitialAmount, d.CurrentAmount, d.InterestRate, d.CreditLimit, d.Notes)
		if err != nil {
			fmt.Printf("Error inserting debt: %v\n", err)
			return err
		}
	} else {
		query := `UPDATE debts SET debt_name=?, debt_owner=?, start_date=?, end_date=?, initial_amount=?, current_amount=?, interest_rate=?, credit_limit=?, notes=? WHERE id=?;`
		_, err = tx.Exec(query, d.DebtName, d.DebtOwner, d.StartDate, d.EndDate, d.InitialAmount, d.CurrentAmount, d.InterestRate, d.CreditLimit, d.Notes, d.ID)
		if err != nil {
			fmt.Printf("Error updating debt: %v\n", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		return err
	}
	return nil // Successful execution
}

func (d *Debt) Delete() {
	fmt.Println("Debt.Delete()")
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println("Error starting transaction: ", err)
	}
	query := `DELETE FROM debts WHERE id=?;`
	_, err = tx.Exec(query, d.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error deleting debt: ", err)
	}
	tx.Commit()
}

func GetDebt(id int) (Debt, error) {
	fmt.Println("GetDebt()")
	var d Debt
	query := `SELECT id, debt_name, debt_owner, start_date, end_date, initial_amount, current_amount, interest_rate, credit_limit, notes FROM debts WHERE id=?;`
	err := database.DB.QueryRow(query, id).Scan(&d.ID, &d.DebtName, &d.DebtOwner, &d.StartDate, &d.EndDate, &d.InitialAmount, &d.CurrentAmount, &d.InterestRate, &d.CreditLimit, &d.Notes)
	if err != nil {
		fmt.Printf("Error getting debt: %v\n", err)
		return Debt{}, err // Return zero value of Debt and the error
	}
	// Successfully retrieved and stored the debt in d
	return d, nil
}

func GetAllDebts() ([]Debt, error) {
	fmt.Println("GetAllDebts()")
	var debts []Debt
	query := `SELECT id, debt_name, debt_owner, start_date, end_date, initial_amount, current_amount, interest_rate, credit_limit, notes FROM debts;`
	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Printf("Error getting debts: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d Debt
		err = rows.Scan(&d.ID, &d.DebtName, &d.DebtOwner, &d.StartDate, &d.EndDate, &d.InitialAmount, &d.CurrentAmount, &d.InterestRate, &d.CreditLimit, &d.Notes)
		if err != nil {
			fmt.Printf("Error scanning debt: %v\n", err)
			return nil, err // Return immediately on error, providing what's been accumulated so far could be misleading
		}
		debts = append(debts, d)
	}
	if err = rows.Err(); err != nil { // Check for errors encountered during iteration
		fmt.Printf("Error iterating through debts: %v\n", err)
		return nil, err
	}

	fmt.Println(debts) // This might be a lot of output if there are many debts; consider removing or modifying this line.
	return debts, nil
}

func GetAllDebtsWithoutError() []Debt {
	fmt.Println("GetAllDebtsWithoutError()")
	debts, _ := GetAllDebts()
	return debts
}

func GetDebtsIndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting the debts index")
	component := DebtsIndex()
	// text/html
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

func GetDebtItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting a debt item")
	id := r.URL.Query().Get("id")
	switch {
	case id == "":
		// return debt item
		d := Debt{}
		component := DebtLineItem(d)
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
	default:
		// return debt item
		intID, err := strconv.Atoi(id)
		if err != nil {
			// Handle the error of Atoi conversion
			// For example, you might want to send a HTTP error response
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		d, err := GetDebt(intID)
		if err != nil {
			http.Error(w, "No Debt Found", http.StatusBadRequest)
			return
		}
		fmt.Println(d)
		component := DebtLineItem(d)
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
	}
}

func UpdateDebtItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting all debt items")
	id := r.URL.Query().Get("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		d, err := GetDebt(intID)
		if err != nil {
			http.Error(w, "No Debt Found", http.StatusBadRequest)
			return
		}
		component := EditDebts(d)
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
	case "POST":
		d, err := GetDebt(intID)
		if err != nil {
			http.Error(w, "Cant find Debt to Update", http.StatusBadRequest)
			return
		}

		dPointer := &d
		fmt.Println(dPointer)

		component := DebtLineItem(*dPointer)
		w.Header().Set("Content-Type", "text/html")
		component.Render(r.Context(), w)
		// default:
		// 	component := DebtLineItems(d)
		// 	w.Header().Set("Content-Type", "text/html")
		// 	component.Render(r.Context(), w)
	}
}

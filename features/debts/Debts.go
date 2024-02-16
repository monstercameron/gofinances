package debts

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
	_, err := database.DB.Exec("CREATE TABLE IF NOT EXISTS debts (id INTEGER PRIMARY KEY, name TEXT, owner TEXT, start_date TEXT, end_date TEXT, initial REAL, current REAL, notes TEXT)")
	if err != nil {
		fmt.Println("Failed to create 'debts' table: ", err)
	}
}

type Debt struct {
	Id        int
	Name      string
	Owner     string
	StartDate string
	EndDate   string
	Initial   float64
	Current   float64
	Notes     string
}

func (d *Debt) Save() {
	fmt.Println("Debt.Save()")
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println("Error starting transaction: ", err)
	}
	var id int
	query := `SELECT id FROM debts WHERE id=?;`
	err = tx.QueryRow(query, d.Id).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
	}
	if id == 0 {
		query := `INSERT INTO debts (name, owner, start_date, end_date, initial, current, notes) VALUES (?, ?, ?, ?, ?, ?, ?);`
		_, err = tx.Exec(query, d.Name, d.Owner, d.StartDate, d.EndDate, d.Initial, d.Current, d.Notes)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error inserting debt: ", err)
		}
	} else {
		query := `UPDATE debts SET name=?, owner=?, start_date=?, end_date=?, initial=?, current=?, notes=? WHERE id=?;`
		_, err = tx.Exec(query, d.Name, d.Owner, d.StartDate, d.EndDate, d.Initial, d.Current, d.Notes, d.Id)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error updating debt: ", err)
		}
	}
	tx.Commit()
}

func (d *Debt) Delete() {
	fmt.Println("Debt.Delete()")
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println("Error starting transaction: ", err)
	}
	query := `DELETE FROM debts WHERE id=?;`
	_, err = tx.Exec(query, d.Id)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error deleting debt: ", err)
	}
	tx.Commit()
}

func GetDebt(id int) (Debt, error) {
	fmt.Println("GetDebt()")
	var d Debt
	query := `SELECT id, name, owner, start_date, end_date, initial, current, notes FROM debts WHERE id=?;`
	err := database.DB.QueryRow(query, id).Scan(&d.Id, &d.Name, &d.Owner, &d.StartDate, &d.EndDate, &d.Initial, &d.Current, &d.Notes)
	if err != nil {
		fmt.Println("Error getting debt: ", err)
		return d, nil
	}
	// store the debt in d
	return d, nil
}

func GetAllDebts() ([]Debt, error) {
	fmt.Println("GetAllDebts()")
	var debts []Debt
	query := `SELECT id, name, owner, start_date, end_date, initial, current, notes FROM debts;`
	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Println("Error getting debts: ", err)
		return debts, err
	}
	defer rows.Close()
	for rows.Next() {
		var d Debt
		err = rows.Scan(&d.Id, &d.Name, &d.Owner, &d.StartDate, &d.EndDate, &d.Initial, &d.Current, &d.Notes)
		if err != nil {
			fmt.Println("Error scanning debt: ", err)
			return debts, err
		}
		debts = append(debts, d)
	}
	fmt.Println(debts)
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

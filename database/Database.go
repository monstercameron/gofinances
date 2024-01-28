package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// DB is a global variable for the database connection
var DB *sql.DB

func init() {
	fmt.Println("Database.init()")

	var err error
	DB, err = sql.Open("sqlite3", "./database/Database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table for recurring bills
	/* Based on structs/RecurringBills.go:
	type RecurringBill struct {
		Id         int
		Name       string
		Amount     float64
		DayOfMonth int
		Owner      string
		Notes      string
	}
	*/
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS recurring_bills (id INTEGER PRIMARY KEY, name TEXT, amount REAL, day_of_month INTEGER, owner TEXT, notes TEXT)")
	if err != nil {
		log.Fatal(err)
	}

}

func SimpleTest(DB *sql.DB) {
	res := DB.QueryRow("SELECT 1+1 AS solution")
	var solution int
	if err := res.Scan(&solution); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database.SimpleTest(): expected 2, got", solution)
}

func Test(DB *sql.DB) {
	fmt.Println("Database.Test()")

	// Check if table recurring_bills exists
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='recurring_bills';`
	err := DB.QueryRow(query).Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Table 'recurring_bills' does not exist.")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Table 'recurring_bills' exists.")
	}
}

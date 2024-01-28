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
    fmt.Println("Database.init(): \t\tInitializing database...")

    var err error
    DB, err = sql.Open("sqlite3", "./database/Database.db")
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    } else {
        fmt.Println("Database.Init(): \t\tDatabase connection established.")
    }

    // Create table for recurring bills
    _, err = DB.Exec("CREATE TABLE IF NOT EXISTS recurring_bills (id INTEGER PRIMARY KEY, name TEXT, amount REAL, day_of_month INTEGER, owner TEXT, notes TEXT)")
    if err != nil {
        log.Fatalf("Failed to create 'recurring_bills' table: %v", err)
    } else {
        fmt.Println("Database.Init(): \t\t'recurring_bills' table created.")
    }

    fmt.Println("Database.Init(): \t\tDatabase initialized.")
}


func SimpleTest(DB *sql.DB) {
	fmt.Println("Database.SimpleTest(): \t\tTesting database connection...")
	res := DB.QueryRow("SELECT 1+1 AS solution")
	var solution int
	if err := res.Scan(&solution); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database.SimpleTest(): \t\texpected 2, got", solution)
}

func Test(DB *sql.DB) {
	fmt.Println("Database.Test(): \t\tTesting database connection...")

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

package settings

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/monstercameron/gofinances/database"
)

func init() {
	fmt.Println("settings.init()")

	// create users table if not exists
	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT
	);`
	_, err := database.DB.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

// what settings do I want to store?
// add other people to the account
// choose a list of financial strategies
// add open ai API key
// export sqlite database
// import sqlite database

type SettingsPageUser struct {
	Name string
	Id   int
}

func (m *SettingsPageUser) Save() {
	fmt.Println("SettingsPageUser.Save()")
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println(err)
	}

	var id int
	query := `SELECT id FROM users WHERE id=?;`
	err = tx.QueryRow(query, m.Id).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		fmt.Println(err)
	}

	if err == sql.ErrNoRows {
		query = `INSERT INTO users (name) VALUES (?);`
		_, err = tx.Exec(query, m.Name)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
		}
	} else {
		query = `UPDATE users SET name=? WHERE id=?;`
		_, err = tx.Exec(query, m.Name, m.Id)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
	}
}

func (m *SettingsPageUser) Delete() {
	fmt.Println("SettingsPageUser.Delete()")
	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println(err)
	}

	query := `DELETE FROM users WHERE id=?;`
	_, err = tx.Exec(query, m.Id)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
	}
}

func GetSettingsPageUserById(id int) *SettingsPageUser {
	fmt.Println("SettingsPageUser.Get()")
	query := `SELECT id, name FROM users WHERE id=?;`
	row := database.DB.QueryRow(query, id)
	row.Scan(&id)
	var user SettingsPageUser
	err := row.Scan(&user.Id, &user.Name)
	if err != nil {
		fmt.Println(err)
	}
	return &user
}

func GetAllSettingsUsers() []SettingsPageUser {
	fmt.Println("GetAllUsers()")
	// Get all users from the database
	rows, err := database.DB.Query("SELECT id, name FROM users;")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	// Create a slice to hold the users
	users := []SettingsPageUser{}

	// Iterate over the rows, adding each user to the slice
	for rows.Next() {
		var user SettingsPageUser
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}

	return users
}

func GetSettingsUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSettingsUsers()")
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all users from the database
	rows, err := database.DB.Query("SELECT id, name FROM users;")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the users
	users := []SettingsPageUser{}

	// Iterate over the rows, adding each user to the slice
	for rows.Next() {
		var user SettingsPageUser
		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Set response headers and status
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Marshal the users to JSON and write to the response
	component := SettingsPageIndex(GetAllSettingsUsers())
	component.Render(r.Context(), w)
}

func GetSettingsUserInput(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSettingsUserInput()")
	// Ensure the request method is GET
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("id")
		if id == "" {

			// Set response headers and status
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)

			// Render the user input form
			component := SettingsPageUserInputField()
			component.Render(r.Context(), w)
		} else {
			intID, err := strconv.Atoi(id)
			if err != nil {
				// Respond with error if 'id' is not a valid integer
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}
			// Get the user from the database
			user := GetSettingsPageUserById(intID)

			// Set response headers and status
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)

			// Render the user input form
			component := SettingsPageUserInputFieldUpdate(user)
			component.Render(r.Context(), w)
		}
	} else if r.Method == http.MethodPost {
		name := r.FormValue("settingsusername")
		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		user := SettingsPageUser{}
		user.Name = name
		user.Save()

		// Set response headers and status
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		// Render the user input form
		component := SettingsPageListItem(&user)
		component.Render(r.Context(), w)
	} else if r.Method == http.MethodPut {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		intID, err := strconv.Atoi(id)
		if err != nil {
			// Respond with error if 'id' is not a valid integer
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}
		// Get the user from the database
		user := GetSettingsPageUserById(intID)

		name := r.FormValue("settingsusername")
		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		user.Name = name
		user.Save()

		// Set response headers and status
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		// Render the user input form
		component := SettingsPageUserInputFieldUpdate(user)
		component.Render(r.Context(), w)
	} else if r.Method == http.MethodDelete {
		fmt.Println("delete")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}
		intID, err := strconv.Atoi(id)
		if err != nil {
			// Respond with error if 'id' is not a valid integer
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}
		// Get the user from the database
		user := GetSettingsPageUserById(intID)
		user.Delete()

		// Set response headers and status
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("hx-trigger", "fetchSettingsUsers")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func GetSettingsUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSettingsUser()")
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user from the database
	users := GetAllSettingsUsers()

	// Set response headers and status
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Marshal the user to JSON and write to the response
	component := SettingsPageList(users)
	component.Render(r.Context(), w)
}

package settings

import (
	"database/sql"
	"fmt"
	"github.com/monstercameron/gofinances/database"
	"net/http"
	"strconv"
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

func (u *SettingsPageUser) Delete() error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM users WHERE id=?;`
	if _, err := tx.Exec(query, u.Id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return tx.Commit()
}

func GetSettingsPageUserById(id int) (*SettingsPageUser, error) {
	fmt.Println("SettingsPageUser.Get()")
	query := `SELECT id, name FROM users WHERE id=?;`
	row := database.DB.QueryRow(query, id)

	var user SettingsPageUser
	err := row.Scan(&user.Id, &user.Name)
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return nil, err
	}
	return &user, nil
}

func GetAllSettingsUsers() ([]SettingsPageUser, error) {
	query := `SELECT id, name FROM users;`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}
	defer rows.Close()

	var users []SettingsPageUser
	for rows.Next() {
		var user SettingsPageUser
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func GetAllSettingsUsersItems() []SettingsPageUser {
	users, err := GetAllSettingsUsers()
	if err != nil {
		return []SettingsPageUser{}
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
	users, err := GetAllSettingsUsers()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers and status
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Marshal the users to JSON and write to the response
	component := SettingsPageIndex(users)
	component.Render(r.Context(), w)
}

func GetSettingsUserActions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSettingsUser()")
	switch r.Method {
	case http.MethodGet:
		HandleGetSettingsUser(w, r)
	case http.MethodPost:
		HandlePostSettingsUser(w, r)
	case http.MethodPut:
		HandlePutSettingsUser(w, r)
	case http.MethodDelete:
		HandleDeleteSettingsUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleGetSettingsUser(w http.ResponseWriter, r *http.Request) {
	// Implementation of GET method
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
		user, err := GetSettingsPageUserById(intID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set response headers and status
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		// Render the user input form
		component := SettingsPageUserInputFieldUpdate(user)
		component.Render(r.Context(), w)
	}
}

func HandlePostSettingsUser(w http.ResponseWriter, r *http.Request) {
	// Implementation of POST method
	name := r.FormValue("settingsusername")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	user := SettingsPageUser{}
	user.Name = name
	user.Save()

	// Set response headers and status
	w.Header().Set("hx-trigger", "fetchSettingsUsers")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Render the user input form
	component := SettingsPageListItem(&user)
	component.Render(r.Context(), w)
}

func HandlePutSettingsUser(w http.ResponseWriter, r *http.Request) {
	// Implementation of PUT method
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
	user, err := GetSettingsPageUserById(intID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("settingsusername")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	user.Name = name
	user.Save()

	// Set response headers and status
	w.Header().Set("hx-trigger", "fetchSettingsUsers")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Render the user input form
	component := SettingsPageUserInputFieldUpdate(user)
	component.Render(r.Context(), w)
}

func HandleDeleteSettingsUser(w http.ResponseWriter, r *http.Request) {
	// Implementation of DELETE method
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
	user, err := GetSettingsPageUserById(intID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Delete()

	// Set response headers and status
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("hx-trigger", "fetchSettingsUsers")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func GetSettingsUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSettingsUser()")
	// Ensure the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user from the database
	users, err := GetAllSettingsUsers()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers and status
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Marshal the user to JSON and write to the response
	component := SettingsPageList(users)
	component.Render(r.Context(), w)
}

func GetSettingsPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SettingsController.GetSettingsPage(): \t\tGetting Settings Page...")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Settings Page from SettingsController.GetSettingsPage()"))
}

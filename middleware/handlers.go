package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pos-backend-go/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Admin struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var loginAdmin models.LoginAdmin
	err := json.NewDecoder(r.Body).Decode(&loginAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query the database to validate the admin login
	var storedPassword string
	err = db.QueryRow("SELECT password FROM admin_table WHERE email = $1", loginAdmin.Email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Incorrect Credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Compare the stored password with the provided password (You should hash passwords)
	if loginAdmin.Password != storedPassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Admin login successful")
}

func AdminChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var changePassword models.LoginAdmin
	err := json.NewDecoder(r.Body).Decode(&changePassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate the admin (you can use session management or tokens)
	// Here, we'll assume the admin is authenticated

	// Hash the new password (You should hash passwords)
	newPassword := changePassword.Password

	// Update the password in the database
	_, err = db.Exec("UPDATE admin_table SET password = $1 WHERE email = $2", newPassword, changePassword.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Password changed successfully")
}

// func AddAdminCredentials(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// create the postgres db connection
// 	db := createConnection()

// 	// close the db connection
// 	defer db.Close()

// 	var newAdmin models.LoginAdmin
// 	err := json.NewDecoder(r.Body).Decode(&newAdmin)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Insert the new admin credentials into the database
// 	_, err = db.Exec("INSERT INTO admin_table (email, password) VALUES ($1, $2)", newAdmin.Email, newAdmin.Password)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintln(w, "Admin credentials added successfully")
// }

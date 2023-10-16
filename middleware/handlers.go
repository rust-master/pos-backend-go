package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pos-backend-go/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Admin struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type response struct {
	Message string `json:"message,omitempty"`
	Jwt     string `json:"jwt"`
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create the postgres db connection
	db := createConnection()
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

	// Your JWT secret key should be kept secret and not hard-coded here. You can load it from an environment variable.
	jwtSecret := []byte("your-secret-key")

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": loginAdmin.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (adjust as needed)
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "could not create JWT", http.StatusInternalServerError)
		return
	}

	// Set the token as a cookie
	cookie := http.Cookie{
		Name:    "myjwt",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 48),
	}

	http.SetCookie(w, &cookie)

	res := response{
		Message: "Admin login successful",
		Jwt:     cookie.Value,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func VerifyJWT(jwtString string) (*jwt.Token, error) {
	// Your JWT secret key should be kept secret and not hard-coded here. You can load it from an environment variable.
	jwtSecret := []byte("your-secret-key")

	// Parse the JWT token using the secret key
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func AdminChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// if r.Method != http.MethodPut {
	// 	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 	return
	// }

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusBadRequest)
		return
	}

	// Split the header value to get the actual token part
	// The header value should be in the format "Bearer <token>"
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	jwtToken := authHeaderParts[1]

	token, errt := VerifyJWT(jwtToken)

	if errt != nil {
		http.Error(w, errt.Error(), http.StatusUnauthorized)
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

	res := response{
		Message: "Password changed successfully",
		Jwt:     token.Signature,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
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

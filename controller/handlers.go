package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pos-backend-go/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	var err error
	connStr := "user=postgres dbname=pos_backend_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}
	fmt.Println("Database connected successfully")

	return db
}

// func createConnection() *sql.DB {
// 	err := godotenv.Load(".env")

// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}

// 	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

// 	if err != nil {
// 		panic(err)
// 	}

// 	err = db.Ping()

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Successfully connected!")
// 	return db
// }

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

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	var newProduct models.Products
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new admin credentials into the database
	_, err = db.Exec("INSERT INTO product_table (name, description, unitprice, unitsinstock) VALUES ($1, $2, $3, $4)", newProduct.Name, newProduct.Description, newProduct.UnitPrice, newProduct.UnitsInStock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := response{
		Message: "Product Added successfully",
		Jwt:     token.Signature,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a database connection
	db := createConnection()
	defer db.Close()

	// Query to select all products
	rows, err := db.Query("SELECT * FROM product_table ORDER BY product_code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Products // Assuming you have a "models.Products" struct

	// Iterate through the rows and scan into the products slice
	for rows.Next() {
		var product models.Products
		if err := rows.Scan(&product.ProductCode, &product.Name, &product.Description, &product.UnitPrice, &product.UnitsInStock); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the products slice to JSON and send it as the response
	jsonResponse, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func GetProductByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the product code from the URL parameter
	productCode := r.URL.Query().Get("product_code")

	// Check if the product code is provided
	if productCode == "" {
		http.Error(w, "Product code is missing in the request", http.StatusBadRequest)
		return
	}

	// Create a database connection
	db := createConnection()

	// Close the database connection when done
	defer db.Close()

	// Query to select the product by its code
	row := db.QueryRow("SELECT * FROM product_table WHERE product_code = $1", productCode)

	var product models.Products // Assuming you have a "models.Product" struct

	// Scan the row into the product struct
	if err := row.Scan(&product.ProductCode, &product.Name, &product.Description, &product.UnitPrice, &product.UnitsInStock); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Convert the product to JSON and send it as the response
	jsonResponse, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusBadRequest)
		return
	}

	// Split the header value to get the actual token part
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

	// Create the Postgres DB connection
	db := createConnection()
	defer db.Close()

	var updateProduct models.Products
	err := json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check for product_code in the query parameters
	productCode := r.URL.Query().Get("product_code")
	if productCode == "" {
		http.Error(w, "Missing product_code query parameter", http.StatusBadRequest)
		return
	}

	// Update the product in the database
	_, err = db.Exec(
		"UPDATE product_table SET name = $1, description = $2, unitprice = $3, unitsinstock = $4 WHERE product_code = $5",
		updateProduct.Name,
		updateProduct.Description,
		updateProduct.UnitPrice,
		updateProduct.UnitsInStock,
		productCode,
	)
	if err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := response{
		Message: "Product updated successfully",
		Jwt:     token.Signature,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func DeleteProductByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	// Get the product code from the URL parameter
	product_code := r.URL.Query().Get("product_code")

	// Check if the product code is provided
	if product_code == "" {
		http.Error(w, "Product code is missing in the request", http.StatusBadRequest)
		return
	}

	// Create a database connection
	db := createConnection()

	// Close the database connection when done
	defer db.Close()

	// Delete the product from the database based on the product code
	_, err := db.Exec("DELETE FROM product_table WHERE product_code = $1", product_code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := response{
		Message: "Product with code " + product_code + " deleted successfully",
		Jwt:     token.Signature,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

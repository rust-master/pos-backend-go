package router

import (
	"pos-backend-go/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/products", middleware.GetAllProducts).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/getProductByCode", middleware.GetProductByCode).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/products", middleware.CreateProduct).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/delete-product", middleware.DeleteProductByCode).Methods("DELETE", "OPTIONS")

	// router.HandleFunc("/api/customers", middleware.CreateUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/customers", middleware.CreateUser).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/orders", middleware.CreateUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/orders", middleware.CreateUser).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/invoices", middleware.CreateUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/invoices", middleware.CreateUser).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/invoice-lines", middleware.CreateUser).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/invoice-lines", middleware.CreateUser).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/login", middleware.AdminLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/change-password", middleware.AdminChangePassword).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/addadmin", middleware.AddAdminCredentials).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}

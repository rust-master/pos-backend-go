package router

import (
	"pos-backend-go/controller"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/products", controller.GetAllProducts).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/getProductByCode", controller.GetProductByCode).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/products", controller.CreateProduct).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/delete-product", controller.DeleteProductByCode).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/update-product", controller.UpdateProduct).Methods("PUT", "OPTIONS")

	// router.HandleFunc("/api/customers", middleware.CreateUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/customers", middleware.CreateUser).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/orders", middleware.CreateUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/orders", middleware.CreateUser).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/invoices", middleware.CreateUser).Methods("GET", "OPTIONS")
	// router.HandleFunc("/api/invoices", middleware.CreateUser).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/invoice-lines", middleware.CreateUser).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/invoice-lines", middleware.CreateUser).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/login", controller.AdminLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/change-password", controller.AdminChangePassword).Methods("POST", "OPTIONS")
	// router.HandleFunc("/api/addadmin", middleware.AddAdminCredentials).Methods("POST", "OPTIONS")

	// router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	// router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}
